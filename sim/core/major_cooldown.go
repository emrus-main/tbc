package core

import (
	"sort"
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const (
	CooldownPriorityLow       = -1.0
	CooldownPriorityDefault   = 0.0
	CooldownPriorityDrums     = 2.0
	CooldownPriorityBloodlust = 1.0
)

const (
	CooldownTypeUnknown = 0
	CooldownTypeMana    = 1
	CooldownTypeDPS     = 2
)

// Condition for whether a cooldown can/should be activated.
// Returning false prevents the cooldown from being activated.
type CooldownActivationCondition func(*Simulation, *Character) bool

// Function for activating a cooldown.
// Returns whether the activation was successful.
type CooldownActivation func(*Simulation, *Character)

// Function for making a CooldownActivation.
//
// We need a function that returns a CooldownActivation rather than a
// CooldownActivation, so captured local variables can be reset on Sim reset.
type CooldownActivationFactory func(*Simulation) CooldownActivation

type MajorCooldown struct {
	// Spell that is cast when this MCD is activated.
	Spell *Spell

	// Unique ID for this cooldown, used to look it up.
	ActionID

	// Primary cooldown for checking whether this cooldown is ready.
	Cooldown Cooldown

	// Secondary cooldown, used for shared cooldowns.
	SharedCooldown Cooldown

	// Whether this MCD requires the GCD.
	UsesGCD bool

	// How long before this cooldown takes effect after activation.
	// Not used yet, but will eventually be important for planning cooldown
	// schedules.
	CastTime time.Duration

	// Cooldowns with higher priority get used first. This is important when some
	// cooldowns have a non-zero cast time. For example, Drums should be used
	// before Bloodlust.
	Priority float64

	// Internal category, used for filtering. For example, mages want to disable
	// all DPS cooldowns during their regen rotation.
	Type int32

	// Whether the cooldown meets all hard requirements for activation (e.g. resource cost).
	// Note chat whether the cooldown is off CD is automatically checked, so it does not
	// need to be checked again by this function.
	CanActivate CooldownActivationCondition

	// Whether the cooldown meets all optional conditions for activation. These
	// conditions will be ignored when the user specifies their own activation time.
	// This is for things like mana thresholds, which are optimizations for better
	// automatic timing.
	ShouldActivate CooldownActivationCondition

	// Factory for creating the activate function on every Sim reset.
	ActivationFactory CooldownActivationFactory

	// Fixed timings at which to use this cooldown. If these are specified, they
	// are used instead of ShouldActivate.
	timings []time.Duration

	// Number of times this MCD was used so far in the current iteration.
	numUsages int

	// Internal lambda function to use the cooldown.
	activate CooldownActivation

	// Whether this MCD is currently disabled.
	disabled bool
}

func (mcd *MajorCooldown) ReadyAt() time.Duration {
	return BothTimersReadyAt(mcd.Cooldown.Timer, mcd.SharedCooldown.Timer)
}

func (mcd *MajorCooldown) IsReady(sim *Simulation) bool {
	return BothTimersReady(mcd.Cooldown.Timer, mcd.SharedCooldown.Timer, sim)
}

func (mcd *MajorCooldown) TimeToReady(sim *Simulation) time.Duration {
	return MaxTimeToReady(mcd.Cooldown.Timer, mcd.SharedCooldown.Timer, sim)
}

func (mcd *MajorCooldown) IsEnabled() bool {
	return !mcd.disabled
}

func (mcd *MajorCooldown) GetTimings() []time.Duration {
	return mcd.timings[:]
}

// Public version of TryActivate for manual activation by Agent code.
// Note that this version will work even if the MCD is disabled.
func (mcd *MajorCooldown) TryActivate(sim *Simulation, character *Character) bool {
	return mcd.tryActivateHelper(sim, character)
}

func (mcd *MajorCooldown) tryActivateInternal(sim *Simulation, character *Character) bool {
	if mcd.disabled {
		return false
	}

	return mcd.tryActivateHelper(sim, character)
}

// Activates this MCD, if all the conditions pass.
// Returns whether the MCD was activated.
func (mcd *MajorCooldown) tryActivateHelper(sim *Simulation, character *Character) bool {
	if mcd.UsesGCD && !character.GCD.IsReady(sim) {
		return false
	}

	if !mcd.CanActivate(sim, character) {
		return false
	}

	var shouldActivate bool
	if mcd.numUsages < len(mcd.timings) {
		shouldActivate = sim.CurrentTime >= mcd.timings[mcd.numUsages]
	} else {
		shouldActivate = mcd.ShouldActivate(sim, character)
	}

	if shouldActivate {
		mcd.activate(sim, character)
		mcd.numUsages++
		if sim.Log != nil {
			character.Log(sim, "Major cooldown used: %s", mcd.ActionID)
		}
	}

	return shouldActivate
}

type majorCooldownManager struct {
	// The Character whose cooldowns are being managed.
	character *Character

	// User-specified cooldown configs.
	cooldownConfigs proto.Cooldowns

	// Cached list of major cooldowns sorted by priority, for resetting quickly.
	initialMajorCooldowns []MajorCooldown

	// Major cooldowns, ordered by next available. This should always contain
	// the same cooldows as initialMajorCooldowns, but the order will change over
	// the course of the sim.
	majorCooldowns []*MajorCooldown
}

func newMajorCooldownManager(cooldowns *proto.Cooldowns) majorCooldownManager {
	cds := proto.Cooldowns{}
	if cooldowns != nil {
		cds = *cooldowns
	}

	return majorCooldownManager{
		cooldownConfigs: cds,
	}
}

func (mcdm *majorCooldownManager) initialize(character *Character) {
	mcdm.character = character
}

func (mcdm *majorCooldownManager) finalize(character *Character) {
	if mcdm.initialMajorCooldowns == nil {
		mcdm.initialMajorCooldowns = []MajorCooldown{}
	}

	// Match user-specified cooldown configs to existing cooldowns.
	for i, _ := range mcdm.initialMajorCooldowns {
		mcd := &mcdm.initialMajorCooldowns[i]
		mcd.timings = []time.Duration{}

		if mcdm.cooldownConfigs.Cooldowns != nil {
			for _, cooldownConfig := range mcdm.cooldownConfigs.Cooldowns {
				configID := ProtoToActionID(*cooldownConfig.Id)
				if configID.SameAction(mcd.ActionID) {
					mcd.timings = make([]time.Duration, len(cooldownConfig.Timings))
					for t, timing := range cooldownConfig.Timings {
						mcd.timings[t] = DurationFromSeconds(timing)
					}
					break
				}
			}
		}
	}

	mcdm.majorCooldowns = make([]*MajorCooldown, len(mcdm.initialMajorCooldowns))
}

// Adds a delay to the first usage of all CDs so that armor debuffs have time
// to be applied. MCDs that have a user-specified timing are not delayed.
//
// This function should be called from Agent.Init().
func (mcdm *majorCooldownManager) DelayDPSCooldownsForArmorDebuffs() {
	if !mcdm.character.Env.GetPrimaryTarget().HasAuraWithTag(SunderExposeAuraTag) {
		return
	}

	const delay = time.Second * 10
	for i, _ := range mcdm.initialMajorCooldowns {
		mcd := &mcdm.initialMajorCooldowns[i]
		if len(mcd.timings) == 0 && mcd.Type == CooldownTypeDPS {
			mcd.timings = append(mcd.timings, delay)
		}
	}
}

func (mcdm *majorCooldownManager) reset(sim *Simulation) {
	// Need to create all cooldowns before calling ActivationFactory on any of them,
	// so that any cooldown can do lookups on other cooldowns.
	for i, _ := range mcdm.majorCooldowns {
		newMCD := &MajorCooldown{}
		*newMCD = mcdm.initialMajorCooldowns[i]
		mcdm.majorCooldowns[i] = newMCD
	}

	for i, _ := range mcdm.majorCooldowns {
		mcdm.majorCooldowns[i].activate = mcdm.majorCooldowns[i].ActivationFactory(sim)
		if mcdm.majorCooldowns[i].activate == nil {
			panic("Nil cooldown activation returned!")
		}
	}

	// For initial sorting.
	mcdm.UpdateMajorCooldowns()
}

// Registers a major cooldown to the Character, which will be automatically
// used when available.
func (mcdm *majorCooldownManager) AddMajorCooldown(mcd MajorCooldown) {
	if mcdm.character.Env.IsFinalized() {
		panic("Major cooldowns may not be added once finalized!")
	}

	spell := mcd.Spell
	if spell != nil {
		mcd.ActionID = spell.ActionID
		mcd.Cooldown = spell.CD
		mcd.SharedCooldown = spell.SharedCD
		mcd.UsesGCD = spell.DefaultCast.GCD > 0
		mcd.CastTime = spell.DefaultCast.CastTime
		if mcd.ActivationFactory == nil {
			mcd.ActivationFactory = func(sim *Simulation) CooldownActivation {
				return func(sim *Simulation, character *Character) {
					spell.Cast(sim, sim.GetPrimaryTarget())
				}
			}
		}
	}

	if mcd.IsEmptyAction() {
		panic("Major cooldown must have an ActionID!")
	}

	if mcd.CanActivate == nil {
		mcd.CanActivate = func(sim *Simulation, character *Character) bool {
			return true
		}
	}
	if mcd.ShouldActivate == nil {
		mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
			return true
		}
	}

	mcdm.initialMajorCooldowns = append(mcdm.initialMajorCooldowns, mcd)
}

func (mcdm *majorCooldownManager) GetInitialMajorCooldown(actionID ActionID) MajorCooldown {
	for _, mcd := range mcdm.initialMajorCooldowns {
		if mcd.SameAction(actionID) {
			return mcd
		}
	}

	return MajorCooldown{}
}

func (mcdm *majorCooldownManager) GetMajorCooldown(actionID ActionID) *MajorCooldown {
	for _, mcd := range mcdm.majorCooldowns {
		if mcd.SameAction(actionID) {
			return mcd
		}
	}

	return nil
}

// Returns all MCDs.
func (mcdm *majorCooldownManager) GetMajorCooldowns() []*MajorCooldown {
	return mcdm.majorCooldowns
}

func (mcdm *majorCooldownManager) GetMajorCooldownIDs() []*proto.ActionID {
	ids := make([]*proto.ActionID, len(mcdm.initialMajorCooldowns))
	for i, mcd := range mcdm.initialMajorCooldowns {
		ids[i] = mcd.ActionID.ToProto()
	}
	return ids
}

func (mcdm *majorCooldownManager) HasMajorCooldown(actionID ActionID) bool {
	return mcdm.GetMajorCooldown(actionID) != nil
}

func (mcdm *majorCooldownManager) DisableMajorCooldown(actionID ActionID) {
	mcd := mcdm.GetMajorCooldown(actionID)
	if mcd != nil {
		mcd.disabled = true
	}
}

func (mcdm *majorCooldownManager) EnableMajorCooldown(actionID ActionID) {
	mcd := mcdm.GetMajorCooldown(actionID)
	if mcd != nil {
		mcd.disabled = false
	}
}

// Disabled all MCDs that are currently enabled, and returns a list of the MCDs
// which were disabled by this call.
// If cooldownType is not CooldownTypeUnknown, then will be restricted to cooldowns of that type.
func (mcdm *majorCooldownManager) DisableAllEnabledCooldowns(cooldownType int32) []*MajorCooldown {
	disabledMCDs := []*MajorCooldown{}
	for _, mcd := range mcdm.majorCooldowns {
		if mcd.IsEnabled() && (cooldownType == CooldownTypeUnknown || mcd.Type == cooldownType) {
			mcdm.DisableMajorCooldown(mcd.ActionID)
			disabledMCDs = append(disabledMCDs, mcd)
		}
	}
	return disabledMCDs
}

func (mcdm *majorCooldownManager) EnableAllCooldowns(mcdsToEnable []*MajorCooldown) {
	for _, mcd := range mcdsToEnable {
		mcdm.EnableMajorCooldown(mcd.ActionID)
	}
}

func (mcdm *majorCooldownManager) TryUseCooldowns(sim *Simulation) {
	anyCooldownsUsed := false
	for curIdx := 0; curIdx < len(mcdm.majorCooldowns) && mcdm.majorCooldowns[curIdx].IsReady(sim); curIdx++ {
		mcd := mcdm.majorCooldowns[curIdx]
		if mcd.tryActivateInternal(sim, mcdm.character) {
			anyCooldownsUsed = true

			if mcd.UsesGCD {
				// If we used a MCD that uses the GCD (like drums), hold off on using
				// any remaining MCDs so they aren't wasted.
				break
			}
		}
	}

	if anyCooldownsUsed {
		// Re-sort by availability.
		// TODO: Probably a much faster way to do this, especially since we know which cooldowns need to be re-ordered.
		// Maybe not because MCDs with shared cooldowns put each other on CD.
		mcdm.UpdateMajorCooldowns()
	}
}

// This function should be called if the CD for a major cooldown changes outside
// of the TryActivate() call.
func (mcdm *majorCooldownManager) UpdateMajorCooldowns() {
	sort.Slice(mcdm.majorCooldowns, func(i, j int) bool {
		// Since we're just comparing and don't actually care about the remaining CD, ok to use 0 instead of sim.CurrentTime.
		cdA := mcdm.majorCooldowns[i].ReadyAt()
		cdB := mcdm.majorCooldowns[j].ReadyAt()
		return cdA < cdB || (cdA == cdB && mcdm.majorCooldowns[i].Priority > mcdm.majorCooldowns[j].Priority)
	})
}

// Add a major cooldown to the given agent, which provides a temporary boost to a single stat.
// This is use for effects like Icon of the Silver Crescent and Bloodlust Brooch.
func RegisterTemporaryStatsOnUseCD(character *Character, auraLabel string, tempStats stats.Stats, duration time.Duration, config SpellConfig) {
	aura := character.NewTemporaryStatsAura(auraLabel, config.ActionID, tempStats, duration)

	config.Cast.DisableCallbacks = true
	config.ApplyEffects = func(sim *Simulation, _ *Target, _ *Spell) {
		aura.Activate(sim)
	}
	spell := character.RegisterSpell(config)

	character.AddMajorCooldown(MajorCooldown{
		Spell: spell,
		Type:  CooldownTypeDPS,
	})
}

// Helper function to make an ApplyEffect for a temporary stats on-use cooldown.
func MakeTemporaryStatsOnUseCDRegistration(auraLabel string, tempStats stats.Stats, duration time.Duration, config SpellConfig, cdFunc func(*Character) Cooldown, sharedCDFunc func(*Character) Cooldown) ApplyEffect {
	return func(agent Agent) {
		localConfig := config
		character := agent.GetCharacter()
		if cdFunc != nil {
			localConfig.Cast.CD = cdFunc(character)
		}
		if sharedCDFunc != nil {
			localConfig.Cast.SharedCD = sharedCDFunc(character)
		}
		RegisterTemporaryStatsOnUseCD(character, auraLabel, tempStats, duration, localConfig)
	}
}

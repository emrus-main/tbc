package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func NewShaman(character core.Character, talents proto.ShamanTalents, selfBuffs SelfBuffs) Shaman {
	shaman := Shaman{
		Character: character,
		Talents:   talents,
		SelfBuffs: selfBuffs,
	}

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Intellect,
		ModifiedStat: stats.SpellCrit,
		Modifier: func(intellect float64, spellCrit float64) float64 {
			return spellCrit + (intellect/78.1)*core.SpellCritRatingPerCritChance
		},
	})

	if shaman.Talents.UnrelentingStorm > 0 {
		coeff := 0.02 * float64(shaman.Talents.UnrelentingStorm)
		shaman.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.MP5,
			Modifier: func(intellect float64, mp5 float64) float64 {
				return mp5 + intellect*coeff
			},
		})
	}

	if selfBuffs.WaterShield {
		shaman.AddStat(stats.MP5, 50)
	}

	shaman.registerElementalMasteryCD()

	return shaman
}

// Which buffs this shaman is using.
type SelfBuffs struct {
	Bloodlust    bool
	WaterShield  bool
	TotemOfWrath bool
	WrathOfAir   bool
	ManaSpring   bool
}

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	Talents   proto.ShamanTalents
	SelfBuffs SelfBuffs

	ElementalFocusStacks byte
	NextTotemDrop        time.Duration // track when to drop totems again

	// "object pool" for shaman spells that are currently being cast.
	lightningBoltSpell   core.SingleTargetDirectDamageSpell
	lightningBoltSpellLO core.SingleTargetDirectDamageSpell

	chainLightningSpell    core.MultiTargetDirectDamageSpell
	chainLightningSpellLOs []core.MultiTargetDirectDamageSpell

	// Precomputed templated cast generator for quickly resetting cast fields.
	lightningBoltCastTemplate   core.SingleTargetDirectDamageSpellTemplate
	lightningBoltLOCastTemplate core.SingleTargetDirectDamageSpellTemplate

	chainLightningCastTemplate    core.MultiTargetDirectDamageSpellTemplate
	chainLightningLOCastTemplates []core.MultiTargetDirectDamageSpellTemplate
}

// Implemented by each Shaman spec.
type ShamanAgent interface {
	core.Agent

	// The Shaman controlled by this Agent.
	GetShaman() *Shaman
}

func (shaman *Shaman) GetCharacter() *core.Character {
	return &shaman.Character
}

func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}
func (shaman *Shaman) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
	if shaman.SelfBuffs.Bloodlust {
		partyBuffs.Bloodlust += 1
	}

	if shaman.SelfBuffs.TotemOfWrath {
		partyBuffs.TotemOfWrath += 1
	}

	if shaman.SelfBuffs.ManaSpring {
		partyBuffs.ManaSpringTotem = core.MaxTristate(partyBuffs.ManaSpringTotem, proto.TristateEffect_TristateEffectRegular)
	}

	if shaman.SelfBuffs.WrathOfAir {
		woaValue := proto.TristateEffect_TristateEffectRegular
		if ItemSetCycloneRegalia.CharacterHasSetBonus(shaman.GetCharacter(), 2) {
			woaValue = proto.TristateEffect_TristateEffectImproved
		}
		partyBuffs.WrathOfAirTotem = core.MaxTristate(partyBuffs.WrathOfAirTotem, woaValue)
	}
}

func (shaman *Shaman) Init(sim *core.Simulation) {
	// Precompute all the spell templates.
	shaman.lightningBoltCastTemplate = shaman.newLightningBoltTemplate(sim, false)
	shaman.lightningBoltLOCastTemplate = shaman.newLightningBoltTemplate(sim, true)

	shaman.chainLightningCastTemplate = shaman.newChainLightningTemplate(sim, false)

	numHits := core.MinInt32(3, sim.GetNumTargets())
	shaman.chainLightningSpellLOs = make([]core.MultiTargetDirectDamageSpell, numHits)
	shaman.chainLightningLOCastTemplates = []core.MultiTargetDirectDamageSpellTemplate{}
	for i := int32(0); i < numHits; i++ {
		shaman.chainLightningLOCastTemplates = append(shaman.chainLightningLOCastTemplates, shaman.newChainLightningTemplate(sim, true))
	}
}

func (shaman *Shaman) Reset(sim *core.Simulation) {
	shaman.Character.Reset(sim)

	shaman.NextTotemDrop = time.Second * 120 // 2 min until drop totems

	// Reset all spells so any pending casts are cleaned up
	shaman.lightningBoltSpell = core.SingleTargetDirectDamageSpell{}
	shaman.lightningBoltSpellLO = core.SingleTargetDirectDamageSpell{}
	shaman.chainLightningSpell = core.MultiTargetDirectDamageSpell{}

	numHits := core.MinInt32(3, sim.GetNumTargets())
	shaman.chainLightningSpellLOs = make([]core.MultiTargetDirectDamageSpell, numHits)
}

func (shaman *Shaman) Advance(sim *core.Simulation, elapsedTime time.Duration) {
	// Shaman should never be outside the 5s window, use combat regen
	shaman.Character.CombatManaRegen(sim, elapsedTime)
	shaman.Character.Advance(sim, elapsedTime)
}

// TryDropTotems will check to see if totems need to be re-cast.
//  If they do time.Duration will be returned will be >0.
func (shaman *Shaman) TryDropTotems(sim *core.Simulation) time.Duration {
	if sim.CurrentTime < shaman.NextTotemDrop {
		return 0
	}

	// 110 = totem of wrath
	// 90 = mana stream
	// 240 = wrath of air
	// currently hardcoded to include 25% cost reduction from resto talents
	var manaCost float64 = 110 + 90 + 240

	cast := &core.SingleTargetDirectDamageSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				Name:         "Shaman Totems",
				ActionID:     core.ActionID{SpellID: 36936}, // just using totem of wrath
				Character:    shaman.GetCharacter(),
				BaseManaCost: manaCost,
				ManaCost:     manaCost,
				CastTime:     time.Second * 3,
				SpellSchool:  stats.NatureSpellPower,
			},
		},
		Effect: core.DirectDamageSpellEffect{
			SpellEffect: core.SpellEffect{
				Target:               sim.GetPrimaryTarget(),
				BonusSpellCritRating: -100000, // dont let this spell hit/crit to accidentally proc shaman abilities
				BonusSpellHitRating:  -100000,
			},
		},
	}
	cast.Act(sim)

	shaman.NextTotemDrop = sim.CurrentTime + time.Second*120
	return cast.CastTime
}

var ElementalMasteryAuraID = core.NewAuraID()
var ElementalMasteryCooldownID = core.NewCooldownID()

func (shaman *Shaman) registerElementalMasteryCD() {
	if !shaman.Talents.ElementalMastery {
		return
	}

	shaman.AddMajorCooldown(core.MajorCooldown{
		CooldownID: ElementalMasteryCooldownID,
		Cooldown:   time.Minute * 3,
		ActivationFactory: func(sim *core.Simulation) core.CooldownActivation {
			return func(sim *core.Simulation, character *core.Character) bool {
				character.AddAura(sim, core.Aura{
					ID:      ElementalMasteryAuraID,
					Name:    "Elemental Mastery",
					Expires: core.NeverExpires,
					OnCast: func(sim *core.Simulation, cast *core.Cast) {
						cast.ManaCost = 0
						cast.GuaranteedCrit = true
					},
					OnCastComplete: func(sim *core.Simulation, cast *core.Cast) {
						// Remove the buff and put skill on CD
						character.SetCD(ElementalMasteryCooldownID, sim.CurrentTime+time.Minute*3)
						character.RemoveAura(sim, ElementalMasteryAuraID)
						character.UpdateMajorCooldowns(sim)
					},
				})
				return true
			}
		},
	})
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Strength:  103,
		stats.Agility:   61,
		stats.Stamina:   113,
		stats.Intellect: 109,
		stats.Spirit:    122,
		stats.Mana:      2678,
		stats.SpellCrit: 47.89,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Intellect: 104,
		stats.Mana:      2678,
		stats.Spirit:    135,
		stats.SpellCrit: 47.89,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassShaman}] = stats.Stats{
		stats.Intellect: 104,
		stats.Mana:      2678,
		stats.Spirit:    135,
		stats.SpellCrit: 47.89,
	}

	trollStats := stats.Stats{
		stats.Intellect: 104,
		stats.Mana:      2678,
		stats.Spirit:    135,
		stats.SpellCrit: 47.89,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll10, Class: proto.Class_ClassShaman}] = trollStats
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll30, Class: proto.Class_ClassShaman}] = trollStats
}

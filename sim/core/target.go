package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type Encounter struct {
	Duration float64
	Targets  []*Target
}

func NewEncounter(options proto.Encounter) Encounter {
	encounter := Encounter{
		Duration: options.Duration,
		Targets:  []*Target{},
	}

	for targetIndex, targetOptions := range options.Targets {
		target := NewTarget(*targetOptions, int32(targetIndex))
		encounter.Targets = append(encounter.Targets, target)
	}

	encounter.finalize()

	return encounter
}

func (encounter *Encounter) finalize() {
	for _, target := range encounter.Targets {
		target.finalize()
	}
}

func (encounter *Encounter) doneIteration(simDuration time.Duration) {
	for i, _ := range encounter.Targets {
		target := encounter.Targets[i]
		target.doneIteration(simDuration)
	}
}

func (encounter *Encounter) GetMetricsProto(numIterations int32) *proto.EncounterMetrics {
	metrics := &proto.EncounterMetrics{
		Targets: make([]*proto.TargetMetrics, len(encounter.Targets)),
	}

	i := 0
	for _, target := range encounter.Targets {
		metrics.Targets[i] = target.GetMetricsProto(numIterations)
		i++
	}

	return metrics
}

// Target is an enemy that can be the target of attacks/spells.
// Currently targets are basically just lvl 73 target dummies.
type Target struct {
	// Index of this target among all the targets. Primary target has index 0,
	// 2nd target has index 1, etc.
	Index int32

	armor int32

	MobType proto.MobType

	// Provides aura tracking behavior. Targets need auras to handle debuffs.
	auraTracker

	// Whether finalize() has been called yet for this Character.
	// All fields above this may not be altered once finalized is set.
	finalized bool
}

func NewTarget(options proto.Target, targetIndex int32) *Target {
	target := &Target{
		Index:       targetIndex,
		armor:       options.Armor,
		MobType:     options.MobType,
		auraTracker: newAuraTracker(true),
	}
	// TODO: Do something with this
	target.auraTracker.playerID = -1

	if options.Debuffs != nil {
		applyDebuffEffects(target, *options.Debuffs)
	}

	return target
}

func applyDebuffEffects(target *Target, debuffs proto.Debuffs) {
	if debuffs.Misery {
		target.AddPermanentAura(func(sim *Simulation) Aura {
			return miseryAura()
		})
	}

	if debuffs.JudgementOfWisdom {
		target.AddPermanentAura(func(sim *Simulation) Aura {
			return judgementOfWisdomAura()
		})
	}

	if debuffs.ImprovedSealOfTheCrusader {
		target.AddPermanentAura(func(sim *Simulation) Aura {
			return improvedSealOfTheCrusaderAura()
		})
	}

	if debuffs.CurseOfElements != proto.TristateEffect_TristateEffectMissing {
		target.AddPermanentAura(func(sim *Simulation) Aura {
			return curseOfElementsAura(debuffs.CurseOfElements)
		})
	}

	if debuffs.IsbUptime > 0.0 {
		target.AddPermanentAura(func(sim *Simulation) Aura {
			return ImprovedShadowBoltAura(debuffs.IsbUptime)
		})
	}
}

var MiseryDebuffID = NewDebuffID()

func miseryAura() Aura {
	return Aura{
		ID:   MiseryDebuffID,
		Name: "Misery",
		OnBeforeSpellHit: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect) {
			spellEffect.DamageMultiplier *= 1.05
		},
		OnPeriodicDamage: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect, tickDamage *float64) {
			*tickDamage *= 1.05
		},
	}
}

var JudgementOfWisdomDebuffID = NewDebuffID()

func judgementOfWisdomAura() Aura {
	const mana = 74 / 2 // 50% proc
	return Aura{
		ID:   JudgementOfWisdomDebuffID,
		Name: "Judgement of Wisdom",
		OnSpellHit: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect) {
			if spellCast.ActionID.ItemID == ItemIDTheLightningCapacitor {
				return // TLC cant proc JoW
			}

			character := spellCast.Character
			// Only apply to agents that have mana.
			if character.MaxMana() > 0 {
				character.AddStat(stats.Mana, mana)
				if sim.Log != nil {
					sim.Log("(%d) +Judgement Of Wisdom: 37 mana (74 @ 50%% proc)\n", character.ID)
				}
			}
		},
	}
}

var ImprovedSealOfTheCrusaderDebuffID = NewDebuffID()

func improvedSealOfTheCrusaderAura() Aura {
	return Aura{
		ID:   ImprovedSealOfTheCrusaderDebuffID,
		Name: "Improved Seal of the Crusader",
		OnBeforeSpellHit: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect) {
			spellEffect.BonusSpellCritRating += 3 * SpellCritRatingPerCritChance
			// FUTURE: melee crit bonus, research actual value
		},
	}
}

var CurseOfElementsDebuffID = NewDebuffID()

func curseOfElementsAura(coe proto.TristateEffect) Aura {
	mult := 1.1
	if coe == proto.TristateEffect_TristateEffectImproved {
		mult = 1.13
	}
	return Aura{
		ID:   CurseOfElementsDebuffID,
		Name: "Curse of the Elements",
		OnBeforeSpellHit: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect) {
			if spellCast.SpellSchool == stats.NatureSpellPower ||
				spellCast.SpellSchool == stats.HolySpellPower ||
				spellCast.SpellSchool == stats.AttackPower {
				return // does not apply to these schools
			}
			spellEffect.DamageMultiplier *= mult
		},
		OnPeriodicDamage: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect, tickDamage *float64) {
			if spellCast.SpellSchool == stats.NatureSpellPower ||
				spellCast.SpellSchool == stats.HolySpellPower ||
				spellCast.SpellSchool == stats.AttackPower {
				return // does not apply to these schools
			}
			*tickDamage *= mult
		},
	}
}

var ImprovedShadowBoltID = NewDebuffID()

func ImprovedShadowBoltAura(uptime float64) Aura {
	mult := (1 + uptime*0.2)
	return Aura{
		ID:   ImprovedShadowBoltID,
		Name: "Improved Shadow Bolt",
		OnBeforeSpellHit: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect) {
			if spellCast.SpellSchool != stats.ShadowSpellPower {
				return // does not apply to these schools
			}
			spellEffect.DamageMultiplier *= mult
		},
		OnPeriodicDamage: func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect, tickDamage *float64) {
			if spellCast.SpellSchool != stats.ShadowSpellPower {
				return // does not apply to these schools
			}
			*tickDamage *= mult
		},
	}
}

func (target *Target) finalize() {
	if target.finalized {
		return
	}
	target.finalized = true

	target.auraTracker.finalize()
}

func (target *Target) Reset(sim *Simulation) {
	target.auraTracker.reset(sim)
}

func (target *Target) Advance(sim *Simulation, elapsedTime time.Duration) {
	target.auraTracker.advance(sim)
}

func (target *Target) doneIteration(simDuration time.Duration) {
	target.auraTracker.doneIteration(simDuration)
}

func (target *Target) GetMetricsProto(numIterations int32) *proto.TargetMetrics {
	return &proto.TargetMetrics{
		Auras: target.auraTracker.GetMetricsProto(numIterations),
	}
}
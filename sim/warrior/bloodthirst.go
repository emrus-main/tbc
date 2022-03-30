package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var BloodthirstCooldownID = core.NewCooldownID()
var BloodthirstActionID = core.ActionID{SpellID: 30335, CooldownID: BloodthirstCooldownID}

func (warrior *Warrior) newBloodthirstTemplate(_ *core.Simulation) core.SimpleSpellTemplate {
	warrior.bloodthirstCost = 30
	if ItemSetDestroyerBattlegear.CharacterHasSetBonus(&warrior.Character, 4) {
		warrior.bloodthirstCost -= 5
	}

	ability := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				ActionID:    BloodthirstActionID,
				Character:   &warrior.Character,
				SpellSchool: core.SpellSchoolPhysical,
				GCD:         core.GCDDefault,
				Cooldown:    time.Second * 6,
				IgnoreHaste: true,
				BaseCost: core.ResourceCost{
					Type:  stats.Rage,
					Value: warrior.bloodthirstCost,
				},
				Cost: core.ResourceCost{
					Type:  stats.Rage,
					Value: warrior.bloodthirstCost,
				},
			},
			OutcomeRollCategory: core.OutcomeRollCategorySpecial,
			CritRollCategory:    core.CritRollCategoryPhysical,
			CritMultiplier:      warrior.critMultiplier(true),
		},
		Effect: core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spellCast *core.SpellCast) float64 {
					return hitEffect.MeleeAttackPower(spellCast) * 0.45
				},
				TargetSpellCoefficient: 0, // Doesn't scale with +damage on target?
			},
		},
	}

	if ItemSetOnslaughtBattlegear.CharacterHasSetBonus(&warrior.Character, 4) {
		ability.Effect.DamageMultiplier *= 1.05
	}

	refundAmount := warrior.bloodthirstCost * 0.8
	ability.Effect.OnSpellHit = func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellEffect) {
		if !spellEffect.Landed() {
			warrior.AddRage(sim, refundAmount, core.ActionID{OtherID: proto.OtherAction_OtherActionRefund})
		}
	}

	return core.NewSimpleSpellTemplate(ability)
}

func (warrior *Warrior) NewBloodthirst(_ *core.Simulation, target *core.Target) *core.SimpleSpell {
	bt := &warrior.bloodthirst
	warrior.bloodthirstTemplate.Apply(bt)

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	bt.Effect.Target = target
	if warrior.StanceMatches(DefensiveStance) {
		bt.Effect.ThreatMultiplier *= 1 + 0.21*float64(warrior.Talents.TacticalMastery)
	}

	return bt
}

func (warrior *Warrior) CanBloodthirst(sim *core.Simulation) bool {
	return warrior.Talents.Bloodthirst && warrior.CurrentRage() >= warrior.bloodthirstCost && !warrior.IsOnCD(BloodthirstCooldownID, sim.CurrentTime)
}

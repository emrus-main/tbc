package hunter

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var AimedShotCooldownID = core.NewCooldownID()
var AimedShotActionID = core.ActionID{SpellID: 27065, CooldownID: AimedShotCooldownID}

func (hunter *Hunter) registerAimedShotSpell(sim *core.Simulation) {
	baseCost := 370.0

	hunter.AimedShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    AimedShotActionID,
		SpellSchool: core.SpellSchoolPhysical,
		SpellExtras: core.SpellExtrasMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.NewCast{
				Cost: baseCost * (1 - 0.02*float64(hunter.Talents.Efficiency)),
				// Actual aimed shot has a 2.5s cast time, but we only use it as an instant precast.
				//CastTime:       time.Millisecond * 2500,
				//Cooldown:       time.Second * 6,
				//GCD:            core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskRangedSpecial,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: hunter.talonOfAlarDamageMod(core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (hitEffect.RangedAttackPower(spell.Character)+hitEffect.RangedAttackPowerOnTarget())*0.2 +
						hunter.AutoAttacks.Ranged.BaseDamage(sim) +
						hunter.AmmoDamageBonus +
						hitEffect.BonusWeaponDamage(spell.Character) +
						870
				},
				TargetSpellCoefficient: 1,
			}),
			OutcomeApplier: core.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, sim.GetPrimaryTarget())),
		}),
	})
}

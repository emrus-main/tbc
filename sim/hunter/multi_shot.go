package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var MultiShotCooldownID = core.NewCooldownID()
var MultiShotActionID = core.ActionID{SpellID: 27021, CooldownID: MultiShotCooldownID}

func (hunter *Hunter) registerMultiShotSpell(sim *core.Simulation) {
	baseCost := 275.0

	baseEffect := core.SpellEffect{
		ProcMask: core.ProcMaskRangedSpecial,

		BonusCritRating:  float64(hunter.Talents.ImprovedBarrage) * 4 * core.MeleeCritRatingPerCritChance,
		DamageMultiplier: 1 + 0.04*float64(hunter.Talents.Barrage),
		ThreatMultiplier: 1,

		BaseDamage: hunter.talonOfAlarDamageMod(core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return (hitEffect.RangedAttackPower(spell.Character)+hitEffect.RangedAttackPowerOnTarget())*0.2 +
					hunter.AutoAttacks.Ranged.BaseDamage(sim) +
					hunter.AmmoDamageBonus +
					hitEffect.BonusWeaponDamage(spell.Character) +
					205
			},
			TargetSpellCoefficient: 1,
		}),
		OutcomeApplier: core.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, sim.GetPrimaryTarget())),

		OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			hunter.rotation(sim, false)
		},
	}

	numHits := core.MinInt32(3, sim.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = sim.GetTarget(i)
	}

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    MultiShotActionID,
		SpellSchool: core.SpellSchoolPhysical,
		SpellExtras: core.SpellExtrasMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.NewCast{
				Cost: baseCost *
					(1 - 0.02*float64(hunter.Talents.Efficiency)) *
					core.TernaryFloat64(ItemSetDemonStalker.CharacterHasSetBonus(&hunter.Character, 4), 0.9, 1),

				GCD: core.GCDDefault + hunter.latency,
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.NewCast) {
				cast.CastTime = hunter.MultiShotCastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			Cooldown:    time.Second * 10,
		},

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (hunter *Hunter) MultiShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*500)/hunter.RangedSwingSpeed()) + hunter.latency
}

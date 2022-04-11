package hunter

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ScorpidStingActionID = core.ActionID{SpellID: 3043}

func (hunter *Hunter) registerScorpidStingSpell(sim *core.Simulation) {
	hunter.ScorpidStingAura = core.ScorpidStingAura(sim.GetPrimaryTarget())

	baseCost := hunter.BaseMana() * 0.09

	hunter.ScorpidSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    ScorpidStingActionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.NewCast{
				Cost: baseCost * (1 - 0.02*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ThreatMultiplier: 1,
			OutcomeApplier:   core.OutcomeFuncRangedHit(),
			OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					hunter.ScorpidStingAura.Activate(sim)
				}
			},
		}),
	})
}

package rogue

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (rogue *Rogue) registerExposeArmorSpell() {
	baseCost := 25.0
	refundAmount := 0.4 * float64(rogue.Talents.QuickRecovery)

	rogue.ExposeArmorAura = core.ExposeArmorAura(rogue.Env.GetPrimaryTarget(), rogue.Talents.ImprovedExposeArmor)

	rogue.ExposeArmor = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26866, Tag: 5},
		SpellSchool: core.SpellSchoolPhysical,
		SpellExtras: core.SpellExtrasMeleeMetrics | rogue.finisherFlags(),

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			ModifyCast:  rogue.applyDeathmantle,
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   rogue.OutcomeFuncMeleeSpecialHit(),
			OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					rogue.ExposeArmorAura.Activate(sim)
					rogue.ApplyFinisher(sim, spell.ActionID)
					if sim.GetRemainingDuration() <= time.Second*30 {
						rogue.doneEA = true
					}
				} else {
					if refundAmount > 0 {
						rogue.AddEnergy(sim, spell.CurCast.Cost*refundAmount, core.ActionID{SpellID: 31245})
					}
				}
			},
		}),
	})
}

func (rogue *Rogue) MaintainingExpose(target *core.Target) bool {
	return !rogue.doneEA && (rogue.Talents.ImprovedExposeArmor == 2 || !target.HasActiveAura(core.SunderArmorAuraLabel))
}

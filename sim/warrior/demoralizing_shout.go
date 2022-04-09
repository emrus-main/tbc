package warrior

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var DemoralizingShoutActionID = core.ActionID{SpellID: 25203}

func (warrior *Warrior) registerDemoralizingShoutSpell(sim *core.Simulation) {
	warrior.shoutCost = 10.0
	if ItemSetBoldArmor.CharacterHasSetBonus(&warrior.Character, 2) {
		warrior.shoutCost -= 2
	}
	warrior.demoShoutCost = warrior.shoutCost - float64(warrior.Talents.FocusedRage)

	ability := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				ActionID:    DemoralizingShoutActionID,
				Character:   &warrior.Character,
				SpellSchool: core.SpellSchoolPhysical,
				GCD:         core.GCDDefault,
				IgnoreHaste: true,
				BaseCost: core.ResourceCost{
					Type:  stats.Rage,
					Value: warrior.demoShoutCost,
				},
				Cost: core.ResourceCost{
					Type:  stats.Rage,
					Value: warrior.demoShoutCost,
				},
			},
		},
	}

	baseEffect := core.SpellEffect{
		ThreatMultiplier: 1,
		FlatThreatBonus:  56,
		OutcomeApplier:   core.OutcomeFuncMagicHit(),
	}

	numHits := sim.GetNumTargets()
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = sim.GetTarget(i)

		demoShoutAura := core.DemoralizingShoutAura(effects[i].Target, warrior.Talents.BoomingVoice, warrior.Talents.ImprovedDemoralizingShout)
		if i == 0 {
			warrior.DemoralizingShoutAura = demoShoutAura
		}

		effects[i].OnSpellHit = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				demoShoutAura.Activate(sim)
			}
		}
	}

	warrior.DemoralizingShout = warrior.RegisterSpell(core.SpellConfig{
		Template:     ability,
		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (warrior *Warrior) CanDemoralizingShout(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.demoShoutCost
}

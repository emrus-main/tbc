package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var ShieldSlamCooldownID = core.NewCooldownID()
var ShieldSlamActionID = core.ActionID{SpellID: 30356, CooldownID: ShieldSlamCooldownID}

func (warrior *Warrior) newShieldSlamTemplate(_ *core.Simulation) core.SimpleSpellTemplate {
	warrior.shieldSlamCost = 20.0 - float64(warrior.Talents.FocusedRage)

	warrior.canShieldSlam = warrior.Talents.ShieldSlam && warrior.Equip[proto.ItemSlot_ItemSlotOffHand].WeaponType == proto.WeaponType_WeaponTypeShield

	ability := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				ActionID:            ShieldSlamActionID,
				Character:           &warrior.Character,
				OutcomeRollCategory: core.OutcomeRollCategorySpecial,
				CritRollCategory:    core.CritRollCategoryPhysical,
				SpellSchool:         core.SpellSchoolPhysical,
				GCD:                 core.GCDDefault,
				Cooldown:            time.Second * 6,
				IgnoreHaste:         true,
				BaseCost: core.ResourceCost{
					Type:  stats.Rage,
					Value: warrior.shieldSlamCost,
				},
				Cost: core.ResourceCost{
					Type:  stats.Rage,
					Value: warrior.shieldSlamCost,
				},
				CritMultiplier: warrior.critMultiplier(true),
			},
		},
		Effect: core.SpellHitEffect{
			SpellEffect: core.SpellEffect{
				ProcMask:               core.ProcMaskMeleeMHSpecial, // TODO: Is this right?
				DamageMultiplier:       1,
				StaticDamageMultiplier: 1,
				ThreatMultiplier:       1,
				FlatThreatBonus:        305,
			},
			DirectInput: core.DirectDamageInput{
				MinBaseDamage: 420,
				MaxBaseDamage: 440,
			},
		},
	}

	if ItemSetOnslaughtArmor.CharacterHasSetBonus(&warrior.Character, 4) {
		ability.Effect.SpellEffect.StaticDamageMultiplier *= 1.1
	}

	refundAmount := warrior.shieldSlamCost * 0.8
	ability.Effect.OnSpellHit = func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellEffect) {
		if !spellEffect.Landed() {
			warrior.AddRage(sim, refundAmount, core.ActionID{OtherID: proto.OtherAction_OtherActionRefund})
		}
	}

	return core.NewSimpleSpellTemplate(ability)
}

func (warrior *Warrior) NewShieldSlam(_ *core.Simulation, target *core.Target) *core.SimpleSpell {
	ss := &warrior.shieldSlam
	warrior.shieldSlamTemplate.Apply(ss)

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	ss.Effect.Target = target
	ss.Effect.DirectInput.FlatDamageBonus = warrior.GetStat(stats.BlockValue)

	return ss
}

func (warrior *Warrior) CanShieldSlam(sim *core.Simulation) bool {
	return warrior.canShieldSlam && warrior.CurrentRage() >= warrior.shieldSlamCost && !warrior.IsOnCD(ShieldSlamCooldownID, sim.CurrentTime)
}

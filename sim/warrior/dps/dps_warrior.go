package dps

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/warrior"
)

func RegisterDpsWarrior() {
	core.RegisterAgentFactory(
		proto.Player_Warrior{},
		proto.Spec_SpecWarrior,
		func(character core.Character, options proto.Player) core.Agent {
			return NewDpsWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Warrior)
			if !ok {
				panic("Invalid spec value for Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDpsWarrior(character core.Character, options proto.Player) *DpsWarrior {
	warOptions := options.GetWarrior()

	war := &DpsWarrior{
		Warrior: warrior.NewWarrior(character, *warOptions.Talents, warrior.WarriorInputs{
			ShoutType:            warOptions.Options.Shout,
			PrecastShout:         warOptions.Options.PrecastShout,
			PrecastShoutSapphire: warOptions.Options.PrecastShoutSapphire,
			PrecastShoutT2:       warOptions.Options.PrecastShoutT2,
			RampageCDThreshold:   core.DurationFromSeconds(warOptions.Rotation.RampageCdThreshold),
		}),
		Rotation: *warOptions.Rotation,
		Options:  *warOptions.Options,

		slamLatency: core.DurationFromSeconds(warOptions.Rotation.SlamLatency / 1000),
	}
	if war.Talents.ImprovedSlam != 2 {
		war.Rotation.UseSlam = false
	}

	war.EnableRageBar(warOptions.Options.StartingRage, core.TernaryFloat64(war.Talents.EndlessRage, 1.25, 1), func(sim *core.Simulation) {
		if war.GCD.IsReady(sim) {
			war.TryUseCooldowns(sim)
			if war.GCD.IsReady(sim) {
				war.tryQueueSlam(sim)
				war.doRotation(sim)
			}
		}
	})
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation) *core.Spell {
			if war.Rotation.UseCleave {
				if war.CurrentRage() < float64(war.Rotation.HsRageThreshold) {
					war.DequeueCleave(sim)
					return nil
				} else {
					return war.TryCleave(sim)
				}
			} else {
				if war.CurrentRage() < float64(war.Rotation.HsRageThreshold) {
					war.DequeueHeroicStrike(sim)
					return nil
				} else {
					return war.TryHeroicStrike(sim)
				}
			}
		},
	})

	return war
}

type DpsWarrior struct {
	*warrior.Warrior

	Options  proto.Warrior_Options
	Rotation proto.Warrior_Rotation

	castFirstSunder bool

	doSlamNext  bool
	castSlamAt  time.Duration
	slamLatency time.Duration
}

func (war *DpsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *DpsWarrior) Initialize() {
	war.Warrior.Initialize()

	if war.Options.UseRecklessness {
		war.RegisterRecklessnessCD()
	}

	war.DelayDPSCooldownsForArmorDebuffs()
}

func (war *DpsWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
	war.BerserkerStanceAura.Activate(sim)
	war.Stance = warrior.BerserkerStance

	war.doSlamNext = false
	war.castFirstSunder = false
	war.castSlamAt = 0
}

package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func init() {
	core.AddItemEffect(32485, ApplyAshtongueTalismanOfValor)

	core.AddItemSet(&ItemSetBoldArmor)
	core.AddItemSet(&ItemSetWarbringerArmor)
	core.AddItemSet(&ItemSetWarbringerBattlegear)
	core.AddItemSet(&ItemSetDestroyerArmor)
	core.AddItemSet(&ItemSetDestroyerBattlegear)
	core.AddItemSet(&ItemSetOnslaughtArmor)
	core.AddItemSet(&ItemSetOnslaughtBattlegear)
}

var ItemSetBoldArmor = core.ItemSet{
	Name: "Bold Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// All of your shout abilities cost 2 less rage.
			// Handled in demoralizing_shout.go.
		},
		4: func(agent core.Agent) {
			// Your Charge ability generates an additional 5 rage.
		},
	},
}

var WarbringerArmor4PcAuraID = core.NewAuraID()
var WarbringerArmor4PcProcAuraID = core.NewAuraID()

var ItemSetWarbringerArmor = core.ItemSet{
	Name: "Warbringer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// You have a chance each time you parry to gain Blade Turning, absorbing 200 damage for 15 sec.
		},
		4: func(agent core.Agent) {
			// Your Revenge ability causes your next damaging ability to do 10% more damage.
			character := agent.GetCharacter()

			procAura := core.Aura{
				ID:       WarbringerArmor4PcProcAuraID,
				ActionID: core.ActionID{SpellID: 37516},
				Duration: core.NeverExpires,
				OnBeforeSpellHit: func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellHitEffect) {
					spellEffect.DamageMultiplier *= 1.1
				},
				OnSpellHit: func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellEffect) {
					if spellEffect.Damage > 0 {
						character.RemoveAura(sim, WarbringerArmor4PcProcAuraID)
					}
				},
			}

			character.AddPermanentAura(func(sim *core.Simulation) core.Aura {
				return core.Aura{
					ID: WarbringerArmor4PcAuraID,
					OnSpellHit: func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellEffect) {
						if !spellCast.SameAction(RevengeActionID) {
							return
						}

						character.AddAura(sim, procAura)
					},
				}
			})
		},
	},
}

var ItemSetWarbringerBattlegear = core.ItemSet{
	Name: "Warbringer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your whirlwind ability costs 5 less rage.
			// Handled in whirlwind.go.
		},
		4: func(agent core.Agent) {
			// You gain an additional 2 rage each time one of your attacks is parried or dodged.
		},
	},
}

var ItemSetDestroyerArmor = core.ItemSet{
	Name: "Destroyer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Each time you use your Shield Block ability, you gain 100 block value against a single attack in the next 6 sec.
		},
		4: func(agent core.Agent) {
			// You have a chance each time you are hit to gain 200 haste rating for 10 sec.
		},
	},
}

var ItemSetDestroyerBattlegear = core.ItemSet{
	Name: "Destroyer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Overpower ability now grants you 100 attack power for 5 sec.
		},
		4: func(agent core.Agent) {
			// Your Bloodthirst and Mortal Strike abilities cost 5 less rage.
		},
	},
}

var ItemSetOnslaughtArmor = core.ItemSet{
	Name: "Onslaught Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the health bonus from your Commanding Shout ability by 170.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Shield Slam ability by 10%.
			// Handled in shield_slam.go.
		},
	},
}

var ItemSetOnslaughtBattlegear = core.ItemSet{
	Name: "Onslaught Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Reduces the rage cost of your Execute ability by 3.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Mortal Strike and Bloodthirst abilities by 5%.
			// Handled in bloodthirst.go and mortal_strike.go.
		},
	},
}

var AshtongueTalismanOfValorAuraID = core.NewAuraID()
var AshtongueTalismanOfValorProcAuraID = core.NewAuraID()

func ApplyAshtongueTalismanOfValor(agent core.Agent) {
	character := agent.GetCharacter()
	character.AddPermanentAura(func(sim *core.Simulation) core.Aura {
		applyStatAura := character.NewTemporaryStatsAuraApplier(AshtongueTalismanOfValorProcAuraID, core.ActionID{ItemID: 32485}, stats.Stats{stats.Strength: 55}, time.Second*12)

		return core.Aura{
			ID: AshtongueTalismanOfValorAuraID,
			OnSpellHit: func(sim *core.Simulation, spellCast *core.SpellCast, spellEffect *core.SpellEffect) {
				if !spellCast.SameAction(ShieldSlamActionID) && !spellCast.SameAction(BloodthirstActionID) {
					return
				}

				if sim.RandomFloat("AshtongueTalismanOfValor") < 0.25 {
					applyStatAura(sim)
				}
			},
		}
	})
}

package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Keep these (and their functions) in alphabetical order.
func init() {
	core.AddItemSet(ItemSetManaEtched)
	core.AddItemSet(ItemSetNetherstrike)
	core.AddItemSet(ItemSetSpellstrike)
	core.AddItemSet(ItemSetTheTwinStars)
	core.AddItemSet(ItemSetWindhawk)
}

var ItemSetManaEtched = core.ItemSet{
	Name:  "Mana Etched",
	Items: map[int32]struct{}{28193: {}, 27465: {}, 27907: {}, 27796: {}, 28191: {}},
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellHit, 35)
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			const spellBonus = 110.0
			const duration = time.Second * 15

			character.AddPermanentAura(core.Aura{
				ID:      core.MagicIDManaEtched,
				OnCastComplete: func(sim *core.Simulation, cast core.DirectCastAction) {
					if sim.Rando.Float64("unmarked") < 0.02 {
						core.AddAuraWithTemporaryStats(sim, character, core.MagicIDManaEtchedInsight, stats.SpellPower, spellBonus, duration)
					}
				},
			})
		},
	},
}

var ItemSetNetherstrike = core.ItemSet{
	Name:  "Netherstrike",
	Items: map[int32]struct{}{29519: {}, 29521: {}, 29520: {}},
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellPower, 23)
		},
	},
}

var ItemSetSpellstrike = core.ItemSet{
	Name:    "Spellstrike",
	Items:   map[int32]struct{}{24266: {}, 24262: {}},
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			const spellBonus = 92.0
			const duration = time.Second * 10

			character.AddPermanentAura(core.Aura{
				ID:      core.MagicIDSpellstrike,
				OnCastComplete: func(sim *core.Simulation, cast core.DirectCastAction) {
					if sim.Rando.Float64("spellstrike") < 0.05 {
						core.AddAuraWithTemporaryStats(sim, character, core.MagicIDSpellstrikeInfusion, stats.SpellPower, spellBonus, duration)
					}
				},
			})
		},
	},
}

var ItemSetTheTwinStars = core.ItemSet{
	Name:  "The Twin Stars",
	Items: map[int32]struct{}{31338: {}, 31339: {}},
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.SpellPower, 15)
		},
	},
}

var ItemSetWindhawk = core.ItemSet{
	Name:  "Windhawk",
	Items: map[int32]struct{}{29524: {}, 29523: {}, 29522: {}},
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MP5, 8)
		},
	},
}

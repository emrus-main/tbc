package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/items"
)

const SpellIDLB12 int32 = 25449

func (shaman *Shaman) newLightningBoltTemplate(sim *core.Simulation, isLightningOverload bool) core.DirectCastAction {
	baseManaCost := 300.0
	if shaman.Equip[items.ItemSlotRanged].ID == TotemOfThePulsingEarth {
		baseManaCost -= 27.0
	}

	spellTemplate := shaman.newElectricSpellTemplate(
		"Lightning Bolt",
		core.ActionID{
			SpellID: SpellIDLB12,
		},
		baseManaCost,
		time.Millisecond*2500,
		isLightningOverload)

	hitInput := core.DirectCastDamageInput{
		MinBaseDamage: 571,
		MaxBaseDamage: 652,
		SpellCoefficient: 0.794,
		DamageMultiplier: 1,
	}
	shaman.applyElectricSpellHitInputModifiers(&hitInput, isLightningOverload)

	spellTemplate.HitInputs = []core.DirectCastDamageInput{hitInput}

	if !isLightningOverload && shaman.Talents.LightningOverload > 0 {
		lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.04
		spellTemplate.OnSpellHit = func(sim *core.Simulation, cast *core.Cast, result *core.DirectCastDamageResult) {
			if shaman.Talents.ElementalFocus && result.Crit {
				shaman.ElementalFocusStacks = 2
			}

			if sim.RandomFloat("LB Lightning Overload") < lightningOverloadChance {
				overloadAction := shaman.NewLightningBolt(sim, result.Target, true)
				overloadAction.Act(sim)
			}
		}
	} else {
		spellTemplate.OnSpellHit = func(sim *core.Simulation, cast *core.Cast, result *core.DirectCastDamageResult) {
			if shaman.Talents.ElementalFocus && result.Crit {
				shaman.ElementalFocusStacks = 2
			}
		}
	}

	return spellTemplate
}

func (shaman *Shaman) NewLightningBolt(sim *core.Simulation, target *core.Target, isLightningOverload bool) *core.DirectCastAction {
	var spell *core.DirectCastAction

	if isLightningOverload {
		spell = &shaman.electricSpellLO
		*spell = shaman.lightningBoltLOTemplate
	} else {
		spell = &shaman.electricSpell
		*spell = shaman.lightningBoltTemplate
	}
	spell.HitInputs[0].Target = target

	shaman.applyElectricSpellInitModifiers(spell)
	spell.Init(sim)

	return spell
}

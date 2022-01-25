package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// TODO: magma and nova totems need to apply to all targets probably instead of just the primary target.

const SpellIDSearingTotem int32 = 25533

func (shaman *Shaman) newSearingTotemTemplate(sim *core.Simulation) core.SimpleSpellTemplate {
	spell := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				CritMultiplier: 1.5,
				SpellSchool:    stats.FireSpellPower,
				Character:      &shaman.Character,
				BaseManaCost:   205,
				ManaCost:       205,
				ActionID: core.ActionID{
					SpellID: SpellIDSearingTotem,
				},
				GCDCooldown: time.Second,
			},
		},
		Effect: core.SpellHitEffect{
			SpellEffect: core.SpellEffect{
				DamageMultiplier:       1,
				StaticDamageMultiplier: 1,
				IgnoreHitCheck:         true,
			},
			DotInput: core.DotDamageInput{
				// These are the real tick values, but searing totem doesn't start its next
				// cast until the previous missile hits the target. We don't have an option
				// for target distance yet so just pretend the tick rate is lower.
				//NumberOfTicks:        30,
				//TickLength:           time.Second * 2,
				NumberOfTicks: 24,
				TickLength:    time.Second * 60 / 24,

				TickBaseDamage:       58,
				TickSpellCoefficient: 0.167,
				TicksCanMissAndCrit:  true,
			},
		},
	}
	spell.Effect.DamageMultiplier *= 1 + float64(shaman.Talents.CallOfFlame)*0.05
	spell.ManaCost -= spell.BaseManaCost * float64(shaman.Talents.TotemicFocus) * 0.05
	spell.ManaCost -= spell.BaseManaCost * float64(shaman.Talents.MentalQuickness) * 0.02

	spell.OnCastComplete = func(sim *core.Simulation, cast *core.Cast) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*60
		shaman.tryTwistFireNova(sim)
	}

	return core.NewSimpleSpellTemplate(spell)
}

func (shaman *Shaman) NewSearingTotem(sim *core.Simulation, target *core.Target) *core.SimpleSpell {
	// Initialize cast from precomputed template.
	searingTotem := &shaman.FireTotemSpell
	shaman.searingTotemTemplate.Apply(searingTotem)

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	searingTotem.Effect.Target = target
	searingTotem.Init(sim)

	return searingTotem
}

const SpellIDMagmaTotem int32 = 25552

// This is probably not worth simming since no other spell in the game does this and AM isn't
// even a popular choice for arcane mages.
func (shaman *Shaman) newMagmaTotemTemplate(sim *core.Simulation) core.SimpleSpellTemplate {
	spell := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				CritMultiplier: 1.5,
				SpellSchool:    stats.FireSpellPower,
				Character:      &shaman.Character,
				BaseManaCost:   800,
				ManaCost:       800,
				ActionID: core.ActionID{
					SpellID: SpellIDMagmaTotem,
				},
				GCDCooldown: time.Second,
			},
		},
		AOECap: 1600,
	}
	spell.ManaCost -= spell.BaseManaCost * float64(shaman.Talents.TotemicFocus) * 0.05
	spell.ManaCost -= spell.BaseManaCost * float64(shaman.Talents.MentalQuickness) * 0.02

	baseEffect := core.SpellHitEffect{
		SpellEffect: core.SpellEffect{
			DamageMultiplier:       1,
			StaticDamageMultiplier: 1,
			IgnoreHitCheck:         true,
		},
		DotInput: core.DotDamageInput{
			NumberOfTicks:        10,
			TickLength:           time.Second * 2,
			TickBaseDamage:       97,
			TickSpellCoefficient: 0.067,
			TicksCanMissAndCrit:  true,
		},
	}
	baseEffect.StaticDamageMultiplier *= 1 + float64(shaman.Talents.CallOfFlame)*0.05

	spell.OnCastComplete = func(sim *core.Simulation, cast *core.Cast) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*20
		shaman.tryTwistFireNova(sim)
	}

	numHits := sim.GetNumTargets()
	effects := make([]core.SpellHitEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = sim.GetTarget(i)
	}
	spell.Effects = effects

	return core.NewSimpleSpellTemplate(spell)
}

func (shaman *Shaman) NewMagmaTotem(sim *core.Simulation) *core.SimpleSpell {
	// Initialize cast from precomputed template.
	magmaTotem := &shaman.FireTotemSpell
	shaman.magmaTotemTemplate.Apply(magmaTotem)

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	magmaTotem.Init(sim)

	return magmaTotem
}

const SpellIDNovaTotem int32 = 25537

var CooldownIDNovaTotem = core.NewCooldownID()

// This is probably not worth simming since no other spell in the game does this and AM isn't
// even a popular choice for arcane mages.
func (shaman *Shaman) newNovaTotemTemplate(sim *core.Simulation) core.SimpleSpellTemplate {
	spell := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				CritMultiplier: 1.5,
				SpellSchool:    stats.FireSpellPower,
				Character:      &shaman.Character,
				BaseManaCost:   800,
				ManaCost:       800,
				ActionID: core.ActionID{
					SpellID:    SpellIDNovaTotem,
					CooldownID: CooldownIDNovaTotem,
				},
				Cooldown:    time.Second * 15,
				GCDCooldown: time.Second,
			},
		},
		AOECap: 9975,
	}
	spell.ManaCost -= spell.BaseManaCost * float64(shaman.Talents.TotemicFocus) * 0.05
	spell.ManaCost -= spell.BaseManaCost * float64(shaman.Talents.MentalQuickness) * 0.02

	baseEffect := core.SpellHitEffect{
		SpellEffect: core.SpellEffect{
			DamageMultiplier:       1,
			StaticDamageMultiplier: 1,
			IgnoreHitCheck:         true,
		},
		DotInput: core.DotDamageInput{
			NumberOfTicks:        1,
			TickLength:           time.Second * 5,
			TickBaseDamage:       692,
			TickSpellCoefficient: 0.214,
			TicksCanMissAndCrit:  true,
		},
	}
	baseEffect.StaticDamageMultiplier *= 1 + float64(shaman.Talents.CallOfFlame)*0.05
	baseEffect.DotInput.TickLength -= time.Duration(shaman.Talents.ImprovedFireTotems) * time.Second

	tickLength := baseEffect.DotInput.TickLength
	spell.OnCastComplete = func(sim *core.Simulation, cast *core.Cast) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + tickLength
		shaman.tryTwistFireNova(sim)
	}

	numHits := sim.GetNumTargets()
	effects := make([]core.SpellHitEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = sim.GetTarget(i)
	}
	spell.Effects = effects

	return core.NewSimpleSpellTemplate(spell)
}

func (shaman *Shaman) NewNovaTotem(sim *core.Simulation) *core.SimpleSpell {
	// If we drop nova while another totem is running, cancel it.
	if shaman.FireTotemSpell.IsInUse() {
		shaman.FireTotemSpell.Cancel(sim)
	}

	// Initialize cast from precomputed template.
	novaTotem := &shaman.FireTotemSpell
	shaman.novaTotemTemplate.Apply(novaTotem)

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	novaTotem.Init(sim)

	return novaTotem
}
package warrior

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

type Warrior struct {
	core.Character

	Talents proto.WarriorTalents

	// Current state
	Stance             Stance
	heroicStrikeQueued bool
	revengeTriggered   bool

	// Cached values
	heroicStrikeCost float64
	canShieldSlam    bool

	castBattleStance    func(*core.Simulation)
	castDefensiveStance func(*core.Simulation)
	castBerserkerStance func(*core.Simulation)

	bloodthirstTemplate core.SimpleSpellTemplate
	bloodthirst         core.SimpleSpell

	demoralizingShoutTemplate core.SimpleSpellTemplate
	demoralizingShout         core.SimpleSpell

	devastateTemplate core.SimpleSpellTemplate
	devastate         core.SimpleSpell

	heroicStrikeTemplate core.SimpleSpellTemplate
	heroicStrike         core.SimpleSpell

	revengeTemplate core.SimpleSpellTemplate
	revenge         core.SimpleSpell

	shieldSlamTemplate core.SimpleSpellTemplate
	shieldSlam         core.SimpleSpell

	sunderArmorTemplate core.SimpleSpellTemplate
	sunderArmor         core.SimpleSpell

	thunderClapTemplate core.SimpleSpellTemplate
	thunderClap         core.SimpleSpell

	whirlwindTemplate core.SimpleSpellTemplate
	whirlwind         core.SimpleSpell
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (warrior *Warrior) Init(sim *core.Simulation) {
	warrior.castBattleStance = warrior.makeCastStance(sim, BattleStance, warrior.BattleStanceAura())
	warrior.castDefensiveStance = warrior.makeCastStance(sim, DefensiveStance, warrior.DefensiveStanceAura())
	warrior.castBerserkerStance = warrior.makeCastStance(sim, BerserkerStance, warrior.BerserkerStanceAura())

	warrior.bloodthirstTemplate = warrior.newBloodthirstTemplate(sim)
	warrior.demoralizingShoutTemplate = warrior.newDemoralizingShoutTemplate(sim)
	warrior.devastateTemplate = warrior.newDevastateTemplate(sim)
	warrior.heroicStrikeTemplate = warrior.newHeroicStrikeTemplate(sim)
	warrior.revengeTemplate = warrior.newRevengeTemplate(sim)
	warrior.shieldSlamTemplate = warrior.newShieldSlamTemplate(sim)
	warrior.sunderArmorTemplate = warrior.newSunderArmorTemplate(sim)
	warrior.thunderClapTemplate = warrior.newThunderClapTemplate(sim)
	warrior.whirlwindTemplate = warrior.newWhirlwindTemplate(sim)
}

func (warrior *Warrior) Reset(sim *core.Simulation) {
	warrior.heroicStrikeQueued = false
	warrior.revengeTriggered = false
}

func NewWarrior(character core.Character, talents proto.WarriorTalents) *Warrior {
	warrior := &Warrior{
		Character: character,
		Talents:   talents,
	}

	warrior.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.MeleeCrit,
		Modifier: func(agility float64, meleecrit float64) float64 {
			return meleecrit + (agility/33)*core.MeleeCritRatingPerCritChance
		},
	})
	warrior.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Agility,
		ModifiedStat: stats.Dodge,
		Modifier: func(agility float64, dodge float64) float64 {
			return dodge + (agility/30)*core.DodgeRatingPerDodgeChance
		},
	})
	warrior.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.AttackPower,
		Modifier: func(strength float64, attackPower float64) float64 {
			return attackPower + strength*2
		},
	})
	warrior.Character.AddStatDependency(stats.StatDependency{
		SourceStat:   stats.Strength,
		ModifiedStat: stats.BlockValue,
		Modifier: func(strength float64, blockValue float64) float64 {
			return blockValue + strength/20
		},
	})

	return warrior
}

func (warrior *Warrior) secondaryCritModifier(applyImpale bool) float64 {
	secondaryModifier := 0.0
	if applyImpale {
		secondaryModifier += 0.1 * float64(warrior.Talents.Impale)
	}
	return secondaryModifier
}
func (warrior *Warrior) critMultiplier(applyImpale bool) float64 {
	return warrior.MeleeCritMultiplier(1.0, warrior.secondaryCritModifier(applyImpale))
}
func (warrior *Warrior) spellCritMultiplier(applyImpale bool) float64 {
	return warrior.SpellCritMultiplier(1.0, warrior.secondaryCritModifier(applyImpale))
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    146,
		stats.Agility:     93,
		stats.Stamina:     132,
		stats.Intellect:   34,
		stats.Spirit:      53,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDwarf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    147,
		stats.Agility:     92,
		stats.Stamina:     136,
		stats.Intellect:   32,
		stats.Spirit:      50,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    140,
		stats.Agility:     99,
		stats.Stamina:     132,
		stats.Intellect:   38,
		stats.Spirit:      51,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    145,
		stats.Agility:     96,
		stats.Stamina:     133,
		stats.Intellect:   33,
		stats.Spirit:      56,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceNightElf, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    142,
		stats.Agility:     101,
		stats.Stamina:     132,
		stats.Intellect:   33,
		stats.Spirit:      51,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceOrc, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    148,
		stats.Agility:     93,
		stats.Stamina:     135,
		stats.Intellect:   30,
		stats.Spirit:      54,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTauren, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    150,
		stats.Agility:     91,
		stats.Stamina:     135,
		stats.Intellect:   28,
		stats.Spirit:      53,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	trollStats := stats.Stats{
		stats.Health:      4264,
		stats.Strength:    146,
		stats.Agility:     98,
		stats.Stamina:     134,
		stats.Intellect:   29,
		stats.Spirit:      52,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll10, Class: proto.Class_ClassWarrior}] = trollStats
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll30, Class: proto.Class_ClassWarrior}] = trollStats
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassWarrior}] = stats.Stats{
		stats.Health:      4264,
		stats.Strength:    144,
		stats.Agility:     94,
		stats.Stamina:     134,
		stats.Intellect:   31,
		stats.Spirit:      56,
		stats.AttackPower: 190,
		stats.MeleeCrit:   1.14 * core.MeleeCritRatingPerCritChance,
		stats.Dodge:       0.75 * core.DodgeRatingPerDodgeChance,
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type Agent interface {
	GetWarrior() *Warrior
}
import { Sim } from '../sim';
import { Gear } from '../api/gear';
import { raceNames } from '../api/names';
import { Buffs } from '../api/newapi';
import { Class } from '../api/newapi';
import { Consumes } from '../api/newapi';
import { Encounter } from '../api/newapi';
import { EquipmentSpec } from '../api/newapi';
import { Race } from '../api/newapi';
import { Spec } from '../api/newapi';
import { Stat } from '../api/newapi';
import { Stats } from '../api/stats';
import { specToEligibleRaces } from '../api/utils';
import { Actions } from '../components/actions';
import { CharacterStats } from '../components/character_stats';
import { CustomStatsPicker } from '../components/custom_stats_picker';
import { EnumPicker } from '../components/enum_picker';
import { GearPicker } from '../components/gear_picker';
import { IconInput } from '../components/icon_picker';
import { IconPicker } from '../components/icon_picker';
import { NumberPicker } from '../components/number_picker';
import { Results } from '../components/results';
import { SavedDataManager } from '../components/saved_data_manager';
import { newTalentsPicker } from '../talents/factory';

import { Theme, ThemeConfig } from './theme';

export interface DefaultThemeConfig<ClassType extends Class> extends ThemeConfig<ClassType> {
  displayStats: Array<Stat>;
  iconSections: Record<string, Array<IconInput>>;
  showTargetArmor: boolean;
  showNumTargets: boolean;
  presets: {
    gear: Array<{
      name: string,
      tooltip?: string,
      equipment: EquipmentSpec,
    }>;
    encounters: Array<{
      name: string,
      tooltip?: string,
      encounter: Encounter,
    }>;
    talents: Array<{
      name: string,
      tooltip?: string,
      talents: string,
    }>;
  },
}

export interface GearAndStats {
  gear: Gear,
  customStats: Stats,
}

export interface Settings {
  buffs: Buffs,
  consumes: Consumes,
  race: Race,
}

export class DefaultTheme<ClassType extends Class> extends Theme<ClassType> {
  private readonly _config: DefaultThemeConfig<ClassType>;

  constructor(parentElem: HTMLElement, config: DefaultThemeConfig<ClassType>) {
    super(parentElem, config)
    this._config = config;
    this.parentElem.innerHTML = layoutHTML;

    const results = new Results(this.parentElem.getElementsByClassName('default-results')[0] as HTMLElement);
    const actions = new Actions(this.parentElem.getElementsByClassName('default-actions')[0] as HTMLElement, this.sim, results, config.epStats, config.epReferenceStat);

    const characterStats = new CharacterStats(this.parentElem.getElementsByClassName('default-stats')[0] as HTMLElement, config.displayStats);

    const gearPicker = new GearPicker(this.parentElem.getElementsByClassName('gear-picker')[0] as HTMLElement, this.sim);
    const customStatsPicker = new CustomStatsPicker(this.parentElem.getElementsByClassName('custom-stats-picker')[0] as HTMLElement, this.sim, config.epStats);

    const talentsPicker = newTalentsPicker(config.spec, this.parentElem.getElementsByClassName('talents-picker')[0] as HTMLElement, this.sim);

    const settingsTab = document.getElementsByClassName('settings-inputs')[0] as HTMLElement;
    Object.keys(config.iconSections).forEach(pickerName => {
      const section = config.iconSections[pickerName];

      const sectionElem = document.createElement('section');
      sectionElem.classList.add('settings-section', pickerName + '-section');
      sectionElem.innerHTML = `<label>${pickerName}</label>`;
      settingsTab.appendChild(sectionElem);

      const iconPicker = new IconPicker(sectionElem, pickerName + '-icon-picker', this.sim, section, this);
    });

    const races = specToEligibleRaces[this.sim.spec];
    const racePicker = new EnumPicker(this.parentElem.getElementsByClassName('race-picker')[0] as HTMLElement, this.sim, {
      names: races.map(race => raceNames[race]),
      values: races,
      changedEvent: sim => sim.raceChangeEmitter,
      getValue: sim => sim.getRace(),
      setValue: (sim, newValue) => sim.setRace(newValue),
    });

    const encounterSectionElem = settingsTab.getElementsByClassName('encounter-section')[0] as HTMLElement;
    new NumberPicker(encounterSectionElem, this.sim, {
      label: 'Duration',
      changedEvent: (sim: Sim<any>) => sim.encounterChangeEmitter,
      getValue: (sim: Sim<any>) => sim.getEncounter().duration,
      setValue: (sim: Sim<any>, newValue: number) => {
        const encounter = sim.getEncounter();
        encounter.duration = newValue;
        sim.setEncounter(encounter);
      },
    });

    if (config.showTargetArmor) {
      new NumberPicker(encounterSectionElem, this.sim, {
        label: 'Target Armor',
        changedEvent: (sim: Sim<any>) => sim.encounterChangeEmitter,
        getValue: (sim: Sim<any>) => sim.getEncounter().targetArmor,
        setValue: (sim: Sim<any>, newValue: number) => {
          const encounter = sim.getEncounter();
          encounter.targetArmor = newValue;
          sim.setEncounter(encounter);
        },
      });
    } else {
    }

    if (config.showNumTargets) {
      new NumberPicker(encounterSectionElem, this.sim, {
        label: '# of Targets',
        changedEvent: (sim: Sim<any>) => sim.encounterChangeEmitter,
        getValue: (sim: Sim<any>) => sim.getEncounter().numTargets,
        setValue: (sim: Sim<any>, newValue: number) => {
          const encounter = sim.getEncounter();
          encounter.numTargets = newValue;
          sim.setEncounter(encounter);
        },
      });
    }
  }

  async init(): Promise<void> {
    const savedGearManager = new SavedDataManager<GearAndStats>(this.parentElem.getElementsByClassName('saved-gear-manager')[0] as HTMLElement, this.sim, {
      label: 'Gear',
      getData: (sim: Sim<any>) => {
        return {
          gear: sim.getGear(),
          customStats: sim.getCustomStats(),
        };
      },
      setData: (sim: Sim<any>, newGearAndStats: GearAndStats) => {
        sim.setGear(newGearAndStats.gear);
        sim.setCustomStats(newGearAndStats.customStats);
      },
      changeEmitters: [this.sim.gearChangeEmitter, this.sim.customStatsChangeEmitter],
      equals: (a: GearAndStats, b: GearAndStats) => a.gear.equals(b.gear) && a.customStats.equals(b.customStats),
      toJson: (a: GearAndStats) => {
        return {
          gear: EquipmentSpec.toJson(a.gear.asSpec()),
          customStats: a.customStats.toJson(),
        };
      },
      fromJson: (obj: any) => {
        return {
          gear: this.sim.lookupEquipmentSpec(EquipmentSpec.fromJson(obj['gear'])),
          customStats: Stats.fromJson(obj['customStats']),
        };
      },
    });

    const savedEncounterManager = new SavedDataManager<Encounter>(this.parentElem.getElementsByClassName('saved-encounter-manager')[0] as HTMLElement, this.sim, {
      label: 'Encounter',
      getData: (sim: Sim<any>) => sim.getEncounter(),
      setData: (sim: Sim<any>, newEncounter: Encounter) => sim.setEncounter(newEncounter),
      changeEmitters: [this.sim.encounterChangeEmitter],
      equals: (a: Encounter, b: Encounter) => Encounter.equals(a, b),
      toJson: (a: Encounter) => Encounter.toJson(a),
      fromJson: (obj: any) => Encounter.fromJson(obj),
    });

    const savedSettingsManager = new SavedDataManager<Settings>(this.parentElem.getElementsByClassName('saved-settings-manager')[0] as HTMLElement, this.sim, {
      label: 'Settings',
      getData: (sim: Sim<any>) => {
        return {
          buffs: sim.getBuffs(),
          consumes: sim.getConsumes(),
          race: sim.getRace(),
        };
      },
      setData: (sim: Sim<any>, newSettings: Settings) => {
        sim.setBuffs(newSettings.buffs);
        sim.setConsumes(newSettings.consumes);
        sim.setRace(newSettings.race);
      },
      changeEmitters: [this.sim.buffsChangeEmitter, this.sim.consumesChangeEmitter, this.sim.raceChangeEmitter],
      equals: (a: Settings, b: Settings) => Buffs.equals(a.buffs, b.buffs) && Consumes.equals(a.consumes, b.consumes) && a.race == b.race,
      toJson: (a: Settings) => {
        return {
          buffs: Buffs.toJson(a.buffs),
          consumes: Consumes.toJson(a.consumes),
          race: a.race,
        };
      },
      fromJson: (obj: any) => {
        return {
          buffs: Buffs.fromJson(obj['buffs']),
          consumes: Consumes.fromJson(obj['consumes']),
          race: Number(obj['race']),
        };
      },
    });

    const savedTalentsManager = new SavedDataManager<string>(this.parentElem.getElementsByClassName('saved-talents-manager')[0] as HTMLElement, this.sim, {
      label: 'Talents',
      getData: (sim: Sim<any>) => sim.getTalentsString(),
      setData: (sim: Sim<any>, newTalentsString: string) => sim.setTalentsString(newTalentsString),
      changeEmitters: [this.sim.talentsStringChangeEmitter],
      equals: (a: string, b: string) => a == b,
      toJson: (a: string) => a,
      fromJson: (obj: any) => obj,
    });

    await super.init();

    savedGearManager.loadUserData();
    this._config.presets.gear.forEach(gearConfig => {
      const gear = this.sim.lookupEquipmentSpec(gearConfig.equipment);
      savedGearManager.addSavedData(gearConfig.name, { gear: gear, customStats: new Stats(), }, true, gearConfig.tooltip);
    });

    savedEncounterManager.loadUserData();
    this._config.presets.encounters.forEach(encounterConfig => {
      savedEncounterManager.addSavedData(encounterConfig.name, encounterConfig.encounter, true, encounterConfig.tooltip);
    });

    savedSettingsManager.loadUserData();

    savedTalentsManager.loadUserData();
    this._config.presets.talents.forEach(talentsConfig => {
      savedTalentsManager.addSavedData(talentsConfig.name, talentsConfig.talents, true, talentsConfig.tooltip);
    });
  }
}

const layoutHTML = `
<div class="default-root">
  <section class="default-sidebar">
    <div class="default-title">
      TBC Elemental Shaman Sim
    </div>
    <div class="default-actions">
    </div>
    <div class="default-results">
    </div>
    <div class="default-stats">
    </div>
  </section>
  <section class="default-main">
    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" href="#gear-tab">Gear</a></li>
      <li><a data-toggle="tab" href="#settings-tab">Settings</a></li>
      <li><a data-toggle="tab" href="#talents-tab">Talents</a></li>
    </ul>
    <div class="tab-content">
      <div id="gear-tab" class="tab-pane fade in active">
        <div class="gear-tab">
          <div class="left-gear-panel">
            <div class="gear-picker">
            </div>
          </div>
          <div class="right-gear-panel">
            <div class="custom-stats-picker">
            </div>
            <div class="saved-gear-manager">
            </div>
          </div>
        </div>
      </div>
      <div id="settings-tab" class="tab-pane fade"">
        <div class="settings-tab">
          <div class="settings-inputs">
            <div class="settings-left-bar">
              <section class="settings-section encounter-section">
                <label>Encounter</label>
              </section>
              <section class="settings-section race-picker">
                <label>Race</label>
              </section>
            </div>
          </div>
          <div class="settings-bottom-bar">
            <div class="saved-encounter-manager">
            </div>
            <div class="saved-settings-manager">
            </div>
          </div>
        </div>
      </div>
      <div id="talents-tab" class="tab-pane fade"">
        <div class="talents-tab">
          <div class="talents-picker">
          </div>
          <div class="saved-talents-manager">
          </div>
        </div>
      </div>
    </div>
  </section>
</div>
`;

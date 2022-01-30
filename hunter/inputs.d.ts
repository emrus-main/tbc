import { Spec } from '/tbc/core/proto/common.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { Player } from '/tbc/core/player.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { Hunter_Options_Ammo as Ammo, Hunter_Options_QuiverBonus as QuiverBonus } from '/tbc/core/proto/hunter.js';
export declare const Quiver: {
    extraCssClasses: string[];
    numColumns: number;
    values: ({
        color: string;
        value: QuiverBonus;
        actionId?: undefined;
    } | {
        actionId: ActionId;
        value: QuiverBonus;
        color?: undefined;
    })[];
    equals: (a: QuiverBonus, b: QuiverBonus) => boolean;
    zeroValue: QuiverBonus;
    changedEvent: (player: Player<Spec.SpecHunter>) => TypedEvent<void>;
    getValue: (player: Player<Spec.SpecHunter>) => QuiverBonus;
    setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => void;
};
export declare const WeaponAmmo: {
    extraCssClasses: string[];
    numColumns: number;
    values: ({
        color: string;
        value: Ammo;
        actionId?: undefined;
    } | {
        actionId: ActionId;
        value: Ammo;
        color?: undefined;
    })[];
    equals: (a: Ammo, b: Ammo) => boolean;
    zeroValue: Ammo;
    changedEvent: (player: Player<Spec.SpecHunter>) => TypedEvent<void>;
    getValue: (player: Player<Spec.SpecHunter>) => Ammo;
    setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => void;
};
export declare const LatencyMs: {
    type: "number";
    getModObject: (simUI: IndividualSimUI<any>) => Player<any>;
    config: {
        extraCssClasses: string[];
        label: string;
        labelTooltip: string;
        changedEvent: (player: Player<Spec.SpecHunter>) => TypedEvent<void>;
        getValue: (player: Player<Spec.SpecHunter>) => number;
        setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => void;
    };
};
export declare const HunterRotationConfig: {
    inputs: ({
        type: "boolean";
        cssClass: string;
        getModObject: (simUI: IndividualSimUI<any>) => Player<any>;
        config: {
            label: string;
            labelTooltip: string;
            changedEvent: (player: Player<Spec.SpecHunter>) => TypedEvent<void>;
            getValue: (player: Player<Spec.SpecHunter>) => boolean;
            setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: boolean) => void;
        };
    } | {
        type: "number";
        cssClass: string;
        getModObject: (simUI: IndividualSimUI<any>) => Player<any>;
        config: {
            label: string;
            labelTooltip: string;
            changedEvent: (player: Player<Spec.SpecHunter>) => TypedEvent<void>;
            getValue: (player: Player<Spec.SpecHunter>) => number;
            setValue: (eventID: EventID, player: Player<Spec.SpecHunter>, newValue: number) => void;
        };
    })[];
};
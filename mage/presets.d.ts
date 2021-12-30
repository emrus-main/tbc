import { Consumes } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Mage_Rotation as MageRotation, Mage_Options as MageOptions } from '/tbc/core/proto/mage.js';
export declare const ArcaneTalents: {
    name: string;
    data: string;
};
export declare const FireTalents: {
    name: string;
    data: string;
};
export declare const FrostTalents: {
    name: string;
    data: string;
};
export declare const DefaultFireRotation: MageRotation;
export declare const DefaultFireOptions: MageOptions;
export declare const DefaultFireConsumes: Consumes;
export declare const P1_FIRE_PRESET: {
    name: string;
    tooltip: string;
    gear: EquipmentSpec;
};
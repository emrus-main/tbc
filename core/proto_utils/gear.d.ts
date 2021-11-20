import { ItemSlot } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { EquippedItem } from './equipped_item.js';
declare type InternalGear = Record<ItemSlot, EquippedItem | null>;
/**
 * Represents a full gear set, including items/enchants/gems for every slot.
 *
 * This is an immutable type.
 */
export declare class Gear {
    private readonly gear;
    constructor(gear: Partial<InternalGear>);
    equals(other: Gear): boolean;
    /**
     * Returns a new Gear set with the item equipped.
     *
     * Checks for validity and removes/exchanges items/gems as needed.
     */
    withEquippedItem(newSlot: ItemSlot, newItem: EquippedItem | null): Gear;
    getEquippedItem(slot: ItemSlot): EquippedItem | null;
    asArray(): Array<EquippedItem | null>;
    asSpec(): EquipmentSpec;
}
export {};
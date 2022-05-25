import { Class } from '/tbc/core/proto/common.js';
import { Debuffs } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Raid as RaidProto } from '/tbc/core/proto/api.js';
import { RaidBuffs } from '/tbc/core/proto/common.js';
import { Party } from './party.js';
import { Player } from './player.js';
import { EventID, TypedEvent } from './typed_event.js';
import { Sim } from './sim.js';
export declare const MAX_NUM_PARTIES = 5;
export declare class Raid {
    private buffs;
    private debuffs;
    private tanks;
    private staggerStormstrikes;
    readonly compChangeEmitter: TypedEvent<void>;
    readonly buffsChangeEmitter: TypedEvent<void>;
    readonly debuffsChangeEmitter: TypedEvent<void>;
    readonly tanksChangeEmitter: TypedEvent<void>;
    readonly staggerStormstrikesChangeEmitter: TypedEvent<void>;
    readonly changeEmitter: TypedEvent<void>;
    private parties;
    readonly sim: Sim;
    constructor(sim: Sim);
    size(): number;
    isEmpty(): boolean;
    getParties(): Array<Party>;
    getParty(index: number): Party;
    getPlayers(): Array<Player<any> | null>;
    getPlayer(index: number): Player<any> | null;
    getPlayerFromRaidTarget(raidTarget: RaidTarget): Player<any> | null;
    setPlayer(eventID: EventID, index: number, newPlayer: Player<any> | null): void;
    getClassCount(playerClass: Class): number;
    getBuffs(): RaidBuffs;
    setBuffs(eventID: EventID, newBuffs: RaidBuffs): void;
    getDebuffs(): Debuffs;
    setDebuffs(eventID: EventID, newDebuffs: Debuffs): void;
    getTanks(): Array<RaidTarget>;
    setTanks(eventID: EventID, newTanks: Array<RaidTarget>): void;
    getStaggerStormstrikes(): boolean;
    setStaggerStormstrikes(eventID: EventID, newValue: boolean): void;
    toProto(forExport?: boolean): RaidProto;
    fromProto(eventID: EventID, proto: RaidProto): void;
    clear(eventID: EventID): void;
}

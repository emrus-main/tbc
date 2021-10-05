import { Stat } from '../proto/common.js';
/**
 * Represents values for all character stats (stam, agi, spell power, hit raiting, etc).
 *
 * This is an immutable type.
 */
export declare class Stats {
    private readonly stats;
    constructor(stats?: Array<number>);
    equals(other: Stats): boolean;
    getStat(stat: Stat): number;
    withStat(stat: Stat, value: number): Stats;
    computeEP(epWeights: Stats): number;
    asArray(): Array<number>;
    toJson(): Object;
    static fromJson(obj: any): Stats;
    static fromMap(statsMap: Partial<Record<Stat, number>>): Stats;
}
import { IndividualSimResult } from '../proto/api.js';
import { StatWeightsRequest, StatWeightsResult } from '../proto/api.js';
import { Stat } from '../proto/common.js';
import { Component } from './component.js';
export declare class Results extends Component {
    readonly pendingElem: HTMLDivElement;
    readonly simElem: HTMLDivElement;
    readonly epElem: HTMLDivElement;
    constructor(parent: HTMLElement);
    hideAll(): void;
    setPending(): void;
    setSimResult(result: IndividualSimResult): void;
    setStatWeights(request: StatWeightsRequest, result: StatWeightsResult, epStats: Array<Stat>): void;
}
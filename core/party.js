import { Party as PartyProto } from '/tbc/core/proto/api.js';
import { Player as PlayerProto } from '/tbc/core/proto/api.js';
import { PartyBuffs } from '/tbc/core/proto/common.js';
import { Class } from '/tbc/core/proto/common.js';
import { playerToSpec } from '/tbc/core/proto_utils/utils.js';
import { Player } from './player.js';
import { TypedEvent } from './typed_event.js';
export const MAX_PARTY_SIZE = 5;
// Manages all the settings for a single Party.
export class Party {
    constructor(raid, sim) {
        this.buffs = PartyBuffs.create();
        // Emits when a party member is added/removed/moved.
        this.compChangeEmitter = new TypedEvent();
        this.buffsChangeEmitter = new TypedEvent();
        // Emits when anything in the party changes.
        this.changeEmitter = new TypedEvent();
        this.sim = sim;
        this.raid = raid;
        this.players = [...Array(MAX_PARTY_SIZE).keys()].map(i => null);
        this.playerChangeListener = eventID => this.changeEmitter.emit(eventID);
        [
            this.compChangeEmitter,
            this.buffsChangeEmitter,
        ].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));
    }
    size() {
        return this.players.filter(player => player != null).length;
    }
    isEmpty() {
        return this.size() == 0;
    }
    clear(eventID) {
        this.setBuffs(eventID, PartyBuffs.create());
        for (let i = 0; i < MAX_PARTY_SIZE; i++) {
            this.setPlayer(eventID, i, null);
        }
    }
    // Returns this party's index within the raid [0-4].
    getIndex() {
        return this.raid.getParties().indexOf(this);
    }
    getPlayers() {
        // Make defensive copy.
        return this.players.slice();
    }
    getPlayer(playerIndex) {
        return this.players[playerIndex];
    }
    setPlayer(eventID, playerIndex, newPlayer) {
        if (playerIndex < 0 || playerIndex >= MAX_PARTY_SIZE) {
            throw new Error('Invalid player index: ' + playerIndex);
        }
        if (newPlayer == this.players[playerIndex]) {
            return;
        }
        TypedEvent.freezeAllAndDo(() => {
            const oldPlayer = this.players[playerIndex];
            this.players[playerIndex] = newPlayer;
            if (oldPlayer != null) {
                oldPlayer.changeEmitter.off(this.playerChangeListener);
                oldPlayer.setParty(null);
            }
            if (newPlayer != null) {
                const newPlayerOldParty = newPlayer.getParty();
                if (newPlayerOldParty) {
                    newPlayerOldParty.setPlayer(eventID, newPlayer.getPartyIndex(), null);
                }
                newPlayer.changeEmitter.on(this.playerChangeListener);
                newPlayer.setParty(this);
            }
            this.compChangeEmitter.emit(eventID);
        });
    }
    getBuffs() {
        // Make a defensive copy
        return PartyBuffs.clone(this.buffs);
    }
    setBuffs(eventID, newBuffs) {
        if (PartyBuffs.equals(this.buffs, newBuffs))
            return;
        // Make a defensive copy
        this.buffs = PartyBuffs.clone(newBuffs);
        this.buffsChangeEmitter.emit(eventID);
    }
    toProto() {
        return PartyProto.create({
            players: this.players.map(player => player == null ? PlayerProto.create() : player.toProto()),
            buffs: this.buffs,
        });
    }
    fromProto(eventID, proto) {
        TypedEvent.freezeAllAndDo(() => {
            this.setBuffs(eventID, proto.buffs || PartyBuffs.create());
            for (let i = 0; i < MAX_PARTY_SIZE; i++) {
                if (!proto.players[i] || proto.players[i].class == Class.ClassUnknown) {
                    this.setPlayer(eventID, i, null);
                    continue;
                }
                const playerProto = proto.players[i];
                const spec = playerToSpec(playerProto);
                const currentPlayer = this.players[i];
                // Reuse the current player if possible, so that event handlers are preserved.
                if (currentPlayer && spec == currentPlayer.spec) {
                    currentPlayer.fromProto(eventID, playerProto);
                }
                else {
                    const newPlayer = new Player(spec, this.sim);
                    newPlayer.fromProto(eventID, playerProto);
                    this.setPlayer(eventID, i, newPlayer);
                }
            }
        });
    }
    // Returns JSON representing all the current values.
    toJson() {
        return {
            'players': this.players.map(player => {
                if (player == null) {
                    return null;
                }
                else {
                    return {
                        'spec': player.spec,
                        'player': player.toJson(),
                    };
                }
            }),
            'buffs': PartyBuffs.toJson(this.getBuffs()),
        };
    }
    // Set all the current values, assumes obj is the same type returned by toJson().
    fromJson(eventID, obj) {
        TypedEvent.freezeAllAndDo(() => {
            try {
                this.setBuffs(eventID, PartyBuffs.fromJson(obj['buffs']));
            }
            catch (e) {
                console.warn('Failed to parse party buffs: ' + e);
            }
            if (obj['players']) {
                for (let i = 0; i < MAX_PARTY_SIZE; i++) {
                    const playerObj = obj['players'][i];
                    if (!playerObj) {
                        this.setPlayer(eventID, i, null);
                        continue;
                    }
                    const newSpec = playerObj['spec'];
                    if (this.players[i] != null && this.players[i].spec == newSpec) {
                        this.players[i].fromJson(eventID, playerObj['player']);
                    }
                    else {
                        const newPlayer = new Player(playerObj['spec'], this.sim);
                        newPlayer.fromJson(eventID, playerObj['player']);
                        this.setPlayer(eventID, i, newPlayer);
                    }
                }
            }
        });
    }
}

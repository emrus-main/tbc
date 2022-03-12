import { ActionMetrics } from '/tbc/core/proto_utils/sim_result.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
export class CastMetricsTable extends MetricsTable {
    constructor(config) {
        config.rootCssClass = 'cast-metrics-root';
        super(config, [
            MetricsTable.nameCellConfig((metric) => {
                return {
                    name: metric.name,
                    actionId: metric.actionId,
                };
            }),
            {
                name: 'Casts',
                tooltip: 'Casts',
                sort: ColumnSortType.Descending,
                getValue: (metric) => metric.casts,
                getDisplayString: (metric) => metric.casts.toFixed(1),
            },
            {
                name: 'CPM',
                tooltip: 'Casts / (Encounter Duration / 60 Seconds)',
                getValue: (metric) => metric.castsPerMinute,
                getDisplayString: (metric) => metric.castsPerMinute.toFixed(1),
            },
        ]);
    }
    getGroupedMetrics(resultData) {
        //const actionMetrics = resultData.result.getActionMetrics(resultData.filter);
        const players = resultData.result.getPlayers(resultData.filter);
        if (players.length != 1) {
            return [];
        }
        const player = players[0];
        const actions = player.actions;
        const actionGroups = ActionMetrics.groupById(actions);
        const petGroups = player.pets.map(pet => pet.actions);
        return actionGroups.concat(petGroups);
    }
    mergeMetrics(metrics) {
        return ActionMetrics.merge(metrics, true, metrics[0].player?.petActionId || undefined);
    }
    shouldCollapse(metric) {
        return !metric.player?.isPet;
    }
}

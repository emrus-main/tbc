import { Gem } from '/tbc/core/proto/common.js';
import { GemColor } from '/tbc/core/proto/common.js';

const socketToMatchingColors = new Map<GemColor, Array<GemColor>>();
socketToMatchingColors.set(GemColor.GemColorMeta,   [GemColor.GemColorMeta]);
socketToMatchingColors.set(GemColor.GemColorBlue,   [GemColor.GemColorBlue, GemColor.GemColorPurple, GemColor.GemColorGreen]);
socketToMatchingColors.set(GemColor.GemColorRed,    [GemColor.GemColorRed, GemColor.GemColorPurple, GemColor.GemColorOrange]);
socketToMatchingColors.set(GemColor.GemColorYellow, [GemColor.GemColorYellow, GemColor.GemColorOrange, GemColor.GemColorGreen]);

// Whether the gem matches the given socket color, for the purposes of gaining the socket bonuses.
export function gemMatchesSocket(gem: Gem, socketColor: GemColor) {
  return socketToMatchingColors.has(socketColor) && socketToMatchingColors.get(socketColor)!.includes(gem.color);
}

// Whether the gem is capable of slotting into a socket of the given color.
export function gemEligibleForSocket(gem: Gem, socketColor: GemColor) {
  return (gem.color == GemColor.GemColorMeta) == (socketColor == GemColor.GemColorMeta);
}


// Maps meta gem IDs to functions that check whether they're active.
const metaGemActiveConditions = new Map<number, (numRed: number, numYellow: number, numBlue: number) => boolean>();

// Maps meta gem IDs to string descriptions of the meta conditions.
const metaGemConditionDescriptions = new Map<number, string>();

export function isMetaGemActive(metaGem: Gem, numRed: number, numYellow: number, numBlue: number): boolean {
	if (!metaGemActiveConditions.has(metaGem.id)) {
		// If we don't have a condition for this meta gem, just default to active.
		return true;
	}

	return metaGemActiveConditions.get(metaGem.id)!(numRed, numYellow, numBlue);
}

export function getMetaGemConditionDescription(metaGem: Gem): string {
	return metaGemConditionDescriptions.get(metaGem.id) || '';
}

// Keep these lists in alphabetical order, separated by color.

// Meta
export const BRACING_EARTHSTORM_DIAMOND = 25897;
metaGemConditionDescriptions.set(BRACING_EARTHSTORM_DIAMOND, 'Requires more Red Gems than Blue Gems.');
metaGemActiveConditions.set(BRACING_EARTHSTORM_DIAMOND, (numRed, numYellow, numBlue) => numRed > numBlue);
export const BRUTAL_EARTHSTORM_DIAMOND = 25899;
metaGemConditionDescriptions.set(BRUTAL_EARTHSTORM_DIAMOND, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems');
metaGemActiveConditions.set(BRUTAL_EARTHSTORM_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 2 && numRed >= 2 && numBlue >= 2);
export const CHAOTIC_SKYFIRE_DIAMOND = 34220;
metaGemConditionDescriptions.set(CHAOTIC_SKYFIRE_DIAMOND, 'Requires at least 2 Blue Gems.');
metaGemActiveConditions.set(CHAOTIC_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numBlue >= 2);
export const DESTRUCTIVE_SKYFIRE_DIAMOND = 25890;
metaGemConditionDescriptions.set(DESTRUCTIVE_SKYFIRE_DIAMOND, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue gems.');
metaGemActiveConditions.set(DESTRUCTIVE_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numRed >= 2 && numYellow >= 2 && numBlue >= 2);
export const EMBER_SKYFIRE_DIAMOND = 35503;
metaGemConditionDescriptions.set(EMBER_SKYFIRE_DIAMOND, 'Requires at least 3 Red Gems.');
metaGemActiveConditions.set(EMBER_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numRed >= 3);
export const ENIGMATIC_SKYFIRE_DIAMOND = 25895;
metaGemConditionDescriptions.set(ENIGMATIC_SKYFIRE_DIAMOND, 'Requires more Red Gems than Yellow Gems.');
metaGemActiveConditions.set(ENIGMATIC_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numRed > numYellow);
export const IMBUED_UNSTABLE_DIAMOND = 32641;
metaGemConditionDescriptions.set(IMBUED_UNSTABLE_DIAMOND, 'Requires at least 3 Yellow Gems.');
metaGemActiveConditions.set(IMBUED_UNSTABLE_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 3);
export const INSIGHTFUL_EARTHSTORM_DIAMOND = 25901;
metaGemConditionDescriptions.set(INSIGHTFUL_EARTHSTORM_DIAMOND, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue gems.');
metaGemActiveConditions.set(INSIGHTFUL_EARTHSTORM_DIAMOND, (numRed, numYellow, numBlue) => numRed >= 2 && numYellow >= 2 && numBlue >= 2);
export const MYSTICAL_SKYFIRE_DIAMOND = 25893;
metaGemConditionDescriptions.set(MYSTICAL_SKYFIRE_DIAMOND, 'Requires more Blue Gems than Yellow Gems.');
metaGemActiveConditions.set(MYSTICAL_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numBlue > numYellow);
export const POTENT_UNSTABLE_DIAMOND = 32640;
metaGemConditionDescriptions.set(POTENT_UNSTABLE_DIAMOND, 'Requires more Blue Gems than Yellow Gems.');
metaGemActiveConditions.set(POTENT_UNSTABLE_DIAMOND, (numRed, numYellow, numBlue) => numBlue > numYellow);
export const POWERFUL_EARTHSTORM_DIAMOND = 25896;
metaGemConditionDescriptions.set(POWERFUL_EARTHSTORM_DIAMOND, 'Requires at least 3 Blue Gems.');
metaGemActiveConditions.set(POWERFUL_EARTHSTORM_DIAMOND, (numRed, numYellow, numBlue) => numBlue >= 3);
export const RELENTLESS_EARTHSTORM_DIAMOND = 32409;
metaGemConditionDescriptions.set(RELENTLESS_EARTHSTORM_DIAMOND, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems');
metaGemActiveConditions.set(RELENTLESS_EARTHSTORM_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 2 && numRed >= 2 && numBlue >= 2);
export const SWIFT_SKYFIRE_DIAMOND = 25894;
metaGemConditionDescriptions.set(SWIFT_SKYFIRE_DIAMOND, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.');
metaGemActiveConditions.set(SWIFT_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 2 && numRed >= 1);
export const SWIFT_STARFIRE_DIAMOND = 28557;
metaGemConditionDescriptions.set(SWIFT_STARFIRE_DIAMOND, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.');
metaGemActiveConditions.set(SWIFT_STARFIRE_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 2 && numRed >= 1);
export const SWIFT_WINDFIRE_DIAMOND = 28556;
metaGemConditionDescriptions.set(SWIFT_WINDFIRE_DIAMOND, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.');
metaGemActiveConditions.set(SWIFT_WINDFIRE_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 2 && numRed >= 1);
export const TENACIOUS_EARTHSTORM_DIAMOND = 25898;
metaGemConditionDescriptions.set(TENACIOUS_EARTHSTORM_DIAMOND, 'Requires at least 5 Blue Gems.');
metaGemActiveConditions.set(TENACIOUS_EARTHSTORM_DIAMOND, (numRed, numYellow, numBlue) => numBlue >= 5);
export const THUNDERING_SKYFIRE_DIAMOND = 32410;
metaGemConditionDescriptions.set(THUNDERING_SKYFIRE_DIAMOND, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems');
metaGemActiveConditions.set(THUNDERING_SKYFIRE_DIAMOND, (numRed, numYellow, numBlue) => numYellow >= 2 && numRed >= 2 && numBlue >= 2);

// Orange
export const GLINTING_NOBLE_TOPAZ = 24061;
export const GLINTING_PYRESTONE = 32220;
export const INSCRIBED_NOBLE_TOPAZ = 24058;
export const INSCRIBED_PYRESTONE = 32217;
export const POTENT_NOBLE_TOPAZ = 24059;
export const POTENT_PYRESTONE = 32218;
export const VEILED_NOBLE_TOPAZ = 31867;
export const VEILED_PYRESTONE = 32221;
export const WICKED_NOBLE_TOPAZ = 31868;
export const WICKED_PYRESTONE = 32222;

// Purple
export const GLOWING_NIGHTSEYE = 24056;
export const GLOWING_SHADOWSONG_AMETHYST = 32215;
export const SHIFTING_NIGHTSEYE = 24055;
export const SHIFTING_SHADOWSONG_AMETHYST = 32212;
export const SOVEREIGN_NIGHTSEYE = 24054;
export const SOVEREIGN_SHADOWSONG_AMETHYST = 32211;

// Red
export const BOLD_CRIMSON_SPINEL = 32193;
export const BOLD_LIVING_RUBY = 24027;
export const DELICATE_CRIMSON_SPINEL = 32194;
export const DELICATE_LIVING_RUBY = 24028;
export const RUNED_CRIMSON_SPINEL = 32196;
export const RUNED_LIVING_RUBY = 24030;
export const RUNED_ORNATE_RUBY = 28118;

// Yellow
export const BRILLIANT_DAWNSTONE = 24047;
export const BRILLIANT_LIONSEYE = 32204;
export const RIGID_DAWNSTONE = 24051;
export const RIGID_LIONSEYE = 32206;

// Green
export const JAGGED_SEASPRAY_EMERALD = 32226;
export const JAGGED_TALASITE = 24067;

const gemSocketCssClasses: Partial<Record<GemColor, string>> = {
  [GemColor.GemColorBlue]: 'socket-color-blue',
  [GemColor.GemColorMeta]: 'socket-color-meta',
  [GemColor.GemColorRed]: 'socket-color-red',
  [GemColor.GemColorYellow]: 'socket-color-yellow',
};
export function setGemSocketCssClass(elem: HTMLElement, color: GemColor) {
  Object.values(gemSocketCssClasses).forEach(cssClass => elem.classList.remove(cssClass));

  if (gemSocketCssClasses[color]) {
    elem.classList.add(gemSocketCssClasses[color] as string);
    return;
  }

  throw new Error('No css class for gem socket color: ' + color);
}

const emptyGemSocketIcons: Partial<Record<GemColor, string>> = {
  [GemColor.GemColorBlue]: 'https://wow.zamimg.com/images/icons/socket-blue.gif',
  [GemColor.GemColorMeta]: 'https://wow.zamimg.com/images/icons/socket-meta.gif',
  [GemColor.GemColorRed]: 'https://wow.zamimg.com/images/icons/socket-red.gif',
  [GemColor.GemColorYellow]: 'https://wow.zamimg.com/images/icons/socket-yellow.gif',
};
export function getEmptyGemSocketIconUrl(color: GemColor): string {
  if (emptyGemSocketIcons[color])
    return emptyGemSocketIcons[color] as string;

  throw new Error('No empty socket url for gem socket color: ' + color);
}

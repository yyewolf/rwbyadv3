export type Card = {
	id: string;
	card_type: string;
	level: number;
	rarity: number | undefined;
	individual_value: number;
	edges: {
		type: CardType;
	};
};

export type CardType = {
	name: string;
	categories: string[];
};

export function getCardIconUri(cardType: string): string {
	return `/cdn/cards/${cardType}/icon.webp`;
}

// Function to convert rarity to hexadecimal color (similar to RarityToColor in Go)
export function rarityToColor(rarity: number | undefined): string {
	rarity = rarity || 0; // Default to 0 if rarity is undefined
	const colors = ['#B5B5B5', '#9DF3C5', '#7EB4FF', '#D887F5', '#FFD700'];
	return colors[Math.min(rarity - 1, colors.length - 1)] || colors[0];
}

export function getCardRarityName(rarity: number | undefined): string {
	// Default to 'Common' if rarity is undefined
	const rarities = ['Common', 'Uncommon', 'Rare', 'Epic', 'Legendary'];
	return rarities[Math.min((rarity || 1) - 1, rarities.length - 1)];
}

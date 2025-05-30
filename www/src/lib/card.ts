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
	const colors = ['#808080', '#7CFC00', '#87CEEB', '#BA55D3', '#FFD700', '#FF0000']; // Hex colors for Common, Uncommon, Rare, Very Rare, Legendary, Collector
	return colors[rarity] || '#808080'; // Default to gray if rarity is out of bounds
}

export function getCardRarityName(rarity: number | undefined): string {
	// Default to 'Common' if rarity is undefined
	const rarities = ['Common', 'Uncommon', 'Rare', 'Very Rare', 'Legendary', 'Collector'];
	return rarities[rarity || 0] || 'Common'; // Default to 'Common' if rarity is out of bounds
}

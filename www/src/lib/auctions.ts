import type { Auction } from './types';

// Get current highest bid amount or starting price
export function getCurrentPrice(auction: Auction | null): number {
	if (!auction) return 0;

	// If there are bids, return the highest one
	if (auction.edges.bids && auction.edges.bids.length > 0) {
		return auction.edges.bids[0].price;
	}

	// Otherwise return the starting price
	return 0;
}

// Calculate minimum bid
export function calculateMinimumBid(auction: Auction | null): number {
	if (!auction) return 0;

	// If there are bids, minimum is highest bid + 50
	if (auction.edges.bids && auction.edges.bids.length > 0) {
		return auction.edges.bids[0].price + 50;
	}

	// Otherwise minimum is starting price + 50
	return 50;
}

// Auction timer related functions
export type TimeRemaining = {
	days: number;
	hours: number;
	minutes: number;
	seconds: number;
	hasEnded: boolean;
	formatted: string;
	shortFormatted: string;
};

/**
 * Calculate time remaining for an auction
 * @param endDateString ISO date string for auction end time
 * @param includeSeconds Whether to include seconds in the calculation
 * @returns TimeRemaining object with various formats
 */
export function calculateTimeRemaining(
	endDateString: string,
	includeSeconds = true
): TimeRemaining {
	const now = new Date();
	const endDate = new Date(endDateString);
	const diff = endDate.getTime() - now.getTime();
	const hasEnded = diff <= 0;

	if (hasEnded) {
		return {
			days: 0,
			hours: 0,
			minutes: 0,
			seconds: 0,
			hasEnded: true,
			formatted: 'Auction ended',
			shortFormatted: 'Ended'
		};
	}

	const days = Math.floor(diff / (1000 * 60 * 60 * 24));
	const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
	const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
	const seconds = includeSeconds ? Math.floor((diff % (1000 * 60)) / 1000) : 0;

	// Create formatted strings
	let formatted = '';
	let shortFormatted = '';

	if (days > 0) {
		formatted = `${days}d ${hours}h ${minutes}m${includeSeconds ? ` ${seconds}s` : ''}`;
		shortFormatted = `${days}d ${hours}h`;
	} else if (hours > 0) {
		formatted = `${hours}h ${minutes}m${includeSeconds ? ` ${seconds}s` : ''}`;
		shortFormatted = `${hours}h ${minutes}m`;
	} else {
		formatted = `${minutes}m${includeSeconds ? ` ${seconds}s` : ''}`;
		shortFormatted = `${minutes}m`;
	}

	return {
		days,
		hours,
		minutes,
		seconds,
		hasEnded: false,
		formatted,
		shortFormatted
	};
}

/**
 * Check if an auction is in its final minutes
 * @param timeRemaining TimeRemaining object
 * @param thresholdMinutes Number of minutes considered as "ending soon"
 * @returns boolean
 */
export function isAuctionEndingSoon(timeRemaining: TimeRemaining, thresholdMinutes = 5): boolean {
	if (timeRemaining.hasEnded) return false;

	const totalMinutesLeft =
		timeRemaining.days * 24 * 60 + timeRemaining.hours * 60 + timeRemaining.minutes;
	return totalMinutesLeft <= thresholdMinutes;
}

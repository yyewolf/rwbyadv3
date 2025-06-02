import type { Card } from './card';
import { post, get } from './request';
import type { Response } from './types';

export interface CreateTradeRequest {
	offer: string[]; // Card IDs that the current player is offering
	receive: string[]; // Card IDs that the current player wants to receive
}

/**
 * Get all cards for the currently logged-in user that are available for trade
 * @returns Promise with array of cards
 */
export async function getMyCards(): Promise<Card[]> {
	const resp = await get<Card[]>('/apis/trades/self/cards');
	return resp.data;
}

/**
 * Get all tradable cards for a specific player
 * @param playerId The ID of the player whose cards to fetch
 * @returns Promise with array of cards
 */
export async function getPlayerCards(playerId: string): Promise<Card[]> {
	const resp = await get<Card[]>(`/apis/trades/${playerId}/cards`);
	return resp.data;
}

/**
 * Submit a trade offer to another player
 * @param playerId The ID of the player to trade with
 * @param offer The card IDs that the current user is offering
 * @param receive The card IDs that the current user wants to receive
 * @returns Promise with the created trade
 */
export async function createTrade(
	playerId: string,
	offer: Card[],
	receive: Card[]
): Promise<Response<any>> {
	const tradeData: CreateTradeRequest = {
		offer: offer.map((card) => card.id),
		receive: receive.map((card) => card.id)
	};

	return post<any>(`/apis/trades/${playerId}`, tradeData);
}

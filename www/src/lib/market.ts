import type { Auction, AuctionBid, Listing, PaginatedResponse, Response } from './types';
import { get, post, getPaginated } from './request';
import { marketSSE } from './sseClient';

// Export the SSE client for market real-time updates
export { marketSSE };

/**
 * Get paginated listings based on search query
 * @param query The search query
 * @param page The page number
 * @returns Paginated response with listings
 */
export const getListings = async (
	query: string,
	page: number
): Promise<PaginatedResponse<Listing>> => {
	return await getPaginated<Listing>(
		`/apis/market/api/listings?query=${encodeURIComponent(query)}&p=${page}`
	);
};

/**
 * Get the latest listings
 * @returns Response with array of latest listings
 */
export const getLatestListings = async (): Promise<Response<Listing[]>> => {
	return await get<Listing[]>(`/apis/market/api/latest/listings`);
};

/**
 * Purchase a listing
 * @param listingId The ID of the listing to purchase
 * @returns Response with the purchased listing or null if redirected
 */
export const purchaseListing = async (listingId: string): Promise<Response<Listing | null>> => {
	return await post<Listing | null>(`/apis/market/api/listings/${listingId}`, undefined);
};

/**
 * Refresh the latest listings and update the SSE store
 * This can be called manually when needed, instead of waiting for SSE updates
 */
export const refreshLatestListings = async (): Promise<void> => {
	try {
		const response = await getLatestListings();
		if (response.meta.success) {
			marketSSE.latestListings.set(response.data);
		}
	} catch (error) {
		console.error('Failed to refresh latest listings:', error);
	}
};

/**
 * Get paginated auctions based on search query
 * @param query The search query
 * @param page The page number
 * @returns Paginated response with auctions
 */
export const getAuctions = async (
	query: string,
	page: number
): Promise<PaginatedResponse<Auction>> => {
	return await getPaginated<Auction>(
		`/apis/market/api/auctions?query=${encodeURIComponent(query)}&p=${page}`
	);
};

/**
 * Get a specific auction by ID
 * @param auctionId The ID of the auction to fetch
 * @returns Response with the auction details
 */
export const getAuction = async (auctionId: string): Promise<Response<Auction>> => {
	return await get<Auction>(`/apis/market/api/auctions/${auctionId}`);
};

/**
 * Get the latest auctions
 * @returns Response with array of latest auctions
 */
export const getLatestAuctions = async (): Promise<Response<Auction[]>> => {
	return await get<Auction[]>(`/apis/market/api/latest/auctions`);
};

/**
 * Bid on an auction
 * @param auctionId The ID of the auction to bid on
 * @param bidAmount The amount to bid
 * @returns Response with the updated auction or null if redirected
 */
export const bidOnAuction = async (
	auctionId: string,
	bidAmount: number
): Promise<Response<Auction | null>> => {
	return await post<Auction | null>(`/apis/market/api/auctions/${auctionId}`, {
		bid_amount: bidAmount
	});
};

/**
 * Refresh the latest auctions and update the SSE store
 * This can be called manually when needed, instead of waiting for SSE updates
 */
export const refreshLatestAuctions = async (): Promise<void> => {
	try {
		const response = await getLatestAuctions();
		if (response.meta.success) {
			marketSSE.latestAuctions.set(response.data);
		}
	} catch (error) {
		console.error('Failed to refresh latest auctions:', error);
	}
};

/**
 * Subscribe to real-time updates for a specific auction's bids
 * @param auctionId The ID of the auction to watch for bid updates
 * @param callback Function to call when a new bid is placed
 * @returns Subscription object that can be used to unsubscribe
 */
export const subscribeToAuctionBids = (
	auctionId: string,
	callback: (bidData: AuctionBid) => void
) => {
	return marketSSE.subscribe(`auction_${auctionId}_bid`, callback);
};

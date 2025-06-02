import { writable, type Writable } from 'svelte/store';
import type { Auction, Listing, Response, AuctionBid } from './types';
import { notifyError } from './stores/notificationStore';

/**
 * Interface for SSE connection options
 */
interface SSEOptions {
	reconnectDelay?: number;
	maxReconnectAttempts?: number;
}

/**
 * Available event types for the market SSE
 */
export type MarketEventType =
	| 'latest_listings'
	| 'latest_auctions'
	| 'auction_bid'
	| 'auction_remove'
	| 'listing_remove'; // Event types for market updates

/**
 * Event callback type for SSE events
 */
type EventCallback = (data: any) => void;

/**
 * Subscription object returned when registering for events
 */
export interface Subscription {
	unsubscribe: () => void;
}

/**
 * Market SSE client to handle real-time updates
 */
class MarketSSEClient {
	private eventSource: EventSource | null = null;
	private reconnectAttempts = 0;
	private reconnectDelay = 1000;
	private maxReconnectAttempts = 5;
	private isConnected = false;
	private eventListeners: Map<string, Set<EventCallback>> = new Map();
	private eventListenersCache: Map<string, any> = new Map();

	// Stores for different types of data
	public latestListings: Writable<Listing[]> = writable([]);
	public latestAuctions: Writable<Auction[]> = writable([]);

	/**
	 * Connect to the SSE endpoint
	 * @param options Configuration options for the SSE connection
	 */
	public connect(options?: SSEOptions): void {
		if (this.eventSource) {
			return; // Already connected
		}

		if (options) {
			this.reconnectDelay = options.reconnectDelay ?? this.reconnectDelay;
			this.maxReconnectAttempts = options.maxReconnectAttempts ?? this.maxReconnectAttempts;
		}

		try {
			this.eventSource = new EventSource('/apis/market/sse');

			// Set up event listeners
			this.eventSource.onopen = this.handleOpen.bind(this);
			this.eventSource.onerror = this.handleError.bind(this);

			// Set up message handlers
			this.eventSource.addEventListener('latest_listings', this.handleLatestListings.bind(this));
			this.eventSource.addEventListener('latest_auctions', this.handleLatestAuctions.bind(this));
			this.eventSource.addEventListener('auction_bid', this.handleAuctionBid.bind(this));
			this.eventSource.addEventListener('auction_remove', this.handleAuctionRemove.bind(this));
			this.eventSource.addEventListener('listing_remove', this.handleListingRemove.bind(this));
			this.eventSource.addEventListener('ping', () => {}); // Ignore ping events
		} catch (error) {
			console.error('Failed to connect to SSE:', error);
			this.attemptReconnect();
		}
	}

	/**
	 * Disconnect from the SSE endpoint
	 */
	public disconnect(): void {
		if (this.eventSource) {
			this.eventSource.close();
			this.eventSource = null;
			this.isConnected = false;
			this.reconnectAttempts = 0;
		}
	}

	/**
	 * Handler for successful connection
	 */
	private handleOpen(): void {
		this.isConnected = true;
		this.reconnectAttempts = 0;
		console.log('SSE connection established');
	}

	/**
	 * Handler for connection errors
	 */
	private handleError(event: Event): void {
		if (this.isConnected) {
			console.error('SSE connection error:', event);
		}

		this.isConnected = false;
		this.eventSource?.close();
		this.eventSource = null;
		this.attemptReconnect();
	}

	/**
	 * Attempt to reconnect to the SSE endpoint
	 */
	private attemptReconnect(): void {
		if (this.reconnectAttempts >= this.maxReconnectAttempts) {
			console.error(`Failed to reconnect after ${this.maxReconnectAttempts} attempts`);
			notifyError('Connection to market updates lost. Please refresh the page.');
			return;
		}

		this.reconnectAttempts++;
		console.log(
			`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`
		);

		setTimeout(() => {
			this.connect();
		}, this.reconnectDelay * this.reconnectAttempts);
	}

	/**
	 * Handler for latest_listings events
	 */
	private handleLatestListings(event: MessageEvent): void {
		try {
			// Parse the event data
			const response = JSON.parse(event.data) as Listing[];

			// Update the store with new data
			this.latestListings.set(response);
		} catch (error) {
			console.error('Failed to parse latest listings:', error);
		}
	}

	/**
	 * Handler for latest_auctions events
	 */
	private handleLatestAuctions(event: MessageEvent): void {
		try {
			// Parse the event data
			const response = JSON.parse(event.data) as Auction[];

			// Update the store with new data
			this.latestAuctions.set(response);
		} catch (error) {
			console.error('Failed to parse latest auctions:', error);
		}
	}

	/**
	 * Handler for auction_bid events
	 */
	private handleAuctionBid(event: MessageEvent): void {
		try {
			const bidData = JSON.parse(event.data) as AuctionBid;
			const auctionId = bidData.auction_id;

			// Notify all subscribers for this specific auction bid
			this.eventListenersCache.set(`auction_${auctionId}_bid`, bidData);
			this.notifyEventListeners(`auction_${auctionId}_bid`, bidData);
		} catch (error) {
			console.error('Failed to parse auction bid data:', error);
		}
	}

	/**
	 * Handler for auction_remove events
	 */
	private handleAuctionRemove(event: MessageEvent): void {
		try {
			// The data should be the auction ID as a string
			const auctionId = event.data;
			if (!auctionId) return;

			// Create removal event data
			const removalData = { id: auctionId, removed: true, timestamp: new Date().toISOString() };

			// Notify all subscribers for this specific auction removal
			this.eventListenersCache.set(`auction_${auctionId}_remove`, removalData);
			this.notifyEventListeners(`auction_${auctionId}_remove`, removalData);
		} catch (error) {
			console.error('Failed to process auction removal data:', error);
		}
	}

	/**
	 * Handler for listing_remove events
	 */
	private handleListingRemove(event: MessageEvent): void {
		try {
			// The data should be the listing ID as a string
			const listingId = event.data;
			if (!listingId) return;

			// Create removal event data
			const removalData = { id: listingId, removed: true, timestamp: new Date().toISOString() };

			// Notify all subscribers for this specific listing removal
			this.eventListenersCache.set(`listing_${listingId}_remove`, removalData);
			this.notifyEventListeners(`listing_${listingId}_remove`, removalData);
		} catch (error) {
			console.error('Failed to process listing removal data:', error);
		}
	}

	/**
	 * Subscribe to a specific event
	 * @param eventName Name of the event to subscribe to (e.g. "auction_123_bid")
	 * @param callback Function to call when the event occurs
	 * @returns Subscription object with unsubscribe method
	 */
	public subscribe<T>(eventName: string, callback: (data: T) => void): Subscription {
		if (!this.eventListeners.has(eventName)) {
			this.eventListeners.set(eventName, new Set());
		}

		if (this.eventListenersCache.has(eventName)) {
			// If we have cached data for this event, call the callback immediately
			callback(this.eventListenersCache.get(eventName) as T);
		}

		const listeners = this.eventListeners.get(eventName)!;
		listeners.add(callback as EventCallback);

		return {
			unsubscribe: () => {
				listeners.delete(callback as EventCallback);
				if (listeners.size === 0) {
					this.eventListeners.delete(eventName);
				}
			}
		};
	}

	/**
	 * Notify all listeners for a specific event
	 * @param eventName Name of the event
	 * @param data Data to pass to the event listeners
	 */
	private notifyEventListeners(eventName: string, data: any): void {
		const listeners = this.eventListeners.get(eventName);
		if (!listeners) return;

		listeners.forEach((callback) => {
			try {
				callback(data);
			} catch (error) {
				console.error(`Error in event listener for ${eventName}:`, error);
			}
		});
	}
}

// Create a singleton instance for the whole app
export const marketSSE = new MarketSSEClient();

// {"meta":{"success":true},"data":[{"id":"883c44dd-fe35-43ca-ae0a-ef336d2902a8","card_id":"bb937b0d-d0bb-4995-971c-5a99ef109465","price":100,"edges":{"owned_by":{"id":"144472011924570113","edges":{}},"card":{"id":"bb937b0d-d0bb-4995-971c-5a99ef109465","card_type":"ruby/summer","level":1,"individual_value":21.04781503822545,"edges":{"type":{"name":"Ruby Rose Summer Edition","categories":["Team RWBY","Summer"]},"stats":{"edges":{}}}}}}],"pagination":{"current_page":1,"total_pages":1,"per_page":5,"total_items":1,"has_next":false,"has_prev":false}}

import type { Card } from './card';

export type PaginatedResponse<T> = {
	meta: {
		success: boolean;
	};
	data: T[];
	pagination: {
		current_page: number;
		total_pages: number;
		per_page: number;
		total_items: number;
		has_next: boolean;
		has_prev: boolean;
	};
	error: {
		code?: string;
		message?: string;
	};
};

export type Response<T> = {
	meta: {
		success: boolean;
	};
	data: T;
	error: {
		code?: string;
		message?: string;
		redirect?: string; // Optional redirect URL for 401 errors
	};
};

export type Listing = {
	id: string;
	card_id: string;
	price: number;
	edges: {
		card: Card;
	};
	create_time: string;
	update_time: string;
};

// Type for an auction bid
export type AuctionBid = {
	id: string;
	auction_id: string;
	player_id: string;
	price: number;
	create_time: string;
	edges?: {
		player: {
			id: string;
			username: string;
		};
	};
};

// Type for an auction
export type Auction = {
	id: string;
	card_id: string;
	ends_at: string; // ISO timestamp when the auction ends
	time_extensions: number; // Number of times the auction has been extended
	edges: {
		card: Card;
		owned_by?: {
			id: string;
			username: string;
		};
		bids?: AuctionBid[];
	};
	create_time: string;
	update_time: string;
};

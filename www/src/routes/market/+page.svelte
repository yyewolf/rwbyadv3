<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { gsap } from 'gsap';
	import { getLatestListings, getLatestAuctions, marketSSE } from '$lib/market';
	import type { Listing, Auction } from '$lib/types';
	import ListingModal from '../../components/market/listings/ListingModal.svelte';
	import AuctionModal from '../../components/market/auctions/AuctionModal.svelte';
	import FeaturedListings from '../../components/market/listings/FeaturedListings.svelte';
	import FeaturedAuctions from '../../components/market/auctions/FeaturedAuctions.svelte';

	// References for animation
	let heroSection: HTMLElement;
	let featuredListingCards: HTMLElement;
	let featuredAuctionCards: HTMLElement;

	// Modal states
	let selectedListing: Listing | null = null;
	let showListingModal = false;
	let selectedAuction: Auction | null = null;
	let showAuctionModal = false;

	// State for latest listings and auctions
	let latestListings: Listing[] = [];
	let latestAuctions: Auction[] = [];
	let isListingsLoading = true;
	let isAuctionsLoading = true;
	let listingsError: string | null = null;
	let auctionsError: string | null = null;

	// Subscribe to SSE updates for latest listings and auctions
	const unsubscribeListings = marketSSE.latestListings.subscribe((listings) => {
		latestListings = listings;
		isListingsLoading = false;
	});

	const unsubscribeAuctions = marketSSE.latestAuctions.subscribe((auctions) => {
		latestAuctions = auctions;
		isAuctionsLoading = false;
	});

	// Function to animate cards when they update
	function animateCards(container: HTMLElement | undefined) {
		if (!container) return;

		setTimeout(() => {
			gsap.fromTo(
				container?.querySelectorAll('.card'),
				{ opacity: 0, y: 20, scale: 0.95 },
				{
					opacity: 1,
					y: 0,
					scale: 1,
					stagger: 0.1,
					duration: 0.6,
					ease: 'back.out(1.2)'
				}
			);
		}, 100);
	}

	onMount(async () => {
		// Hero section animation
		gsap.fromTo(
			heroSection,
			{ opacity: 0, y: 30 },
			{ opacity: 1, y: 0, duration: 0.8, ease: 'power2.out' }
		);

		// Connect to SSE and fetch initial data
		marketSSE.connect();

		// Fetch latest listings initially (will be updated by SSE afterward)
		try {
			const listingsResponse = await getLatestListings();
			latestListings = listingsResponse.data;
			isListingsLoading = false;

			// Run animations after data is loaded
			animateCards(featuredListingCards);
		} catch (e) {
			listingsError = e instanceof Error ? e.message : 'Failed to fetch listings';
			isListingsLoading = false;
		}

		// Fetch latest auctions initially
		try {
			const auctionsResponse = await getLatestAuctions();
			latestAuctions = auctionsResponse.data;
			isAuctionsLoading = false;
			// Animate auction cards after they load
			animateCards(featuredAuctionCards);
		} catch (e) {
			auctionsError = e instanceof Error ? e.message : 'Failed to fetch auctions';
			isAuctionsLoading = false;
		}
	});

	// Clean up SSE connection and subscription when component unmounts
	onDestroy(() => {
		unsubscribeListings();
		unsubscribeAuctions();
		marketSSE.disconnect();
	});
</script>

<div>
	<!-- Hero Section -->
	<section
		bind:this={heroSection}
		class="mb-12 overflow-hidden rounded-lg bg-gradient-to-r from-indigo-700 to-purple-600 shadow-lg"
	>
		<div class="relative mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-24 lg:px-8">
			<div class="relative">
				<h1 class="text-4xl font-extrabold tracking-tight text-white sm:text-5xl md:text-6xl">
					Welcome to RWBYzon
				</h1>
				<p class="mt-6 max-w-3xl text-xl text-indigo-100">
					The premier marketplace for RWBY card trading. Find rare cards, complete your collection,
					and connect with fellow enthusiasts.
				</p>
				<div class="mt-8 flex flex-wrap gap-4">
					<a
						href="/market/listings"
						class="inline-flex items-center justify-center rounded-md border border-transparent bg-white px-5 py-3 text-base font-medium text-indigo-700 shadow-sm transition-colors duration-200 hover:bg-indigo-50"
					>
						Browse Listings
					</a>
					<a
						href="/market/auctions"
						class="inline-flex items-center justify-center rounded-md border border-white px-5 py-3 text-base font-medium text-white transition-colors duration-200 hover:bg-indigo-500"
					>
						View Auctions
					</a>
				</div>
			</div>
		</div>

		<!-- Decorative element -->
		<div class="absolute bottom-0 right-0 h-48 w-48 translate-x-1/3 translate-y-1/3 transform">
			<div class="h-full w-full rounded-full bg-indigo-400 opacity-50 blur-3xl filter"></div>
		</div>
	</section>

	<!-- Latest Auctions Section -->
	<section class="mb-16">
		<h2 class="mb-8 text-2xl font-bold text-gray-900">Latest Auctions</h2>
		<FeaturedAuctions
			auctions={latestAuctions}
			isLoading={isAuctionsLoading}
			error={auctionsError}
			bind:featuredCards={featuredAuctionCards}
			on:select={(event) => {
				selectedAuction = latestAuctions.find((a) => a.id === event.detail.auctionId) || null;
				if (selectedAuction) showAuctionModal = true;
			}}
		/>
	</section>

	<!-- Latest Listings Section -->
	<section class="mb-16">
		<h2 class="mb-8 text-2xl font-bold text-gray-900">Latest Listings</h2>
		<FeaturedListings
			listings={latestListings}
			isLoading={isListingsLoading}
			error={listingsError}
			bind:featuredCards={featuredListingCards}
			on:select={(event) => {
				selectedListing = latestListings.find((l) => l.id === event.detail.listingId) || null;
				if (selectedListing) showListingModal = true;
			}}
		/>
	</section>

	{#if selectedListing && showListingModal}
		<ListingModal
			listing={selectedListing}
			show={showListingModal}
			on:close={() => (showListingModal = false)}
		/>
	{/if}

	{#if selectedAuction && showAuctionModal}
		<AuctionModal
			auction={selectedAuction}
			show={showAuctionModal}
			on:close={() => (showAuctionModal = false)}
		/>
	{/if}
</div>

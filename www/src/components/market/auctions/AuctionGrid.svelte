<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { gsap } from 'gsap';
	import { ScrollTrigger } from 'gsap/ScrollTrigger';
	import type { Auction } from '$lib/types';
	import AuctionCard from './AuctionCard.svelte';

	// Register GSAP plugins
	gsap.registerPlugin(ScrollTrigger);

	// Create an event dispatcher
	const dispatch = createEventDispatcher<{
		bid: { auctionId: string };
	}>();

	export let auctions: Auction[] = [];
	export let loading = false;
	export let hasNextPage = false;
	export let loadMore = () => {}; // Function to load more auctions

	let gridContainer: HTMLElement;

	function createScrollTrigger() {
		ScrollTrigger.create({
			trigger: '#auction-grid',
			start: 'bottom bottom+=200px',
			end: 'bottom bottom',
			onEnter: () => {
				if (!loading && hasNextPage) {
					loadMore();
				}
			},
			once: true
		});
	}

	// When auctions change, recreate the scroll trigger
	$: if (auctions.length > 0) {
		setTimeout(() => {
			createScrollTrigger();
		}, 1000);
	}

	// Handle bid event from AuctionCard
	function handleBid(event: CustomEvent<{ auctionId: string }>) {
		dispatch('bid', event.detail);
	}
</script>

<div
	bind:this={gridContainer}
	id="listing-grid"
	class="container grid grid-cols-1 gap-10 p-4 sm:grid-cols-2 md:grid-cols-2 xl:grid-cols-4"
>
	{#each auctions as auction}
		<div class="h-full">
			<AuctionCard {auction} />
		</div>
	{/each}

	{#if hasNextPage}
		<div class="col-span-full flex justify-center py-4">
			<div
				class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]"
				role="status"
			>
				<span
					class="!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]"
				>
					Loading...
				</span>
			</div>
		</div>
	{/if}

	{#if !hasNextPage && !loading}
		<div class="col-span-full py-8 text-center text-gray-500">No more auctions</div>
	{/if}

	{#if auctions.length === 0 && !loading}
		<div class="col-span-full py-8 text-center text-gray-500">No auctions found</div>
	{/if}
</div>

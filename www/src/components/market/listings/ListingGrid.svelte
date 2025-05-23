<script lang="ts">
	import { gsap } from 'gsap';
	import { ScrollTrigger } from 'gsap/ScrollTrigger';
	import type { Listing } from '$lib/types';
	import ListingCard from './ListingCard.svelte';

	// Register GSAP plugins
	gsap.registerPlugin(ScrollTrigger);

	export let listings: Listing[] = [];
	export let loading = false;
	export let hasNextPage = false;
	export let loadMore = () => {}; // Function to load more listings

	let gridContainer: HTMLElement;

	function createScrollTrigger() {
		ScrollTrigger.create({
			trigger: '#listing-grid',
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

	// when listings change, recreate the scroll trigger
	$: if (listings.length > 0) {
		setTimeout(() => {
			createScrollTrigger();
		}, 1000);
	}
</script>

<div
	bind:this={gridContainer}
	id="listing-grid"
	class="container grid grid-cols-1 gap-10 p-4 sm:grid-cols-2 md:grid-cols-2 xl:grid-cols-4"
>
	{#each listings as listing}
		<div class="h-full">
			<ListingCard {listing} />
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
		<div class="col-span-full py-8 text-center text-gray-500">No more listings</div>
	{/if}

	{#if listings.length === 0 && !loading}
		<div class="col-span-full py-8 text-center text-gray-500">No listings found</div>
	{/if}
</div>

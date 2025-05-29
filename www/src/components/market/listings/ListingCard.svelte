<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { rarityToColor, getCardIconUri, getCardRarityName } from '$lib/card';
	import { gsap } from 'gsap';
	import type { Listing } from '$lib/types';
	import { marketSSE } from '$lib/market';

	export let listing: Listing | null = null;

	// DOM reference for animation
	let cardElement: HTMLElement;

	// Removal state
	let isRemoved = '';
	let removalSubscription: { unsubscribe: () => void } | null = null;

	const dispatch = createEventDispatcher<{
		purchase: { listingId: string };
	}>();

	function openPurchaseModal() {
		// This would trigger your modal opening logic
		// For example, dispatching an event or setting a store value
		const event = new CustomEvent('openModal', {
			detail: { listingId: listing?.id }
		});
		document.dispatchEvent(event);

		dispatch('purchase', { listingId: listing?.id ?? '' });
	}

	// Handler for listing removal
	function handleListingRemoval(removalData: any) {
		if (!listing) return;

		isRemoved = removalData;

		// Animate the card to indicate removal
		if (cardElement?.id === `listing-card-${listing.id}`) {
			gsap.to(cardElement, {
				filter: 'blur(3px) grayscale(80%)',
				opacity: 0.7,
				duration: 0.5,
				ease: 'power2.out'
			});
		}
	}

	// Set up subscription on component mount
	onMount(() => {
		if (listing) {
			removalSubscription = marketSSE.subscribe(
				`listing_${listing.id}_remove`,
				handleListingRemoval
			);
		}
	});

	// Clean up subscription on component destroy
	onDestroy(() => {
		if (removalSubscription) {
			removalSubscription.unsubscribe();
		}
	});
</script>

{#if !listing}
	<div class="card overflow-hidden rounded-lg bg-gray-100 opacity-60 shadow-md">
		<div class="h-48 bg-gray-200"></div>
		<div class="p-4">
			<div class="h-6 w-3/4 rounded bg-gray-200"></div>
			<div class="mt-2 flex items-center justify-between">
				<div class="h-4 w-1/3 rounded bg-gray-200"></div>
				<div class="h-4 w-1/4 rounded bg-gray-200"></div>
			</div>
			<div class="mt-4 h-10 rounded bg-gray-200"></div>
		</div>
	</div>
{:else}
	<!-- Render the listing card -->
	<div
		bind:this={cardElement}
		id="listing-card-{listing.id}"
		class="card relative overflow-hidden rounded-lg bg-white shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<div
			class="relative h-48"
			style="background-image: url('{getCardIconUri(
				listing.edges.card.card_type
			)}'); background-repeat: no-repeat; background-position: center;"
		>
			<div
				class="absolute right-2 top-2 rounded-full px-2 py-1 text-xs font-semibold"
				style="background-color: {rarityToColor(listing.edges.card.rarity)}"
			>
				{getCardRarityName(listing.edges.card.rarity)}
			</div>
		</div>
		<div class="p-4">
			<h3
				class="overflow-hidden text-ellipsis whitespace-nowrap text-lg font-semibold text-gray-900"
				title={listing.edges.card.edges.type.name}
			>
				{listing.edges.card.edges.type.name}
			</h3>
			<div class="mt-2 flex items-center justify-between">
				<span class="text-sm text-gray-500"
					>Value: {listing.edges.card.individual_value.toFixed(2)}%</span
				>
				<span class="font-bold text-indigo-600">{listing.price} Ⱡ</span>
			</div>
			{#if !isRemoved}
				<button
					on:click={openPurchaseModal}
					class="mt-3 flex w-full items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition-colors duration-200 hover:bg-indigo-700"
				>
					<svg
						class="mr-2 h-4 w-4"
						fill="currentColor"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M3 1a1 1 0 000 2h1.22l.305 1.222a.997.997 0 00.01.042l1.358 5.43-.893.892C3.74 11.846 4.632 14 6.414 14H15a1 1 0 000-2H6.414l1-1H14a1 1 0 00.894-.553l3-6A1 1 0 0017 3H6.28l-.31-1.243A1 1 0 005 1H3z"
						></path>
					</svg>
					Purchase
				</button>
			{:else}
				<div
					class="mt-3 flex w-full items-center justify-center rounded-md bg-gray-500 px-4 py-2 text-sm font-medium text-white"
				>
					<svg
						class="mr-2 h-4 w-4"
						fill="currentColor"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							fill-rule="evenodd"
							d="M13.477 14.89A6 6 0 015.11 6.524l8.367 8.368zm1.414-1.414L6.524 5.11a6 6 0 018.367 8.367zM18 10a8 8 0 11-16 0 8 8 0 0116 0z"
							clip-rule="evenodd"
						></path>
					</svg>
					No Longer Available
				</div>
			{/if}
		</div>

		{#if isRemoved == listing?.id}
			<!-- Overlay with "Listing Removed" message -->
			<div class="absolute inset-0 flex items-center justify-center rounded-lg">
				<!-- backdrop -->
				<div class="absolute inset-0 z-10 bg-gray-900 opacity-60"></div>
				<div class="z-20 rounded-lg bg-gray-800 bg-opacity-80 px-4 py-3 text-center shadow-lg">
					<svg
						class="mx-auto mb-2 h-12 w-12 text-gray-300"
						fill="currentColor"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							fill-rule="evenodd"
							d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
							clip-rule="evenodd"
						></path>
					</svg>
					<p class="text-lg font-bold text-white">Listing Removed</p>
					<p class="mt-1 text-sm text-gray-300">This Listing is no longer available</p>
				</div>
			</div>
		{/if}
	</div>
{/if}

<style>
	/* Card hover effects */
	.card {
		transition: transform 0.3s ease;
	}

	.card:hover {
		transform: translateY(-5px);
	}
</style>

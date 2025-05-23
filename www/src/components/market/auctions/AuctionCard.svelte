<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { rarityToColor, getCardIconUri, getCardRarityName } from '$lib/card';
	import { gsap } from 'gsap';
	import type { Auction, AuctionBid } from '$lib/types';
	import { getCurrentPrice, isAuctionEndingSoon, calculateTimeRemaining } from '$lib/auctions';
	import { subscribeToAuctionBids, marketSSE } from '$lib/market';

	export let auction: Auction | null = null;

	const dispatch = createEventDispatcher<{
		bid: { auctionId: string };
	}>();

	// Format the end date
	function formatEndDate(dateString: string): string {
		const date = new Date(dateString);
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		}).format(date);
	}

	// Timer state
	let timeRemaining = auction ? calculateTimeRemaining(auction.ends_at, true) : null;
	let timerElement: HTMLElement;
	let countdownTimer: ReturnType<typeof setInterval>;
	let cardElement: HTMLElement;

	// Bid update subscription
	let bidSubscription: { unsubscribe: () => void } | null = null;

	// Removal state
	let isRemoved = '';
	let removalSubscription: { unsubscribe: () => void } | null = null;

	// Update the timer
	function updateTimer() {
		if (!auction) return;

		const prevTimeRemaining = timeRemaining;
		timeRemaining = calculateTimeRemaining(auction.ends_at, true);

		// Animate timer if it's ending soon
		if (timerElement && timeRemaining && isAuctionEndingSoon(timeRemaining, 5)) {
			// For auctions ending soon, highlight the timer
			if (!prevTimeRemaining || prevTimeRemaining.minutes !== timeRemaining.minutes) {
				gsap.to(timerElement, {
					backgroundColor: 'rgba(234, 88, 12, 0.9)',
					scale: 1.05,
					duration: 0.3,
					yoyo: true,
					repeat: 1,
					ease: 'power1.inOut'
				});
			}
		}

		// If auction has ended, stop timer and change style
		if (timeRemaining?.hasEnded) {
			clearInterval(countdownTimer);
			if (timerElement) {
				gsap.to(timerElement, {
					backgroundColor: 'rgba(239, 68, 68, 0.8)', // red-500 with opacity
					color: 'white',
					duration: 0.5
				});
			}
		}
	}

	// Handler for auction bid updates
	function handleBidUpdate(bidData: AuctionBid) {
		if (!auction) return;

		// Update the auction with the new bid data
		auction = {
			...auction,
			edges: {
				...auction.edges,
				bids: [bidData, ...(auction.edges.bids || [])]
			}
		};

		// Animate the card to indicate new bid
		if (cardElement) {
			gsap.fromTo(
				cardElement,
				{
					boxShadow: '0 0 0 3px rgba(217, 119, 6, 0.8)',
					scale: 1.03
				},
				{
					boxShadow: '0 0 0 0px rgba(217, 119, 6, 0)',
					scale: 1,
					duration: 0.8,
					ease: 'elastic.out(1, 0.3)'
				}
			);
		}
	}

	// Handler for auction removal
	function handleAuctionRemoval(removalData: any) {
		if (!auction) return;

		isRemoved = removalData;

		// Animate the card to indicate removal
		if (cardElement?.id === `auction-card-${auction.id}`) {
			gsap.to(cardElement, {
				filter: 'blur(3px) grayscale(80%)',
				opacity: 0.7,
				duration: 0.5,
				ease: 'power2.out'
			});
		}
	}

	onMount(() => {
		// Initialize timer
		updateTimer();

		// Update every 100ms for smooth countdown
		countdownTimer = setInterval(updateTimer, 100);

		// Subscribe to bid updates and removal events for this specific auction
		if (auction) {
			bidSubscription = subscribeToAuctionBids(auction.id, handleBidUpdate);

			// Subscribe to auction removal events
			removalSubscription = marketSSE.subscribe(
				`auction_${auction.id}_remove`,
				handleAuctionRemoval
			);
		}
	});

	onDestroy(() => {
		// Clean up timer
		clearInterval(countdownTimer);

		// Clean up subscriptions
		if (bidSubscription) {
			bidSubscription.unsubscribe();
		}

		if (removalSubscription) {
			removalSubscription.unsubscribe();
		}
	});

	function openBidModal() {
		const event = new CustomEvent('openModal', {
			detail: { auctionId: auction?.id }
		});
		document.dispatchEvent(event);

		dispatch('bid', { auctionId: auction?.id ?? '' });
	}
</script>

{#if !auction}
	<!-- Skeleton loader for auction card -->
	<div class="card overflow-hidden rounded-lg bg-gray-100 opacity-60 shadow-md">
		<div class="h-48 bg-gray-200"></div>
		<div class="p-4">
			<div class="h-6 w-3/4 rounded bg-gray-200"></div>
			<div class="mt-2 flex items-center justify-between">
				<div class="h-4 w-1/3 rounded bg-gray-200"></div>
				<div class="h-4 w-1/4 rounded bg-gray-200"></div>
			</div>
			<div class="mt-2 h-4 w-1/2 rounded bg-gray-200"></div>
			<div class="mt-4 h-10 rounded bg-gray-200"></div>
		</div>
	</div>
{:else}
	<!-- Render the auction card -->
	<div
		bind:this={cardElement}
		id="auction-card-{auction.id}"
		class="card relative overflow-hidden rounded-lg bg-white shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<div
			class="relative h-48"
			style="background-image: url('{getCardIconUri(
				auction.edges.card.card_type
			)}'); background-repeat: no-repeat; background-position: center;"
		>
			<!-- Card rarity badge -->
			<div
				class="absolute right-2 top-2 rounded-full px-2 py-1 text-xs font-semibold"
				style="background-color: {rarityToColor(auction.edges.card.rarity)}"
			>
				{getCardRarityName(auction.edges.card.rarity)}
			</div>
			<!-- Time remaining badge -->
			<div
				bind:this={timerElement}
				class="absolute bottom-2 right-2 rounded-md bg-black bg-opacity-70 px-2 py-1 text-xs font-medium text-white transition-colors duration-300"
				class:bg-amber-600={timeRemaining &&
					isAuctionEndingSoon(timeRemaining, 5) &&
					!timeRemaining.hasEnded}
				class:bg-red-500={timeRemaining?.hasEnded}
			>
				{timeRemaining?.formatted || 'Loading...'}
			</div>
		</div>
		<div class="p-4">
			<!-- Card name with ellipsis overflow -->
			<h3
				class="overflow-hidden text-ellipsis whitespace-nowrap text-lg font-semibold text-gray-900"
				title={auction.edges.card.edges.type.name}
			>
				{auction.edges.card.edges.type.name}
			</h3>

			<!-- Card stats and current bid -->
			<div class="mt-2 flex items-center justify-between">
				<span class="text-sm text-gray-500"
					>Value: {auction.edges.card.individual_value.toFixed(2)}%</span
				>
				<span class="font-bold text-amber-600">{getCurrentPrice(auction)} Ⱡ</span>
			</div>

			<!-- Auction info -->
			<div class="mt-1 flex items-center justify-between">
				<span class="text-xs text-gray-500">
					Ends: {formatEndDate(auction.ends_at)}
				</span>
				<span class="text-xs font-semibold">
					{auction.edges.bids && auction.edges.bids.length ? auction.edges.bids.length : 'No'} bid{auction
						.edges.bids && auction.edges.bids.length !== 1
						? 's'
						: ''}
				</span>
			</div>

			<!-- Bid button -->
			{#if !isRemoved}
				<button
					on:click={openBidModal}
					class="mt-3 flex w-full items-center justify-center rounded-md bg-amber-600 px-4 py-2 text-sm font-medium text-white transition-colors duration-200 hover:bg-amber-700"
				>
					<svg
						class="mr-2 h-4 w-4"
						fill="currentColor"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							fill-rule="evenodd"
							d="M12 7a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0V8.414l-4.293 4.293a1 1 0 01-1.414 0L8 10.414l-4.293 4.293a1 1 0 01-1.414-1.414l5-5a1 1 0 011.414 0L11 10.586 14.586 7H12z"
							clip-rule="evenodd"
						></path>
					</svg>
					Place Bid
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

		{#if isRemoved == auction?.id}
			<!-- Overlay with "Auction Removed" message -->
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
					<p class="text-lg font-bold text-white">Auction Removed</p>
					<p class="mt-1 text-sm text-gray-300">This auction is no longer available</p>
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

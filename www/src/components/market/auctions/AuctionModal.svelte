<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { gsap } from 'gsap';
	import { rarityToColor, getCardIconUri, getCardRarityName } from '$lib/card';
	import type { Auction, AuctionBid } from '$lib/types';
	import {
		type TimeRemaining,
		calculateMinimumBid,
		calculateTimeRemaining,
		getCurrentPrice
	} from '$lib/auctions';
	import { bidOnAuction, subscribeToAuctionBids, marketSSE } from '$lib/market';

	let { auction, show }: { auction: Auction; show: boolean } = $props();

	const dispatch = createEventDispatcher<{
		close: void;
		bid: { auctionId: string; bidAmount: number };
	}>();

	// User's bid amount
	let bidAmount = $state<string>('');
	let minimumBid = $state<number>(0);
	let bidError = $state<string | null>(null);

	// GSAP animation references
	let modalContainer: HTMLElement | null = $state<HTMLElement | null>(null);
	let modalContent: HTMLElement | null = $state<HTMLElement | null>(null);
	const modalTl = gsap.timeline({ paused: true });

	// Countdown timer
	let timeRemaining: TimeRemaining | null = $state<TimeRemaining | null>(null);
	let countdownTimer: ReturnType<typeof setInterval> | null = $state<ReturnType<
		typeof setInterval
	> | null>(null);
	let timerElement: HTMLElement | null = $state<HTMLElement | null>(null);

	// Subscription for real-time bid updates
	let bidSubscription: { unsubscribe: () => void } | null = null;

	// Subscription for auction removal events
	let removalSubscription: { unsubscribe: () => void } | null = null;

	// Close modal function
	function closeModal() {
		const tl = gsap.timeline({
			onComplete: () => {
				dispatch('close');
			}
		});

		tl.to(modalContent, {
			y: 20,
			opacity: 0,
			scale: 0.95,
			duration: 0.2,
			ease: 'power2.in'
		}).to(
			modalContainer,
			{
				opacity: 0,
				duration: 0.2,
				ease: 'power2.in'
			},
			'-=0.1'
		);
	}

	// Format timeRemaining
	function updateTimeRemaining() {
		if (!auction) return;

		// Get the previous state to check for changes
		const prevTimeRemaining = timeRemaining;

		// Use the utility function to calculate time remaining
		timeRemaining = calculateTimeRemaining(auction.ends_at, true);

		if (timeRemaining.hasEnded) {
			if (countdownTimer) {
				clearInterval(countdownTimer);
				countdownTimer = null; // Clear the timer reference
			}

			// Animation for ended auction
			if (timerElement) {
				gsap.to(timerElement, {
					backgroundColor: 'rgba(239, 68, 68, 0.8)', // red-500 with opacity
					color: 'white',
					duration: 1,
					ease: 'power2.inOut'
				});
			}
			return;
		}

		// Animate the timer when seconds change or when it's the first update
		if (
			timerElement &&
			(!prevTimeRemaining || prevTimeRemaining.formatted !== timeRemaining.formatted)
		) {
			// For the last minute, make the timer more noticeable
			if (timeRemaining.minutes === 0 && timeRemaining.seconds <= 60) {
				gsap.to(timerElement, {
					backgroundColor:
						timeRemaining.seconds % 2 === 0 ? 'rgba(234, 88, 12, 0.9)' : 'rgba(217, 119, 6, 0.9)',
					scale: 1.03,
					duration: 0.3,
					yoyo: true,
					repeat: 1,
					ease: 'power1.inOut'
				});
			} else {
				gsap.to(timerElement, {
					scale: 1.03,
					duration: 0.2,
					yoyo: true,
					repeat: 1,
					ease: 'power1.inOut'
				});
			}
		}
	}

	// Place bid
	function placeBid() {
		const numericBid = Number(bidAmount);
		const bidInput = document.getElementById('bidAmount');

		if (!bidAmount || isNaN(numericBid)) {
			bidError = 'Please enter a valid bid amount';
			// Shake animation for invalid input
			if (bidInput) {
				gsap.to(bidInput, {
					x: -5,
					duration: 0.4,
					ease: 'power2.inOut'
				});
			}
			return;
		}

		if (numericBid < minimumBid) {
			bidError = `Bid must be at least ${minimumBid} Liens`;
			// Shake animation for invalid input
			if (bidInput) {
				gsap.to(bidInput, {
					x: -5,
					duration: 0.4,
					ease: 'power2.inOut'
				});
			}
			return;
		}

		// Success animation before closing
		gsap.to(biddingForm, {
			scale: 1.03,
			duration: 0.2,
			ease: 'back.out(1.5)',
			onComplete: () => {
				// Clear error and dispatch bid event
				bidError = '';
				dispatch('bid', {
					auctionId: auction.id,
					bidAmount: numericBid
				});

				// Close modal after bid
				closeModal();
			}
		});
	}

	// References for animated elements
	let cardDetails: HTMLElement | null = $state<HTMLElement | null>(null);
	let biddingForm: HTMLElement | null = $state<HTMLElement | null>(null);
	let actionButtons: HTMLElement | null = $state<HTMLElement | null>(null);

	// Function to scroll modal into view
	function scrollModalIntoView() {
		if (modalContent) {
			setTimeout(() => {
				modalContent?.scrollIntoView({
					behavior: 'smooth',
					block: 'center'
				});
			}, 50);
		}
	}

	// Handle auction bid updates from SSE
	function handleBidUpdate(bidData: AuctionBid) {
		// Update the auction with new bid data
		auction = {
			...auction,
			edges: {
				...auction.edges,
				bids: [bidData, ...(auction.edges.bids || [])]
			}
		};

		// Recalculate minimum bid based on new data
		minimumBid = calculateMinimumBid(auction);
		bidAmount = minimumBid.toString();
	}

	// Handle auction removal events
	function handleAuctionRemoval(removalData: any) {
		// Show a brief notification to the user before closing
		const notificationTl = gsap.timeline({
			onComplete: () => {
				// Close the modal after showing notification
				closeModal();
			}
		});

		// Create a notification element
		const notification = document.createElement('div');
		notification.className =
			'absolute inset-0 flex items-center justify-center bg-gray-900 bg-opacity-70 z-50';
		notification.innerHTML = `
			<div class="bg-gray-800 bg-opacity-90 p-4 rounded-lg text-center">
				<svg class="mx-auto h-12 w-12 text-red-500 mb-2" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
				</svg>
				<p class="text-lg font-bold text-white">Auction Removed</p>
				<p class="text-gray-300 mt-1">This auction is no longer available</p>
			</div>
		`;

		// Add notification to modal content
		if (modalContent) {
			modalContent.appendChild(notification);

			// Apply animations
			notificationTl
				.from(notification, {
					opacity: 0,
					scale: 0.8,
					duration: 0.3,
					ease: 'back.out(1.7)'
				})
				.to(
					modalContent,
					{
						filter: 'grayscale(80%)',
						duration: 0.5,
						ease: 'power2.out'
					},
					'-=0.1'
				)
				.to(notification, {
					opacity: 0,
					duration: 1,
					delay: 1.5
				});
		} else {
			// If no modal content reference, just close immediately
			closeModal();
		}
	}

	// Set up the component when it mounts
	onMount(() => {
		// Set the suggested bid amount
		minimumBid = calculateMinimumBid(auction);
		bidAmount = minimumBid.toString();

		// Initialize with fresh timer data
		timeRemaining = calculateTimeRemaining(auction.ends_at, true);

		// Start countdown timer (updates every second for smoother experience)
		countdownTimer = setInterval(updateTimeRemaining, 1000);

		// Subscribe to real-time bid updates and removal events
		bidSubscription = subscribeToAuctionBids(auction.id, handleBidUpdate);
		removalSubscription = marketSSE.subscribe(`auction_${auction.id}_remove`, handleAuctionRemoval);

		// Animation setup for modal
		modalTl
			.fromTo(modalContainer, { opacity: 0 }, { opacity: 1, duration: 0.25, ease: 'power2.inOut' })
			.fromTo(
				modalContent,
				{ y: 30, opacity: 0, scale: 0.95 },
				{ y: 0, opacity: 1, scale: 1, duration: 0.35, ease: 'back.out(1.2)' },
				'<0.1'
			)
			.fromTo(
				cardDetails?.children || [],
				{ opacity: 0, y: 10 },
				{
					opacity: 1,
					y: 0,
					stagger: 0.05,
					duration: 0.2,
					ease: 'power1.out'
				},
				'-=0.2'
			)
			.fromTo(biddingForm, { opacity: 0, y: 10 }, { opacity: 1, y: 0, duration: 0.3 }, '-=0.1')
			.fromTo(actionButtons, { opacity: 0, y: 10 }, { opacity: 1, y: 0, duration: 0.2 }, '-=0.1');
	});

	// Clean up when component is destroyed
	onDestroy(() => {
		if (countdownTimer) {
			clearInterval(countdownTimer);
		}
		document.body.style.overflow = ''; // Ensure body scroll is restored

		// Clean up subscriptions
		if (bidSubscription) {
			bidSubscription.unsubscribe();
		}

		if (removalSubscription) {
			removalSubscription.unsubscribe();
		}
	});

	// Handle auction bid
	async function handleBid() {
		try {
			// Place the bid
			await bidOnAuction(auction.id, Number(bidAmount));
			closeModal();
		} catch (err) {
			console.error('Bid failed:', err);
		}
	}

	$effect(() => {
		if (show) {
			document.body.style.overflow = 'hidden'; // Prevent body scroll
			modalTl.restart();
			scrollModalIntoView();
		} else {
			document.body.style.overflow = ''; // Restore body scroll
		}
	});
</script>

{#if show}
	<!-- Modal backdrop -->
	<div
		bind:this={modalContainer}
		class="fixed inset-0 z-50 overflow-y-auto"
		onclick={closeModal}
		onkeydown={(e) => {
			if (e.key === 'Escape') {
				closeModal();
			}
		}}
		aria-labelledby="modal-title"
		role="dialog"
		tabindex="0"
		aria-modal="true"
	>
		<!-- Modal content -->
		<div
			class="flex min-h-screen items-center justify-center px-4 py-20 text-center sm:block sm:p-0"
		>
			<!-- Backdrop -->
			<div class="fixed inset-0 bg-gray-500 opacity-75 transition-opacity"></div>

			<!-- Modal positioning trick -->
			<span class="hidden sm:inline-block sm:h-screen sm:align-middle" aria-hidden="true"
				>&#8203;</span
			>

			<!-- Modal content -->
			<div
				bind:this={modalContent}
				class="relative z-50 inline-block transform overflow-hidden rounded-lg bg-white text-left align-bottom shadow-2xl sm:my-8 sm:w-full sm:max-w-lg sm:align-middle"
				onclick={(e) => e.stopPropagation()}
				onkeydown={(e) => {
					if (e.key === 'Enter') {
						placeBid();
					}
				}}
				role="dialog"
				tabindex="0"
			>
				<!-- Modal header with card image -->
				<div
					class="relative z-50 h-52 bg-center bg-no-repeat"
					style="background-image: url('{getCardIconUri(auction.edges.card.card_type)}');"
				>
					<!-- Close button -->
					<button
						class="absolute right-2 top-2 rounded-full bg-black bg-opacity-50 p-1 text-white hover:bg-opacity-70"
						onclick={closeModal}
						aria-label="Close modal"
					>
						<svg
							class="h-5 w-5"
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
					</button>

					<!-- Card rarity badge -->
					<div
						class="absolute left-2 top-2 rounded-full px-2 py-1 text-xs font-bold text-white"
						style="background-color: {rarityToColor(auction.edges.card.rarity)}"
					>
						{getCardRarityName(auction.edges.card.rarity)}
					</div>

					<!-- Time remaining badge -->
					<div
						bind:this={timerElement}
						class="absolute bottom-2 right-2 rounded-md bg-black bg-opacity-70 px-3 py-1 text-sm font-medium text-white transition-colors"
					>
						{timeRemaining?.formatted || 'Loading...'}
					</div>
				</div>

				<!-- Modal body -->
				<div class="p-6">
					<h3 class="mb-2 text-xl font-bold text-gray-900">{auction.edges.card.edges.type.name}</h3>

					<!-- Card details -->
					<div bind:this={cardDetails} class="mb-4 grid grid-cols-2 gap-2 text-sm">
						<div class="text-gray-600">Value:</div>
						<div class="font-medium">{auction.edges.card.individual_value.toFixed(2)}%</div>

						<div class="text-gray-600">Current Price:</div>
						<div class="font-bold text-amber-600">{getCurrentPrice(auction)} Ⱡ</div>

						<div class="text-gray-600">Total Bids:</div>
						<div class="font-medium">{auction.edges.bids?.length || 0}</div>
					</div>

					<!-- Bidding form -->
					<div bind:this={biddingForm} class="mt-6">
						<label for="bidAmount" class="mb-1 block text-sm font-medium text-gray-700">
							Your Bid (minimum {minimumBid} Ⱡ)
						</label>

						<div class="relative mt-1 rounded-md shadow-sm">
							<input
								type="number"
								name="bidAmount"
								id="bidAmount"
								class="block w-full rounded-md border-gray-300 pr-12 focus:border-amber-500 focus:ring-amber-500 sm:text-sm"
								placeholder="Enter bid amount"
								min={minimumBid}
								bind:value={bidAmount}
							/>
							<div class="absolute inset-y-0 right-0 flex items-center pr-3">
								<span class="text-gray-500 sm:text-sm">Ⱡ</span>
							</div>
						</div>

						{#if bidError}
							<p class="mt-1 text-sm text-red-600">{bidError}</p>
						{/if}
					</div>

					<!-- Action buttons -->
					<div bind:this={actionButtons} class="mt-6 flex justify-end space-x-3">
						<button
							type="button"
							class="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2"
							onclick={closeModal}
						>
							Cancel
						</button>
						<button
							type="button"
							class="rounded-md border border-transparent bg-amber-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2"
							onclick={handleBid}
						>
							Place Bid
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	/* Additional modal styling */
	/* Use :global to style input elements from libraries */
	:global(.modal-open) {
		overflow: hidden;
	}
</style>

<script lang="ts">
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';
	import { gsap } from 'gsap';
	import type { Listing } from '$lib/types';
	import { rarityToColor } from '$lib/card';
	import { notifySuccess } from '$lib/stores/notificationStore';
	import { purchaseListing } from '$lib/market';
	import { marketSSE } from '$lib/sseClient';

	export let listing: Listing;

	export let show = false;

	const dispatch = createEventDispatcher();

	let modalContainer: HTMLElement;
	let modalContent: HTMLElement;

	// Format the individual value to 2 decimal places
	$: formattedIV = listing.edges.card.individual_value.toFixed(2);

	const modalTl = gsap.timeline({ paused: true });

	// Subscription for listing removal events
	let removalSubscription: { unsubscribe: () => void } | null = null;

	onMount(() => {
		// Animation setup for modal
		modalTl
			.fromTo(modalContainer, { opacity: 0 }, { opacity: 1, duration: 0.25 })
			.fromTo(modalContent, { y: 20, opacity: 0 }, { y: 0, opacity: 1, duration: 0.25 }, '<');

		// Subscribe to listing removal events
		removalSubscription = marketSSE.subscribe(`listing_${listing.id}_remove`, handleListingRemoval);

		return () => {
			document.body.style.overflow = ''; // Ensure body scroll is restored
		};
	});

	// Clean up subscriptions when component is destroyed
	onDestroy(() => {
		document.body.style.overflow = ''; // Ensure body scroll is restored

		// Clean up subscription
		if (removalSubscription) {
			removalSubscription.unsubscribe();
		}
	});

	// Handle listing removal events
	function handleListingRemoval(removalData: any) {
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
			<div class="bg-gray-800 bg-opacity-90 p-4 rounded-lg text-center z-70">
				<svg class="mx-auto h-12 w-12 text-red-500 mb-2" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
				</svg>
				<p class="text-lg font-bold text-white">Listing Removed</p>
				<p class="text-gray-300 mt-1">This listing is no longer available</p>
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

	// Function to scroll modal into view
	function scrollModalIntoView() {
		if (modalContent) {
			setTimeout(() => {
				modalContent.scrollIntoView({
					behavior: 'smooth',
					block: 'center'
				});
			}, 50);
		}
	}

	// Watch for changes to the show prop
	$: if (show) {
		document.body.style.overflow = 'hidden'; // Prevent body scroll
		modalTl.play();
		scrollModalIntoView();
	} else {
		document.body.style.overflow = ''; // Restore body scroll
	}

	function closeModal() {
		dispatch('close');
	}

	async function doPurchaseListing() {
		try {
			const resp = await purchaseListing(listing.id);
			if (!resp.meta.success) {
				console.error('Purchase failed:', resp.error);
				return;
			}

			notifySuccess(`Successfully purchased ${listing.edges.card.edges.type.name}!`);
			closeModal();
		} catch (error) {
			console.error('Purchase failed:', error);
		}
	}
</script>

{#if show}
	<div
		bind:this={modalContainer}
		class="fixed inset-0 z-50 overflow-y-auto"
		aria-labelledby="modal-title"
		role="dialog"
		tabindex="0"
		aria-modal="true"
		on:click={closeModal}
		on:keydown={(e) => {
			if (e.key === 'Escape') {
				closeModal();
			}
		}}
	>
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
				on:click|stopPropagation={() => {}}
				on:keydown={(e) => {
					if (e.key === 'Enter') {
						doPurchaseListing();
					}
				}}
				role="dialog"
				tabindex="0"
			>
				<div class="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
					<div class="sm:flex sm:items-start">
						<div class="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
							<h3
								class="flex flex-row items-center gap-2 text-lg font-medium leading-6 text-gray-900"
								id="modal-title"
							>
								<svg class="h-4 w-4" viewBox="0 0 10 10" xmlns="http://www.w3.org/2000/svg">
									<circle
										fill={rarityToColor(listing.edges.card.rarity)}
										cx="5"
										cy="5"
										r="4"
										stroke="black"
										stroke-width="0.25"
									></circle>
								</svg>
								{listing.edges.card.edges.type.name} ({formattedIV}%)
							</h3>
							<div class="mt-4 text-gray-700">
								<p>
									Are you sure you want to purchase this card for <span
										class="font-semibold text-indigo-600">{listing.price} Ⱡ</span
									>?
								</p>
								<p class="mt-2 text-sm text-gray-500">This action cannot be undone.</p>
							</div>
						</div>
					</div>
				</div>
				<div class="flex flex-row-reverse gap-4 bg-gray-50 px-4 py-3 sm:gap-3 sm:px-6">
					<button
						on:click={closeModal}
						type="button"
						class="inline-flex w-auto justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm transition-colors duration-200 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:mt-0 sm:w-auto sm:text-sm"
					>
						Close
					</button>
					<button
						on:click={doPurchaseListing}
						type="button"
						class="inline-flex w-auto items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm transition-colors duration-200 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
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
				</div>
			</div>
		</div>
	</div>
{/if}

<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Listing } from '$lib/types';
	import ListingCard from './ListingCard.svelte';

	export let listings: Listing[] = [];
	export let isLoading = false;
	export let error: string | null = null;
	export let featuredCards: HTMLElement | undefined = undefined;

	const dispatch = createEventDispatcher<{
		select: { listingId: string };
	}>();

	function handleSelect(event: CustomEvent) {
		dispatch('select', event.detail);
	}
</script>

{#if isLoading}
	<div class="flex justify-center py-12">
		<div
			class="h-12 w-12 animate-spin rounded-full border-4 border-indigo-600 border-t-transparent"
		></div>
	</div>
{:else if error}
	<div class="rounded-md bg-red-50 p-4">
		<div class="flex">
			<div class="flex-shrink-0">
				<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
					<path
						fill-rule="evenodd"
						d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
						clip-rule="evenodd"
					/>
				</svg>
			</div>
			<div class="ml-3">
				<p class="text-sm text-red-700">{error}</p>
			</div>
		</div>
	</div>
{:else}
	<div class="overflow-x-auto pb-4">
		<div
			bind:this={featuredCards}
			class="flex gap-6 md:grid md:grid-cols-3 lg:grid-cols-5"
			style="min-width: min-content;"
		>
			{#each listings as listing}
				<div class="w-56 flex-shrink-0 md:w-auto">
					<ListingCard {listing} on:purchase={handleSelect} />
				</div>
			{/each}

			{#if listings.length < 5}
				{@const placeholders = Array(5 - listings.length).fill(0)}
				{#each placeholders as _}
					<div class="w-56 flex-shrink-0 md:w-auto">
						<ListingCard />
					</div>
				{/each}
			{/if}
		</div>

		<!-- Scroll indicator for mobile -->
		<div class="mt-4 flex justify-center md:hidden">
			<div class="flex space-x-2">
				<div class="h-2 w-2 rounded-full bg-gray-300"></div>
				<div class="h-2 w-6 rounded-full bg-indigo-600"></div>
				<div class="h-2 w-2 rounded-full bg-gray-300"></div>
			</div>
		</div>
	</div>
{/if}

<script lang="ts">
	import TradeCard from './TradeCard.svelte';
	import type { Card } from '$lib/card';

	let {
		cards,
		loading,
		maxHeight,
		cardSelect
	}: {
		cards: Card[];
		loading: boolean;
		maxHeight: string;
		cardSelect: (cardId: string) => void;
	} = $props();

	// Event handler for card selection
	function handleCardSelect(cardId: string) {
		// Forward the event to parent components
		cardSelect(cardId);
	}
</script>

<div class="trade-grid-container" style="max-height: {maxHeight};">
	{#if loading}
		<div class="flex h-full w-full items-center justify-center">
			<div class="loading-spinner"></div>
		</div>
	{:else if cards.length === 0}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<svg
				class="mb-4 h-12 w-12 text-gray-400"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
				></path>
			</svg>
			<p class="text-center text-gray-500">No cards to be shown yet</p>
		</div>
	{:else}
		<div class="grid-layout">
			{#each cards as card (card.id)}
				<div class="grid-item">
					<TradeCard {card} select={handleCardSelect} />
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.trade-grid-container {
		width: 100%;
		overflow-y: auto;
		scrollbar-width: thin;
		scrollbar-color: rgba(156, 163, 175, 0.5) transparent;
	}

	/* For Webkit browsers like Chrome/Safari */
	.trade-grid-container::-webkit-scrollbar {
		width: 6px;
	}

	.trade-grid-container::-webkit-scrollbar-track {
		background: transparent;
	}

	.trade-grid-container::-webkit-scrollbar-thumb {
		background-color: rgba(156, 163, 175, 0.5);
		border-radius: 6px;
	}

	.grid-layout {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
		gap: 12px;
		padding: 8px;
		width: 100%;
	}

	.grid-item {
		display: flex;
		justify-content: center;
	}

	/* Loading spinner animation */
	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 4px solid rgba(156, 163, 175, 0.3);
		border-radius: 50%;
		border-top-color: #6366f1; /* Indigo 500 */
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	@media (max-width: 640px) {
		.grid-layout {
			grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
			gap: 8px;
		}
	}
</style>

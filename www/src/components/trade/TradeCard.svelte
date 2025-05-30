<script lang="ts">
	import { rarityToColor, getCardIconUri, getCardRarityName } from '$lib/card';
	import type { Card } from '$lib/card';

	export let card: Card | null = null;
	export let select: (cardId: string) => void = () => {};

	// DOM reference for animation
	let cardElement: HTMLElement;
</script>

{#if !card}
	<div class="card overflow-hidden rounded-lg bg-gray-100 opacity-60 shadow-md">
		<div class="h-32 bg-gray-200"></div>
		<div class="p-2">
			<div class="h-5 w-3/4 rounded bg-gray-200"></div>
			<div class="mt-2 flex items-center justify-between">
				<div class="h-3 w-1/3 rounded bg-gray-200"></div>
				<div class="h-3 w-1/4 rounded bg-gray-200"></div>
			</div>
			<div class="mt-3 h-8 rounded bg-gray-200"></div>
		</div>
	</div>
{:else}
	<!-- Render the card -->
	<div
		bind:this={cardElement}
		id="card-{card.id}"
		class="card relative overflow-hidden rounded-lg bg-white shadow-md transition-shadow duration-300 hover:shadow-lg"
	>
		<div
			class="relative h-32"
			style="background-image: url('{getCardIconUri(
				card.card_type
			)}'); background-repeat: no-repeat; background-position: center; background-size: contain;"
		>
			<div
				class="absolute right-2 top-1 rounded-full px-1.5 py-0.5 text-xs font-semibold"
				style="background-color: {rarityToColor(card.rarity)}"
			>
				{getCardRarityName(card.rarity)}
			</div>
		</div>
		<div class="p-2">
			<h3
				class="overflow-hidden text-ellipsis whitespace-nowrap text-sm font-semibold text-gray-900"
				title={card.edges.type.name}
			>
				{card.edges.type.name}
			</h3>
			<div class="mt-1 flex items-center justify-between">
				<span class="text-xs text-gray-500">IV: {card.individual_value.toFixed(2)}%</span>
				<span class="text-xs text-gray-500">Lvl: {card.level}</span>
			</div>
			<button
				on:click={() => {
					select(card.id);
				}}
				class="mt-2 flex w-full items-center justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white transition-colors duration-200 hover:bg-indigo-700"
			>
				<svg
					class="mr-1.5 h-3 w-3"
					fill="currentColor"
					viewBox="0 0 20 20"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path d="M10 12a2 2 0 100-4 2 2 0 000 4z"></path>
					<path
						fill-rule="evenodd"
						d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
						clip-rule="evenodd"
					></path>
				</svg>
				Select
			</button>
		</div>
	</div>
{/if}

<style>
	/* Card hover effects */
	.card {
		transition: transform 0.3s ease;
		max-width: 180px;
		width: 100%;
	}

	.card:hover {
		transform: translateY(-3px);
	}
</style>

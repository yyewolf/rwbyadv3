<script lang="ts">
	import type { Card } from '$lib/card';
	import { getMyCards, getPlayerCards, createTrade } from '$lib/trades';
	import { onMount } from 'svelte';
	import { notifySuccess, notifyError } from '$lib/stores/notificationStore';
	import TradeGrid from '../../components/trade/TradeGrid.svelte';
	import { fly } from 'svelte/transition';

	const urlParams = new URLSearchParams(window.location.search);
	const playerId = urlParams.get('playerId') ?? '';
	const playerName = urlParams.get('playerName') ?? 'Player';

	let myCards: Card[] = $state([]);
	let playerCards: Card[] = $state([]);
	let loadingMyCards = $state(true);
	let loadingPlayerCards = $state(true);
	let selectedMyCards: Card[] = $state([]);
	let selectedPlayerCards: Card[] = $state([]);
	let isSubmitting = $state(false);

	onMount(async () => {
		try {
			myCards = await getMyCards();
		} finally {
			loadingMyCards = false;
		}

		if (playerId) {
			try {
				playerCards = await getPlayerCards(playerId);
			} finally {
				loadingPlayerCards = false;
			}
		} else {
			loadingPlayerCards = false;
		}
	});

	function moveFromAtoB(source: Card[], target: Card[], cardId: string): Card[][] {
		const card = source.find((c) => c.id === cardId);
		if (card) {
			source = source.filter((c) => c.id !== cardId);
			target = [...target, card];
		}
		return [source, target];
	}

	function handleMyCardSelect(cardId: string) {
		[myCards, selectedMyCards] = moveFromAtoB(myCards, selectedMyCards, cardId);
	}

	function handlePlayerCardSelect(cardId: string) {
		[playerCards, selectedPlayerCards] = moveFromAtoB(playerCards, selectedPlayerCards, cardId);
	}

	async function proposeTrade() {
		if (!playerId || (selectedMyCards.length === 0 && selectedPlayerCards.length === 0)) {
			notifyError('Invalid trade: You must select at least one card to trade');
			return;
		}

		isSubmitting = true;
		try {
			let resp = await createTrade(playerId, selectedMyCards, selectedPlayerCards);
			if (resp.error) {
				return;
			}
			notifySuccess(`Trade offer sent to ${playerName}!`);

			// Reset the selected cards
			myCards = [...myCards, ...selectedMyCards];
			playerCards = [...playerCards, ...selectedPlayerCards];
			selectedMyCards = [];
			selectedPlayerCards = [];
		} catch (error) {
			notifyError('Failed to submit trade offer. Please try again.');
			console.error('Trade error:', error);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Trading | RWBY Adventure</title>
</svelte:head>

<div class="container mx-auto px-4 py-6">
	<h1 class="mb-6 text-2xl font-bold text-gray-800">Trading Cards</h1>

	<div class="flex flex-col gap-6 md:flex-row">
		<!-- My Cards Section -->
		<div class="flex-1 rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
			<h2 class="mb-4 text-xl font-semibold">My Cards ({myCards.length})</h2>
			<TradeGrid
				cards={myCards}
				loading={loadingMyCards}
				maxHeight="40vh"
				cardSelect={handleMyCardSelect}
			/>
		</div>

		<!-- Player Cards Section (Only show if playerid is provided) -->
		{#if playerId}
			<div class="flex-1 rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
				<h2 class="mb-4 text-xl font-semibold">Their Cards ({playerCards.length})</h2>
				<TradeGrid
					cards={playerCards}
					loading={loadingPlayerCards}
					maxHeight="40vh"
					cardSelect={handlePlayerCardSelect}
				/>
			</div>
		{/if}
	</div>

	<!-- Trade Offer Section -->
	<div class="mt-8">
		<h2 class="mb-4 text-xl font-bold text-gray-800">Trade Offer</h2>
		<div class="flex flex-col gap-6 md:flex-row">
			<!-- My Offer Section -->
			<div class="flex-1 rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
				<div class="flex items-center justify-between">
					<h3 class="mb-4 text-lg font-semibold text-gray-800">
						My Offer ({selectedMyCards.length})
					</h3>
					{#if selectedMyCards.length > 0}
						<button
							class="mb-4 rounded bg-gray-200 px-3 py-1 text-sm text-gray-700 hover:bg-gray-300"
							onclick={() => {
								myCards = [...myCards, ...selectedMyCards];
								selectedMyCards = [];
							}}
						>
							Clear All
						</button>
					{/if}
				</div>
				<TradeGrid
					cards={selectedMyCards}
					loading={false}
					maxHeight="30vh"
					cardSelect={(cardId) => {
						[selectedMyCards, myCards] = moveFromAtoB(selectedMyCards, myCards, cardId);
					}}
				/>
			</div>

			<!-- I Receive Section (Only show if playerid is provided) -->
			{#if playerId}
				<div class="flex-1 rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
					<div class="flex items-center justify-between">
						<h3 class="mb-4 text-lg font-semibold text-gray-800">
							I Receive ({selectedPlayerCards.length})
						</h3>
						{#if selectedPlayerCards.length > 0}
							<button
								class="mb-4 rounded bg-gray-200 px-3 py-1 text-sm text-gray-700 hover:bg-gray-300"
								onclick={() => {
									playerCards = [...playerCards, ...selectedPlayerCards];
									selectedPlayerCards = [];
								}}
							>
								Clear All
							</button>
						{/if}
					</div>
					<TradeGrid
						cards={selectedPlayerCards}
						loading={false}
						maxHeight="30vh"
						cardSelect={(cardId) => {
							[selectedPlayerCards, playerCards] = moveFromAtoB(
								selectedPlayerCards,
								playerCards,
								cardId
							);
						}}
					/>
				</div>
			{/if}
		</div>

		<!-- Submit Trade Button -->
		{#if playerId && (selectedMyCards.length > 0 || selectedPlayerCards.length > 0)}
			<div class="mt-6 flex justify-center">
				<button
					class="group relative inline-flex items-center justify-center overflow-hidden rounded-md bg-gradient-to-br from-purple-600 to-indigo-600 p-0.5 text-sm font-medium text-white hover:text-white focus:outline-none focus:ring-4 focus:ring-indigo-300 disabled:opacity-70"
					onclick={proposeTrade}
					disabled={isSubmitting}
				>
					<span
						class="relative rounded-md bg-[rgba(0,0,0,0.2)] px-8 py-2.5 transition-all duration-75 ease-in group-hover:bg-opacity-0"
					>
						{#if isSubmitting}
							<span class="flex items-center">
								<svg class="mr-2 h-4 w-4 animate-spin" viewBox="0 0 24 24">
									<circle
										class="opacity-25"
										cx="12"
										cy="12"
										r="10"
										stroke="currentColor"
										stroke-width="4"
									></circle>
									<path
										class="opacity-75"
										fill="currentColor"
										d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
									></path>
								</svg>
								Processing...
							</span>
						{:else}
							Propose Trade to {playerName}
						{/if}
					</span>
				</button>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Add any page-specific styles here */
</style>

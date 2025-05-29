<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let filters = {
		minPrice: '',
		maxPrice: '',
		sortBy: 'ending_soon' // Default sort by ending soon
	};

	const dispatch = createEventDispatcher();

	// Sorting options
	const sortOptions = [
		{ value: 'ending_soon', label: 'Ending Soon' },
		{ value: 'newest', label: 'Newest First' },
		{ value: 'price_low', label: 'Price: Low to High' },
		{ value: 'price_high', label: 'Price: High to Low' },
		{ value: 'bids', label: 'Most Bids' }
	];

	function updateFilters() {
		dispatch('filter', filters);
	}

	function clearFilters() {
		filters = {
			minPrice: '',
			maxPrice: '',
			sortBy: 'ending_soon'
		};
		dispatch('filter', filters);
	}
</script>

<div class="mb-8 rounded-lg bg-white p-4 shadow-md">
	<h3 class="mb-4 text-lg font-semibold text-gray-800">Filter & Sort</h3>

	<div class="mb-4 grid grid-cols-1 gap-4 md:grid-cols-3">
		<!-- Price Range -->
		<div class="flex flex-col">
			<label for="min-price" class="mb-1 block text-sm font-medium text-gray-700">
				Min Price
			</label>
			<input
				type="number"
				id="min-price"
				bind:value={filters.minPrice}
				on:change={updateFilters}
				placeholder="0"
				min="0"
				class="rounded-md border border-gray-300 p-2 focus:border-amber-500 focus:outline-none focus:ring-amber-500"
			/>
		</div>

		<div class="flex flex-col">
			<label for="max-price" class="mb-1 block text-sm font-medium text-gray-700">
				Max Price
			</label>
			<input
				type="number"
				id="max-price"
				bind:value={filters.maxPrice}
				on:change={updateFilters}
				placeholder="No limit"
				min="0"
				class="rounded-md border border-gray-300 p-2 focus:border-amber-500 focus:outline-none focus:ring-amber-500"
			/>
		</div>

		<!-- Sort Options -->
		<div class="flex flex-col">
			<label for="sort-by" class="mb-1 block text-sm font-medium text-gray-700"> Sort By </label>
			<select
				id="sort-by"
				bind:value={filters.sortBy}
				on:change={updateFilters}
				class="rounded-md border border-gray-300 p-2 focus:border-amber-500 focus:outline-none focus:ring-amber-500"
			>
				{#each sortOptions as option}
					<option value={option.value}>{option.label}</option>
				{/each}
			</select>
		</div>
	</div>

	<!-- Action Buttons -->
	<div class="flex justify-end">
		<button
			on:click={clearFilters}
			class="rounded-md bg-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500"
		>
			Clear Filters
		</button>
	</div>
</div>

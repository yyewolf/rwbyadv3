<script lang="ts">
	import { onMount } from 'svelte';
	import { gsap } from 'gsap';

	export let title: string;
	export let searchQuery = '';
	export let placeholderText = 'Begin Typing To Search...';
	export let accentColor = 'indigo';

	// Event dispatcher for search input
	function handleSearchInput(event: Event) {
		const inputEvent = new CustomEvent('search', {
			detail: { query: (event.target as HTMLInputElement).value }
		});

		// Dispatch the event so parent component can handle it
		document.dispatchEvent(inputEvent);
	}

	onMount(() => {
		// Animate the header
		gsap.from('.market-header', {
			opacity: 0,
			y: -20,
			duration: 0.7,
			ease: 'power2.out'
		});

		// Animate the search input
		gsap.from('.search-container', {
			opacity: 0,
			y: 20,
			duration: 0.7,
			delay: 0.2,
			ease: 'power2.out'
		});
	});
</script>

<div class="container mx-auto px-12 lg:px-32">
	<div
		class="market-header flex w-full flex-row items-center justify-between p-4 text-2xl font-bold"
	>
		<span>{title}</span>
	</div>

	<div class="search-container text-md flex w-full flex-row justify-center p-4 font-bold">
		<label for="search" class="sr-only">Search</label>
		<input
			type="search"
			value={searchQuery}
			on:input={handleSearchInput}
			placeholder={placeholderText}
			class={`w-full max-w-lg rounded-md border-2 border-gray-300 p-2 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-${accentColor}-600`}
		/>
	</div>

	<slot />
</div>

<style>
	/* Add any additional styling specific to the browse layout here */
</style>

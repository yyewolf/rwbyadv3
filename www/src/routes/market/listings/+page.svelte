<script lang="ts">
	import ListingGrid from '../../../components/market/listings/ListingGrid.svelte';
	import ListingModal from '../../../components/market/listings/ListingModal.svelte';
	import MarketBrowseLayout from '../../../components/MarketBrowseLayout.svelte';

	import { onMount } from 'svelte';
	import { debounce } from '../../../utils/debounce';
	import { getListings } from '$lib/market';
	import type { Listing } from '$lib/types';

	// State variables
	let listings: Listing[] = [];
	let loading = false;
	let hasNextPage = false;
	let currentPage = 1;
	let searchQuery = '';
	let selectedListing: Listing | null = null;
	let showModal = false;
	let error: string | null = null;

	// Set up debounced search
	const debouncedSearch = debounce((query) => {
		getListings(searchQuery, 1)
			.then((response) => {
				listings = response.data;
				hasNextPage = response.pagination.has_next;
				currentPage = 1;
				loading = false;
			})
			.catch((error) => {
				console.error('Error fetching listings:', error);
				loading = false;
			});
	}, 500);

	onMount(() => {
		// Initial listings load
		loading = true;
		getListings(searchQuery, 1)
			.then((response) => {
				listings = response.data;
				hasNextPage = response.pagination.has_next;
				currentPage = 1;
				loading = false;
			})
			.catch((err) => {
				console.error('Error fetching listings:', err);
				error = err instanceof Error ? err.message : 'Failed to fetch listings';
				loading = false;
			});

		// Set up event listener for search from MarketBrowseLayout
		const searchHandler = (event: CustomEvent) => {
			searchQuery = event.detail.query;
			debouncedSearch(searchQuery);
		};

		// @ts-ignore
		document.addEventListener('search', searchHandler);

		return () => {
			// @ts-ignore
			document.removeEventListener('search', searchHandler);
		};
	});

	// Load more listings for infinite scroll
	function loadMore() {
		if (!loading && hasNextPage) {
			loading = true;
			console.log('Loading more listings...');
			getListings(searchQuery, currentPage + 1)
				.then((response) => {
					listings = [...listings, ...response.data];
					hasNextPage = response.pagination.has_next;
					currentPage += 1;
					loading = false;
				})
				.catch((err) => {
					console.error('Error loading more listings:', err);
					loading = false;
					error = err instanceof Error ? err.message : 'Failed to load more listings';
				});
		}
	}

	// Open modal with selected listing
	function openListingModal(event: CustomEvent) {
		const listingId = event.detail.listingId;
		selectedListing = listings.find((l) => l.id === listingId) ?? null;

		if (selectedListing) {
			showModal = true;
		}
	}

	// Listen for custom events from Card components
	onMount(() => {
		document.addEventListener('openModal', (e) => openListingModal(e as CustomEvent));

		return () => {
			document.removeEventListener('openModal', (e) => openListingModal(e as CustomEvent));
		};
	});
</script>

<MarketBrowseLayout
	title="Listings"
	{searchQuery}
	placeholderText="Begin Typing To Search Listings..."
	accentColor="indigo"
>
	<ListingGrid {listings} {loading} {hasNextPage} {loadMore} />

	{#if selectedListing && showModal}
		<ListingModal listing={selectedListing} show={showModal} on:close={() => (showModal = false)} />
	{/if}
</MarketBrowseLayout>

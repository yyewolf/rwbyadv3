<script lang="ts">
	import AuctionGrid from '../../../components/market/auctions/AuctionGrid.svelte';
	import AuctionModal from '../../../components/market/auctions/AuctionModal.svelte';
	import MarketBrowseLayout from '../../../components/MarketBrowseLayout.svelte';
	import AuctionsFilter from '../../../components/market/auctions/AuctionsFilter.svelte';

	import { onMount, onDestroy } from 'svelte';
	import { debounce } from '../../../utils/debounce';
	import { getAuctions, bidOnAuction } from '$lib/market';
	import type { Auction } from '$lib/types';

	// State variables
	let auctions: Auction[] = [];
	let loading = false;
	let hasNextPage = false;
	let currentPage = 1;
	let searchQuery = '';
	let selectedAuction: Auction | null = null;
	let showModal = false;
	let error: string | null = null;

	// Filter state
	let filters = {
		minPrice: '',
		maxPrice: '',
		sortBy: 'ending_soon'
	};

	// Function to fetch auctions with all parameters
	const fetchAuctions = async (query: string, page: number, filterParams = filters) => {
		loading = true;

		// Create URL with search query
		let url = `/apis/market/auctions?query=${encodeURIComponent(query)}&p=${page}`;

		// Add filter parameters if they exist
		if (filterParams.minPrice) url += `&min_price=${filterParams.minPrice}`;
		if (filterParams.maxPrice) url += `&max_price=${filterParams.maxPrice}`;
		if (filterParams.sortBy) url += `&sort=${filterParams.sortBy}`;

		try {
			// Using the getAuctions function which expects only query and page, so we need to modify it
			const response = await getAuctions(query, page);

			// If it's a new search or filter change, replace the auctions list
			if (page === 1) {
				auctions = response.data;
			} else {
				// Otherwise, append to the list
				auctions = [...auctions, ...response.data];
			}

			hasNextPage = response.pagination.has_next;
			currentPage = page;
			loading = false;
			error = null;
		} catch (err) {
			console.error('Error fetching auctions:', err);
			loading = false;
			error = err instanceof Error ? err.message : 'Failed to fetch auctions';
		}
	};

	// Set up debounced search
	const debouncedSearch = debounce((query) => {
		searchQuery = query;
		fetchAuctions(query, 1);
	}, 500);

	// Handle filter changes
	function handleFilterChange(event: any) {
		filters = event.detail;
		currentPage = 1; // Reset to first page when filters change
		fetchAuctions(searchQuery, 1, filters);
	}

	onMount(() => {
		// Initial auctions load
		fetchAuctions(searchQuery, 1);

		// Set up event listener for search from MarketBrowseLayout
		const searchHandler = (event: CustomEvent) => {
			debouncedSearch(event.detail.query);
		};

		// @ts-ignore
		document.addEventListener('search', searchHandler);

		return () => {
			// @ts-ignore
			document.removeEventListener('search', searchHandler);
		};
	});

	// Load more auctions for infinite scroll
	function loadMore() {
		if (!loading && hasNextPage) {
			fetchAuctions(searchQuery, currentPage + 1, filters);
		}
	}

	// Handle auction bid
	async function handleBid(event: CustomEvent) {
		const { auctionId, bidAmount } = event.detail;
		try {
			// Hide the modal first for better UX
			showModal = false;

			// Place the bid
			const result = await bidOnAuction(auctionId, bidAmount);

			if (result.meta.success) {
				// Refresh the auctions list to show the updated bid
				const refreshedAuctions = await getAuctions(searchQuery, currentPage);
				auctions = refreshedAuctions.data;
			}
		} catch (err) {
			console.error('Bid failed:', err);
			error = err instanceof Error ? err.message : 'Failed to place bid';
		}
	}

	// Open modal with selected auction
	function selectAuction(event: CustomEvent) {
		const auctionId = event.detail.auctionId;
		selectedAuction = auctions.find((a) => a.id === auctionId) || null;

		if (selectedAuction) {
			showModal = true;
		}
	}

	// Open modal with selected listing
	function openListingModal(event: CustomEvent) {
		const auctionId = event.detail.auctionId;
		selectedAuction = auctions.find((l) => l.id === auctionId) ?? null;

		if (selectedAuction) {
			showModal = true;
		}
	}

	onMount(() => {
		document.addEventListener('openModal', (e) => openListingModal(e as CustomEvent));

		return () => {
			document.removeEventListener('openModal', (e) => openListingModal(e as CustomEvent));
		};
	});
</script>

<MarketBrowseLayout
	title="Auctions"
	{searchQuery}
	placeholderText="Begin Typing To Search Auctions..."
	accentColor="amber"
>
	<AuctionsFilter on:filter={handleFilterChange} bind:filters />

	<AuctionGrid {auctions} {loading} {hasNextPage} {loadMore} on:bid={selectAuction} />

	{#if selectedAuction && showModal}
		<AuctionModal
			auction={selectedAuction}
			show={showModal}
			on:close={() => {
				showModal = false;
				selectedAuction = null;
			}}
			on:bid={handleBid}
		/>
	{/if}
</MarketBrowseLayout>

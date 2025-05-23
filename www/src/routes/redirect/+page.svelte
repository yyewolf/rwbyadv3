<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	// Define allowed redirect destinations
	const allowedRedirects: Record<string, string> = {
		market: '/market',
		home: '/',
		collection: '/collection',
		deck: '/deck'
	};

	// Loading state
	let loading = true;
	let error: string | null = null;
	let destination: string = '';
	let countdown: number = 3;
	let timer: ReturnType<typeof setInterval>;

	onMount(() => {
		try {
			// Get the 'to' parameter from the URL
			const to = $page.url.searchParams.get('to');

			// Check if parameter exists and is in our allowed list
			if (to && allowedRedirects[to]) {
				// Set destination name for display
				destination = to.charAt(0).toUpperCase() + to.slice(1);

				// Start countdown
				timer = setInterval(() => {
					countdown--;
					if (countdown <= 0) {
						clearInterval(timer);
						goto(allowedRedirects[to]);
					}
				}, 1000);
			} else {
				// Handle invalid or missing redirect parameter
				error = 'Invalid redirect destination';
				// Default fallback - redirect to home after a delay
				timer = setInterval(() => {
					countdown--;
					if (countdown <= 0) {
						clearInterval(timer);
						goto('/');
					}
				}, 1000);
			}
		} catch (e) {
			error = 'An error occurred during redirection';
			console.error(e);

			// Default fallback with countdown
			timer = setInterval(() => {
				countdown--;
				if (countdown <= 0) {
					clearInterval(timer);
					goto('/');
				}
			}, 1000);
		} finally {
			loading = false;
		}

		return () => {
			if (timer) clearInterval(timer);
		};
	});
</script>

<div class="flex h-screen w-full items-center justify-center bg-gray-100">
	<div
		class="mx-4 w-full max-w-md rounded-lg border border-gray-200 bg-white p-8 text-center shadow-lg"
	>
		{#if loading}
			<div class="mb-4 text-2xl font-semibold text-gray-800">Preparing your redirect...</div>
			<div
				class="mx-auto h-16 w-16 animate-spin rounded-full border-4 border-b-blue-700 border-l-blue-600 border-r-blue-600 border-t-blue-500"
			></div>
		{:else if error}
			<div class="mb-6">
				<div class="mb-2 text-4xl text-red-600">⚠️</div>
				<div class="mb-4 text-2xl font-semibold text-red-600">{error}</div>
				<p class="mb-4 text-gray-600">
					Redirecting to home page in <span class="font-bold text-blue-600">{countdown}</span> seconds...
				</p>
				<div class="h-2.5 w-full rounded-full bg-gray-200">
					<div
						class="h-2.5 rounded-full bg-blue-600 transition-all duration-1000"
						style="width: {(countdown / 3) * 100}%"
					></div>
				</div>
			</div>
		{:else}
			<div class="mb-6">
				<div class="mb-2 text-4xl text-blue-600">✓</div>
				<div class="mb-4 text-2xl font-semibold text-gray-800">
					Redirecting you to {destination}
				</div>
				<p class="mb-4 text-gray-600">
					You will be redirected in <span class="font-bold text-blue-600">{countdown}</span> seconds...
				</p>
				<div class="h-2.5 w-full rounded-full bg-gray-200">
					<div
						class="h-2.5 rounded-full bg-blue-600 transition-all duration-1000"
						style="width: {(countdown / 3) * 100}%"
					></div>
				</div>
			</div>
			<div
				class="mx-auto h-16 w-16 animate-spin rounded-full border-4 border-b-blue-700 border-l-blue-600 border-r-blue-600 border-t-blue-500"
			></div>
		{/if}
	</div>
</div>

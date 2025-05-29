<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import gsap from 'gsap';

	// Define error types
	const errorTypes: Record<
		string,
		{ code: string; title: string; message: string; hint: string; icon: string }
	> = {
		'404': {
			code: '404',
			title: 'Not Found',
			message: 'The resource you requested could not be found.',
			hint: 'Please check the URL and try again.',
			icon: 'question-mark'
		},
		'403': {
			code: '403',
			title: 'Forbidden',
			message: 'You do not have permission to access this resource.',
			hint: 'Please log in or contact support if you believe this is an error.',
			icon: 'lock'
		},
		'401': {
			code: '401',
			title: 'Unauthorized',
			message: 'Authentication is required to access this resource.',
			hint: 'Please log in to continue.',
			icon: 'key'
		},
		'500': {
			code: '500',
			title: 'Internal Server Error',
			message: 'An unexpected error occurred on the server.',
			hint: 'Please try again in a few seconds or contact support.',
			icon: 'server'
		},
		dungeon_not_found: {
			code: '404',
			title: 'Dungeon Not Found',
			message: 'The dungeon you requested could not be found.',
			hint: 'Check if the dungeon ID is correct or try accessing your dungeons list.',
			icon: 'dungeon'
		},
		auth_error: {
			code: '403',
			title: 'Authentication Error',
			message: 'There was a problem with the authentication process.',
			hint: 'Try again in a few seconds or use a different login method.',
			icon: 'user'
		},
		github_star_error: {
			code: 'GitHub',
			title: 'GitHub Star Error',
			message: 'There was a problem verifying your GitHub star.',
			hint: 'Make sure you have starred the repository and try again.',
			icon: 'github'
		},
		trade_error: {
			code: '404',
			title: 'Trade Not Found',
			message: 'The trade you requested could not be found or has been completed.',
			hint: 'Check if the trade ID is correct or try accessing your trades list.',
			icon: 'exchange'
		},
		market_error: {
			code: '404',
			title: 'Market Item Not Found',
			message: 'The market item you requested could not be found or has been sold.',
			hint: 'Return to the market to browse available items.',
			icon: 'shopping'
		},
		default: {
			code: 'Error',
			title: 'Unknown Error',
			message: 'An unknown error occurred.',
			hint: 'Please try again or contact support.',
			icon: 'alert'
		}
	};

	// Function to get error details from URL param
	function getErrorDetails() {
		const errorParam = $page.url.searchParams.get('error');
		return errorTypes[errorParam ?? 'default'];
	}

	// Current error
	const error = getErrorDetails();

	// Function to handle return to home
	function goHome() {
		window.location.href = '/';
	}

	// Animation timeline
	let tl: gsap.core.Timeline;

	onMount(() => {
		// Initialize GSAP timeline
		tl = gsap.timeline({ defaults: { ease: 'power2.out' } });

		// Animate elements sequentially
		tl.from('.error-container', {
			opacity: 0,
			y: 30,
			duration: 0.7
		})
			.from(
				'.error-icon',
				{
					opacity: 0,
					scale: 0.5,
					duration: 0.8,
					rotation: -10
				},
				'-=0.3'
			)
			.from(
				'.error-code',
				{
					opacity: 0,
					y: -20,
					duration: 0.6
				},
				'-=0.5'
			)
			.from(
				'.error-title',
				{
					opacity: 0,
					y: -15,
					duration: 0.6
				},
				'-=0.4'
			)
			.from(
				'.error-message',
				{
					opacity: 0,
					y: 15,
					duration: 0.6
				},
				'-=0.3'
			)
			.from(
				'.error-hint',
				{
					opacity: 0,
					scale: 0.95,
					duration: 0.7
				},
				'-=0.3'
			)
			.from(
				'.action-buttons button',
				{
					opacity: 0,
					y: 15,
					stagger: 0.15,
					duration: 0.5
				},
				'-=0.4'
			);

		// Add subtle background animation
		gsap.to('.bg-pattern', {
			backgroundPosition: '100px 100px',
			duration: 20,
			repeat: -1,
			ease: 'linear'
		});

		// Add pulsing effect to error code
		gsap.to('.error-code', {
			scale: 1.05,
			duration: 1.5,
			repeat: -1,
			yoyo: true,
			ease: 'sine.inOut'
		});
	});
</script>

<div
	class="bg-pattern flex h-screen w-full items-center justify-center bg-gray-100"
	style="background-image: url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI1IiBoZWlnaHQ9IjUiPgo8cmVjdCB3aWR0aD0iNSIgaGVpZ2h0PSI1IiBmaWxsPSIjZmZmIj48L3JlY3Q+CjxyZWN0IHdpZHRoPSIxIiBoZWlnaHQ9IjEiIGZpbGw9IiNmNWY1ZjUiPjwvcmVjdD4KPC9zdmc+');"
>
	<div
		class="error-container mx-4 w-full max-w-md rounded-lg border border-gray-200 bg-white p-8 text-center shadow-lg"
	>
		<div class="mb-6">
			<div class="error-icon mb-3 text-5xl text-red-600">
				{#if error.icon === 'question-mark'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10"></circle>
						<line x1="12" y1="8" x2="12" y2="12"></line>
						<line x1="12" y1="16" x2="12.01" y2="16"></line>
					</svg>
				{:else if error.icon === 'lock'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
						<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
					</svg>
				{:else if error.icon === 'key'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path
							d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"
						></path>
					</svg>
				{:else if error.icon === 'server'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
						<rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
						<line x1="6" y1="6" x2="6.01" y2="6"></line>
						<line x1="6" y1="18" x2="6.01" y2="18"></line>
					</svg>
				{:else if error.icon === 'dungeon'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
						<polyline points="9 22 9 12 15 12 15 22"></polyline>
					</svg>
				{:else if error.icon === 'user'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
						<circle cx="12" cy="7" r="4"></circle>
					</svg>
				{:else if error.icon === 'github'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path
							d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"
						></path>
						<path d="M13 16.5v-3" stroke-width="2"></path>
					</svg>
				{:else if error.icon === 'exchange'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<polyline points="17 1 21 5 17 9"></polyline>
						<line x1="3" y1="5" x2="21" y2="5"></line>
						<polyline points="7 23 3 19 7 15"></polyline>
						<line x1="21" y1="19" x2="3" y2="19"></line>
					</svg>
				{:else if error.icon === 'shopping'}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path>
						<line x1="3" y1="6" x2="21" y2="6"></line>
						<path d="M16 10a4 4 0 0 1-8 0"></path>
					</svg>
				{:else}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="inline-block h-16 w-16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10"></circle>
						<line x1="12" y1="8" x2="12" y2="12"></line>
						<line x1="12" y1="16" x2="12.01" y2="16"></line>
					</svg>
				{/if}
			</div>
			<div class="error-code mb-2 text-4xl font-bold text-red-600">{error.code}</div>
			<div class="error-title mb-4 text-2xl font-semibold text-gray-800">{error.title}</div>
			<p class="error-message mb-6 text-gray-600">{error.message}</p>

			<div class="error-hint mb-6 rounded-md border border-gray-100 bg-gray-50 px-4 py-3">
				<p class="text-gray-700">{error.hint}</p>
			</div>

			<div class="action-buttons flex flex-col justify-center gap-3 sm:flex-row">
				<button
					on:click={goHome}
					class="rounded-md bg-red-600 px-6 py-2 font-medium text-white transition-colors duration-300 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
				>
					Return to Home
				</button>
			</div>
		</div>
	</div>
</div>

<svelte:head>
	<title>Error | RWBY Adventure</title>
</svelte:head>

<script lang="ts">
	import { onMount } from 'svelte';
	import { gsap } from 'gsap';
	import { ScrollTrigger } from 'gsap/ScrollTrigger';
	import NotificationToast from '../../components/NotificationToast.svelte';
	import { marketSSE } from '$lib/sseClient';

	// Register GSAP plugins
	gsap.registerPlugin(ScrollTrigger);

	let { children } = $props();
	let isMenuOpen = $state(false);
	let searchQuery = $state('');
	let subroute = $state('market');

	// References for animation
	let navbar: HTMLElement;
	let logoText: HTMLElement;
	let navLinks: HTMLElement;
	let searchBar: HTMLElement;
	let mainContent: HTMLElement;

	let menuItems = [
		{ name: 'Home', href: '/market', active: () => subroute === 'market' || subroute === 'home' },
		{ name: 'Listings', href: '/market/listings', active: () => subroute === 'listings' },
		{ name: 'Auctions', href: '/market/auctions', active: () => subroute === 'auctions' },
		{ name: 'About', href: '/market/about', active: () => subroute === 'about' }
	];

	function toggleMenu() {
		isMenuOpen = !isMenuOpen;
	}

	onMount(() => {
		// Initial animations when the page loads
		const timeline = gsap.timeline({ defaults: { ease: 'power2.out' } });

		marketSSE.connect();

		// Check from URL to set subroute
		const path = window.location.pathname;
		if (path.includes('/market/listings')) {
			subroute = 'listings';
		} else if (path.includes('/market/auctions')) {
			subroute = 'auctions';
		} else if (path.includes('/market/about')) {
			subroute = 'about';
		} else {
			subroute = 'market';
		}

		timeline
			.fromTo(navbar, { y: -100, opacity: 0 }, { y: 0, opacity: 1, duration: 0.8 })
			.fromTo(logoText, { x: -20, opacity: 0 }, { x: 0, opacity: 1, duration: 0.5 }, '-=0.5')
			.fromTo(
				navLinks?.children || [],
				{ y: -20, opacity: 0 },
				{ y: 0, opacity: 1, stagger: 0.1, duration: 0.3 },
				'-=0.3'
			)
			.fromTo(
				searchBar,
				{ width: '80%', opacity: 0 },
				{ width: '100%', opacity: 1, duration: 0.5 },
				'-=0.2'
			)
			.fromTo(mainContent, { opacity: 0 }, { opacity: 1, duration: 0.5 }, '-=0.3');

		// Scroll animations
		ScrollTrigger.batch('.animate-on-scroll', {
			onEnter: (elements) => {
				gsap.to(elements, {
					y: 0,
					opacity: 1,
					stagger: 0.15,
					duration: 0.8,
					ease: 'power2.out'
				});
			},
			start: 'top 85%'
		});
	});
</script>

<div class="min-h-screen bg-gray-50">
	<!-- Notification Toast Container -->
	<NotificationToast />

	<!-- Navigation -->
	<nav bind:this={navbar} class="fixed top-0 z-50 w-full bg-white shadow-md">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 justify-between">
				<!-- Logo and Brand -->
				<div class="flex items-center">
					<div class="flex flex-shrink-0 items-center">
						<span bind:this={logoText} class="text-2xl font-bold tracking-tight text-indigo-600"
							>RWBYzon</span
						>
					</div>

					<!-- Desktop Navigation Links -->
					<div class="hidden md:ml-10 md:flex md:space-x-8" bind:this={navLinks}>
						{#each menuItems as item}
							<a
								href={item.href}
								class="rounded-md px-3 py-2 text-sm font-medium transition-colors duration-200
									{item.active()
									? 'border-b-2 border-indigo-600 text-gray-900'
									: 'text-gray-700 hover:border-b-2 hover:border-indigo-400 hover:text-indigo-600'}"
								onclick={() => (subroute = item.name.toLowerCase())}
							>
								{item.name}
							</a>
						{/each}
					</div>
				</div>

				<!-- Search Bar -->
				<div class="flex flex-1 items-center justify-center px-2 lg:ml-6 lg:justify-end">
					<div bind:this={searchBar} class="flex w-full max-w-lg items-center lg:max-w-xs">
						<label for="search" class="sr-only">Search</label>
						<div class="relative w-full">
							<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
								<!-- Search Icon -->
								<svg
									class="h-5 w-5 text-gray-400"
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 20 20"
									fill="currentColor"
									aria-hidden="true"
								>
									<path
										fill-rule="evenodd"
										d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
										clip-rule="evenodd"
									/>
								</svg>
							</div>
							<input
								bind:value={searchQuery}
								id="search"
								name="search"
								class="block w-full rounded-md border border-gray-300 bg-white py-2 pl-10 pr-3 leading-5 placeholder-gray-500 transition-all duration-200 focus:border-indigo-500 focus:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
								placeholder="Search cards, listings, auctions..."
								type="search"
							/>
						</div>
					</div>
				</div>

				<!-- User Menu & Mobile Menu Button -->
				<div class="flex items-center gap-4">
					<!-- User Profile -->
					<div class="hidden md:flex">
						<button
							class="flex items-center text-sm font-medium text-gray-700 hover:text-indigo-600"
						>
							<span
								class="flex h-8 w-8 items-center justify-center rounded-full bg-gray-200 text-gray-500"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-5 w-5"
									viewBox="0 0 20 20"
									fill="currentColor"
								>
									<path
										fill-rule="evenodd"
										d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
										clip-rule="evenodd"
									/>
								</svg>
							</span>
							<span class="ml-2">Profile</span>
						</button>
					</div>

					<!-- Mobile menu button -->
					<div class="-mr-2 flex md:hidden">
						<button
							onclick={toggleMenu}
							type="button"
							class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-indigo-600 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
						>
							<span class="sr-only">Open main menu</span>
							{#if isMenuOpen}
								<!-- X icon -->
								<svg
									class="block h-6 w-6"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
									aria-hidden="true"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M6 18L18 6M6 6l12 12"
									/>
								</svg>
							{:else}
								<!-- Hamburger icon -->
								<svg
									class="block h-6 w-6"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
									aria-hidden="true"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M4 6h16M4 12h16M4 18h16"
									/>
								</svg>
							{/if}
						</button>
					</div>
				</div>
			</div>
		</div>

		<!-- Mobile Menu -->
		{#if isMenuOpen}
			<div class="animate-slideDown border-t border-gray-200 bg-white md:hidden">
				<div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
					{#each menuItems as item}
						<a
							href={item.href}
							class="block rounded-md px-3 py-2 text-base font-medium text-gray-700 transition-all duration-200
								{item.active()
								? 'border-l-4 border-indigo-500 bg-indigo-50 text-indigo-600'
								: 'hover:border-l-4 hover:border-indigo-500 hover:bg-gray-50 hover:text-indigo-600'}"
							onclick={() => (subroute = item.name.toLowerCase())}
						>
							{item.name}
						</a>
					{/each}
				</div>
				<div class="border-t border-gray-200 pb-3 pt-4">
					<div class="flex items-center px-5">
						<div class="flex-shrink-0">
							<span
								class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-200 text-gray-500"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-6 w-6"
									viewBox="0 0 20 20"
									fill="currentColor"
								>
									<path
										fill-rule="evenodd"
										d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
										clip-rule="evenodd"
									/>
								</svg>
							</span>
						</div>
						<div class="ml-3">
							<div class="text-base font-medium text-gray-800">User</div>
							<div class="text-sm font-medium text-gray-500">user@example.com</div>
						</div>
					</div>
					<div class="mt-3 space-y-1 px-2">
						<a
							href="/market/profile"
							class="block rounded-md px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-indigo-600"
							>Your Profile</a
						>
						<a
							href="/market/settings"
							class="block rounded-md px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-indigo-600"
							>Settings</a
						>
						<a
							href="/logout"
							class="block rounded-md px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-indigo-600"
							>Sign out</a
						>
					</div>
				</div>
			</div>
		{/if}
	</nav>

	<!-- Main Content -->
	<div bind:this={mainContent} class="mx-auto max-w-7xl px-4 pb-10 pt-24 sm:px-6 lg:px-8">
		{@render children()}
	</div>

	<!-- Footer -->
	<footer class="border-t border-gray-200 bg-white">
		<div class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
			<div class="grid grid-cols-1 gap-8 md:grid-cols-4">
				<div class="col-span-1 md:col-span-2">
					<span class="text-xl font-bold text-indigo-600">RWBYzon</span>
					<p class="mt-2 text-sm text-gray-500">
						The premier marketplace for RWBY card trading. Find rare cards, complete your
						collection, and connect with fellow enthusiasts.
					</p>
				</div>
				<div>
					<h3 class="text-sm font-semibold uppercase tracking-wider text-gray-600">Resources</h3>
					<ul role="list" class="mt-4 space-y-4">
						<li>
							<a href="/market/guide" class="text-base text-gray-500 hover:text-indigo-600">
								Trading Guide
							</a>
						</li>
						<li>
							<a href="/market/faq" class="text-base text-gray-500 hover:text-indigo-600"> FAQ </a>
						</li>
						<li>
							<a href="/market/rules" class="text-base text-gray-500 hover:text-indigo-600">
								Community Rules
							</a>
						</li>
					</ul>
				</div>
				<div>
					<h3 class="text-sm font-semibold uppercase tracking-wider text-gray-600">Contact</h3>
					<ul role="list" class="mt-4 space-y-4">
						<li>
							<a href="/market/support" class="text-base text-gray-500 hover:text-indigo-600">
								Support
							</a>
						</li>
						<li>
							<p class="text-base text-gray-500 hover:text-indigo-600">Discord</p>
						</li>
						<li>
							<p class="text-base text-gray-500 hover:text-indigo-600">Twitter</p>
						</li>
					</ul>
				</div>
			</div>
			<div class="mt-12 border-t border-gray-200 pt-8">
				<p class="text-center text-base text-gray-400">
					&copy; {new Date().getFullYear()} RWBYzon. All rights reserved.
				</p>
			</div>
		</div>
	</footer>
</div>

<style>
	/* Custom animations */
	:global(.animate-on-scroll) {
		opacity: 0;
		transform: translateY(20px);
	}

	/* Mobile menu animation */
	.animate-slideDown {
		animation: slideDown 0.3s ease-out forwards;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>

<!-- filepath: /home/yewolf/workspace/rwbyadv3/www/src/routes/trade/+layout.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { gsap } from 'gsap';
	import { ScrollTrigger } from 'gsap/ScrollTrigger';
	import NotificationToast from '../../components/NotificationToast.svelte';

	// Register GSAP plugins
	gsap.registerPlugin(ScrollTrigger);

	let { children } = $props();
	let isMenuOpen = $state(false);
	let searchQuery = $state('');
	let subroute = $state('trade');

	// References for animation
	let navbar: HTMLElement;
	let logoText: HTMLElement;
	let navLinks: HTMLElement;
	let searchBar: HTMLElement;
	let mainContent: HTMLElement;

	let menuItems = [
		{ name: 'Trade', href: '/trade', active: () => subroute === 'trade' },
		{ name: 'History', href: '/trade/history', active: () => subroute === 'history' },
		{ name: 'Offers', href: '/trade/offers', active: () => subroute === 'offers' }
	];

	function toggleMenu() {
		isMenuOpen = !isMenuOpen;
	}

	onMount(() => {
		// Initial animations when the page loads
		const timeline = gsap.timeline({ defaults: { ease: 'power2.out' } });

		// Check from URL to set subroute
		const path = window.location.pathname;
		if (path.includes('/trade/history')) {
			subroute = 'history';
		} else if (path.includes('/trade/offers')) {
			subroute = 'offers';
		} else {
			subroute = 'trade';
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
						<span bind:this={logoText} class="text-2xl font-bold tracking-tight text-purple-600"
							>RWBY Trading</span
						>
					</div>

					<!-- Desktop Navigation Links -->
					<div class="hidden md:ml-10 md:flex md:space-x-8" bind:this={navLinks}>
						{#each menuItems as item}
							<a
								href={item.href}
								class="rounded-md px-3 py-2 text-sm font-medium transition-colors duration-200
									{item.active()
									? 'border-b-2 border-purple-600 text-gray-900'
									: 'text-gray-700 hover:border-b-2 hover:border-purple-400 hover:text-purple-600'}"
								onclick={() => (subroute = item.name.toLowerCase())}
							>
								{item.name}
							</a>
						{/each}
					</div>
				</div>

				<!-- User Menu & Mobile Menu Button -->
				<div class="flex items-center gap-4">
					<!-- Mobile menu button -->
					<div class="-mr-2 flex md:hidden">
						<button
							onclick={toggleMenu}
							type="button"
							class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-purple-600 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-purple-500"
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
								? 'border-l-4 border-purple-500 bg-purple-50 text-purple-600'
								: 'hover:border-l-4 hover:border-purple-500 hover:bg-gray-50 hover:text-purple-600'}"
							onclick={() => (subroute = item.name.toLowerCase())}
						>
							{item.name}
						</a>
					{/each}
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
					<span class="text-xl font-bold text-purple-600">RWBY Trading</span>
					<p class="mt-2 text-sm text-gray-500">
						Trade cards with other players, expand your collection and get the cards you need.
					</p>
				</div>
				<div>
					<h3 class="text-sm font-semibold uppercase tracking-wider text-gray-600">Resources</h3>
					<ul role="list" class="mt-4 space-y-4">
						<li>
							<a href="/guide" class="text-base text-gray-500 hover:text-purple-600">
								Trading Guide
							</a>
						</li>
						<li>
							<a href="/faq" class="text-base text-gray-500 hover:text-purple-600"> FAQ </a>
						</li>
						<li>
							<a href="/rules" class="text-base text-gray-500 hover:text-purple-600">
								Trading Rules
							</a>
						</li>
					</ul>
				</div>
				<div>
					<h3 class="text-sm font-semibold uppercase tracking-wider text-gray-600">Links</h3>
					<ul role="list" class="mt-4 space-y-4">
						<li>
							<a href="/market" class="text-base text-gray-500 hover:text-purple-600">
								Marketplace
							</a>
						</li>
						<li>
							<a href="/collection" class="text-base text-gray-500 hover:text-purple-600">
								My Collection
							</a>
						</li>
						<li>
							<a href="/support" class="text-base text-gray-500 hover:text-purple-600"> Support </a>
						</li>
					</ul>
				</div>
			</div>
			<div class="mt-12 border-t border-gray-200 pt-8">
				<p class="text-center text-base text-gray-400">
					&copy; {new Date().getFullYear()} RWBY Adventure. All rights reserved.
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

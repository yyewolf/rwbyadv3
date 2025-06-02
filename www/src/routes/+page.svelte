<!-- filepath: /home/yewolf/workspace/rwbyadv3/www/src/routes/+page.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { gsap } from 'gsap';

	// Hero section animation
	onMount(() => {
		// Animate the hero section
		gsap.from('.hero-content', {
			opacity: 0,
			y: -30,
			duration: 0.8,
			ease: 'power2.out'
		});

		// Animate the feature cards with stagger
		gsap.from('.feature-card', {
			opacity: 0,
			y: 50,
			stagger: 0.2,
			duration: 0.7,
			delay: 0.3,
			ease: 'power2.out'
		});

		// Animate the call to action
		gsap.from('.cta-section', {
			opacity: 0,
			y: 30,
			duration: 0.8,
			delay: 0.9,
			ease: 'power3.out'
		});
	});
</script>

<svelte:head>
	<title>RWBY Adventures</title>
	<meta
		name="description"
		content="Welcome to RWBY Adventures - Collect cards, explore dungeons, and trade with other players!"
	/>
</svelte:head>

<main class="min-h-screen bg-gradient-to-b from-slate-900 to-slate-800 text-white">
	<!-- Hero Section -->
	<section class="hero-section w-full px-4 py-16 sm:py-24 md:py-32">
		<div class="hero-content container mx-auto text-center">
			<h1 class="mb-6 text-4xl font-extrabold tracking-tight sm:text-5xl md:text-6xl">
				<span class="bg-gradient-to-r from-red-500 to-orange-500 bg-clip-text text-transparent"
					>RWBY</span
				> Adventures
			</h1>
			<p class="mx-auto mb-8 max-w-3xl text-lg text-slate-300 sm:text-xl">
				Collect cards, explore dungeons, and trade with other players in the world of RWBY
			</p>
			<div class="flex flex-wrap justify-center gap-4">
				<a
					href="/dungeons"
					class="rounded-lg bg-red-600 px-6 py-3 font-medium text-white transition-colors duration-200 hover:bg-red-700"
				>
					Enter Dungeon
				</a>
				<a
					href="/market/listings"
					class="rounded-lg bg-slate-700 px-6 py-3 font-medium text-white transition-colors duration-200 hover:bg-slate-600"
				>
					Visit Market
				</a>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section class="features-section w-full bg-slate-800/50 px-4 py-16">
		<div class="container mx-auto">
			<h2 class="mb-12 text-center text-3xl font-bold">Game Features</h2>
			<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
				<!-- Card Collection Feature -->
				<div class="feature-card rounded-xl bg-slate-700 p-6 shadow-lg">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-500/20 text-red-500"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="h-6 w-6"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M6.429 9.75L2.25 12l4.179 2.25m0-4.5l5.571 3 5.571-3m-11.142 0L2.25 7.5 12 2.25l9.75 5.25-4.179 2.25m0 0L21.75 12l-4.179 2.25m0 0l4.179 2.25L12 21.75 2.25 16.5l4.179-2.25m11.142 0l-5.571 3-5.571-3"
							/>
						</svg>
					</div>
					<h3 class="mb-2 text-xl font-bold">Card Collection</h3>
					<p class="text-slate-300">
						Collect rare and powerful cards featuring your favorite RWBY characters with unique
						abilities and stats.
					</p>
				</div>

				<!-- Dungeons Feature -->
				<div class="feature-card rounded-xl bg-slate-700 p-6 shadow-lg">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-orange-500/20 text-orange-500"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="h-6 w-6"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M9 12.75l3 3m0 0l3-3m-3 3v-7.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
					<h3 class="mb-2 text-xl font-bold">Explore Dungeons</h3>
					<p class="text-slate-300">
						Venture into dangerous dungeons, face challenging enemies, and earn valuable loot and
						experience.
					</p>
				</div>

				<!-- Trading Feature -->
				<div class="feature-card rounded-xl bg-slate-700 p-6 shadow-lg">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-indigo-500/20 text-indigo-500"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="h-6 w-6"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5"
							/>
						</svg>
					</div>
					<h3 class="mb-2 text-xl font-bold">Player Trading</h3>
					<p class="text-slate-300">
						Trade cards with other players in the marketplace to complete your collection or sell
						rare finds.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Call to Action Section -->
	<section class="cta-section w-full px-4 py-16">
		<div class="container mx-auto text-center">
			<h2 class="mb-6 text-3xl font-bold">Ready to Start Your Adventure?</h2>
			<p class="mx-auto mb-8 max-w-2xl text-lg text-slate-300">
				Join thousands of players in the world of RWBY. Collect cards, build powerful decks, and
				challenge other players!
			</p>
			<a
				href="/dungeons"
				class="inline-block rounded-lg bg-gradient-to-r from-red-500 to-orange-500 px-8 py-4 font-bold text-white transition-colors duration-200 hover:from-red-600 hover:to-orange-600"
			>
				Start Playing Now
			</a>
		</div>
	</section>

	<!-- Footer -->
	<footer class="w-full bg-slate-900 px-4 py-8">
		<div class="container mx-auto text-center text-sm text-slate-400">
			<p>© {new Date().getFullYear()} RWBY Adventures. All rights reserved.</p>
		</div>
	</footer>
</main>

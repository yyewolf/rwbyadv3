<script lang="ts">
	import { onMount } from 'svelte';
	import gsap from 'gsap';

	// Function to handle manual window close
	function closeWindow() {
		window.close();
	}

	// Animation timeline
	let tl: gsap.core.Timeline;

	onMount(() => {
		// Initialize GSAP timeline
		tl = gsap.timeline({ defaults: { ease: 'power2.out' } });

		// Animate elements sequentially
		tl.from('.discord-icon', {
			opacity: 0,
			scale: 0.5,
			duration: 0.8,
			rotation: -10
		})
			.from(
				'.page-title',
				{
					opacity: 0,
					y: -20,
					duration: 0.6
				},
				'-=0.3'
			)
			.from(
				'.message-text',
				{
					opacity: 0,
					y: 15,
					duration: 0.6
				},
				'-=0.3'
			)
			.from(
				'.notification-box',
				{
					opacity: 0,
					scale: 0.95,
					duration: 0.7
				},
				'-=0.3'
			)
			.from(
				'.action-buttons button, .action-buttons a',
				{
					opacity: 0,
					y: 15,
					stagger: 0.15,
					duration: 0.5
				},
				'-=0.4'
			)
			.from(
				'.footer-text',
				{
					opacity: 0,
					duration: 0.6
				},
				'-=0.2'
			);
	});
</script>

<div class="flex h-screen w-full items-center justify-center bg-gray-100">
	<div
		class="mx-4 w-full max-w-md rounded-lg border border-gray-200 bg-white p-8 text-center shadow-lg"
	>
		<div class="mb-6">
			<div class="discord-icon mb-3 text-5xl text-indigo-600">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="inline-block h-16 w-16"
					viewBox="0 0 640 512"
				>
					<path
						fill="currentColor"
						d="M524.531,69.836a1.5,1.5,0,0,0-.764-.7A485.065,485.065,0,0,0,404.081,32.03a1.816,1.816,0,0,0-1.923.91,337.461,337.461,0,0,0-14.9,30.6,447.848,447.848,0,0,0-134.426,0,309.541,309.541,0,0,0-15.135-30.6,1.89,1.89,0,0,0-1.924-.91A483.689,483.689,0,0,0,116.085,69.137a1.712,1.712,0,0,0-.788.676C39.068,183.651,18.186,294.69,28.43,404.354a2.016,2.016,0,0,0,.765,1.375A487.666,487.666,0,0,0,176.02,479.918a1.9,1.9,0,0,0,2.063-.676A348.2,348.2,0,0,0,208.12,430.4a1.86,1.86,0,0,0-1.019-2.588,321.173,321.173,0,0,1-45.868-21.853,1.885,1.885,0,0,1-.185-3.126c3.082-2.309,6.166-4.711,9.109-7.137a1.819,1.819,0,0,1,1.9-.256c96.229,43.917,200.41,43.917,295.5,0a1.812,1.812,0,0,1,1.924.233c2.944,2.426,6.027,4.851,9.132,7.16a1.884,1.884,0,0,1-.162,3.126,301.407,301.407,0,0,1-45.89,21.83,1.875,1.875,0,0,0-1,2.611,391.055,391.055,0,0,0,30.014,48.815,1.864,1.864,0,0,0,2.063.7A486.048,486.048,0,0,0,610.7,405.729a1.882,1.882,0,0,0,.765-1.352C623.729,277.594,590.933,167.465,524.531,69.836ZM222.491,337.58c-28.972,0-52.844-26.587-52.844-59.239S193.056,219.1,222.491,219.1c29.665,0,53.306,26.82,52.843,59.239C275.334,310.993,251.924,337.58,222.491,337.58Zm195.38,0c-28.971,0-52.843-26.587-52.843-59.239S388.437,219.1,417.871,219.1c29.667,0,53.307,26.82,52.844,59.239C470.715,310.993,447.538,337.58,417.871,337.58Z"
					/>
				</svg>
			</div>
			<div class="page-title mb-4 text-2xl font-semibold text-gray-800">
				Success! Return to Discord
			</div>
			<p class="message-text mb-6 text-gray-600">
				Your action has been completed successfully. You can now return to Discord.
			</p>

			<div class="notification-box mb-6 rounded-md border border-indigo-100 bg-indigo-50 px-4 py-3">
				<p class="text-gray-700">
					Thank you for using RWBY Adventure! Your Discord session is waiting for you.
				</p>
			</div>

			<div class="action-buttons flex flex-col justify-center gap-3 sm:flex-row">
				<a
					href="https://discord.com/app"
					target="_blank"
					class="rounded-md bg-indigo-600 px-6 py-2 font-medium text-white transition-colors duration-300 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
				>
					Open Discord
				</a>

				<button
					on:click={closeWindow}
					class="rounded-md bg-gray-200 px-6 py-2 font-medium text-gray-800 transition-colors duration-300 hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2"
				>
					Close Window
				</button>
			</div>
		</div>

		<div class="footer-text mt-8 border-t border-gray-200 pt-4 text-sm text-gray-500">
			<p>You may close this window at any time. It's safe to do so.</p>
		</div>
	</div>
</div>

<svelte:head>
	<title>Return to Discord | RWBY Adventure</title>
</svelte:head>

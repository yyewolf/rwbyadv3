<script lang="ts">
	import { onMount } from 'svelte';

	onMount(() => {
		document.body.innerHTML = '';
	});
</script>

<svelte:head>
	<title>RWBY Adventures' Dungeons</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0" />
	<style>
		body {
			overflow-y: scroll;
			margin: 0;
			width: 100%;
			position: fixed;
		}
	</style>
	<script src="/dungeons/static/app.js"></script>
</svelte:head>

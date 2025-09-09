<script lang="ts">
	import GraphEditor from '$lib/components/GraphEditor.svelte';
	import Toolbar from '$lib/components/Toolbar.svelte';
	import StatusBar from '$lib/components/StatusBar.svelte';
	import { graphActions } from '$lib/stores/graphStore';
	import { onMount } from 'svelte';

	let editorContainer: HTMLElement;

	onMount(() => {
		// Load default graph or show empty editor
		graphActions.initializeEmpty();
	});
</script>

<svelte:head>
	<title>Axon Visual Editor</title>
	<meta name="description" content="Visual programming editor for Axon language" />
</svelte:head>

<div class="app">
	<header>
		<h1>🧠 Axon Visual Editor</h1>
		<Toolbar />
	</header>

	<main bind:this={editorContainer}>
		<GraphEditor />
	</main>

	<footer>
		<StatusBar />
	</footer>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		height: 100vh;
		background: #1a1a1a;
		color: white;
		font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}

	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		background: #2a2a2a;
		border-bottom: 1px solid #444;
		min-height: 60px;
	}

	h1 {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 600;
	}

	main {
		flex: 1;
		position: relative;
		overflow: hidden;
	}

	footer {
		background: #2a2a2a;
		border-top: 1px solid #444;
		min-height: 32px;
	}
</style>

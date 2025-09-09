<script lang="ts">
	import { editorNodes, editorEdges, editorStore } from '$lib/stores/graphStore';
	
	let nodeCount = 0;
	let edgeCount = 0;
	let viewTransform = { x: 0, y: 0, scale: 1 };
	let selectedCount = 0;

	// Subscribe to store updates
	editorNodes.subscribe(nodes => {
		nodeCount = nodes.length;
	});

	editorEdges.subscribe(edges => {
		edgeCount = edges.length;
	});

	editorStore.subscribe(state => {
		viewTransform = state.viewTransform;
		selectedCount = state.selectedNodes.size + state.selectedEdges.size;
	});
</script>

<div class="status-bar">
	<div class="status-section">
		<span class="status-item">Nodes: {nodeCount}</span>
		<span class="status-separator">•</span>
		<span class="status-item">Edges: {edgeCount}</span>
		{#if selectedCount > 0}
			<span class="status-separator">•</span>
			<span class="status-item">Selected: {selectedCount}</span>
		{/if}
	</div>

	<div class="status-section">
		<span class="status-item">
			Zoom: {Math.round(viewTransform.scale * 100)}%
		</span>
		<span class="status-separator">•</span>
		<span class="status-item">
			Position: ({Math.round(viewTransform.x)}, {Math.round(viewTransform.y)})
		</span>
	</div>

	<div class="status-section">
		<span class="status-item ready">Ready</span>
	</div>
</div>

<style>
	.status-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.25rem 1rem;
		background: #2a2a2a;
		border-top: 1px solid #444;
		font-size: 0.75rem;
		color: #ccc;
	}

	.status-section {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.status-item {
		white-space: nowrap;
	}

	.status-item.ready {
		color: #4CAF50;
	}

	.status-separator {
		opacity: 0.5;
	}
</style>
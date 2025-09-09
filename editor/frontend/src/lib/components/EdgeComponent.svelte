<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { EditorEdge } from '$lib/types';

	export let edge: EditorEdge;
	export let path: string;

	const dispatch = createEventDispatcher<{
		select: { edgeId: string };
		delete: { edgeId: string };
	}>();

	function getEdgeColor(type: string) {
		switch (type) {
			case 'data': return '#2196F3';
			case 'exec': return '#FFD700';
			default: return '#757575';
		}
	}

	function getStrokeWidth(type: string) {
		switch (type) {
			case 'data': return 2;
			case 'exec': return 3;
			default: return 1;
		}
	}

	function handleClick(event: MouseEvent) {
		dispatch('select', { edgeId: edge.id });
		event.stopPropagation();
	}

	function handleDoubleClick(event: MouseEvent) {
		dispatch('delete', { edgeId: edge.id });
		event.stopPropagation();
	}
</script>

<g class="edge" class:selected={edge.selected}>
	<!-- Main edge path -->
	<path 
		d={path}
		fill="none"
		stroke={getEdgeColor(edge.type)}
		stroke-width={getStrokeWidth(edge.type)}
		stroke-dasharray={edge.type === 'exec' ? '5,5' : 'none'}
		opacity={edge.selected ? 1 : 0.8}
		class="edge-path"
		on:click={handleClick}
		on:dblclick={handleDoubleClick}
	/>
	
	<!-- Selection indicator (invisible wider path for easier clicking) -->
	<path 
		d={path}
		fill="none"
		stroke="transparent"
		stroke-width="12"
		class="edge-hit-area"
		on:click={handleClick}
		on:dblclick={handleDoubleClick}
	/>
	
	<!-- Arrow marker for data edges -->
	{#if edge.type === 'data'}
		<defs>
			<marker 
				id="arrowhead-{edge.id}" 
				markerWidth="10" 
				markerHeight="7" 
				refX="9" 
				refY="3.5" 
				orient="auto">
				<polygon 
					points="0 0, 10 3.5, 0 7" 
					fill={getEdgeColor(edge.type)} 
				/>
			</marker>
		</defs>
		<path 
			d={path}
			fill="none"
			stroke={getEdgeColor(edge.type)}
			stroke-width={getStrokeWidth(edge.type)}
			opacity={edge.selected ? 1 : 0.8}
			marker-end="url(#arrowhead-{edge.id})"
			class="edge-path-with-arrow"
			pointer-events="none"
		/>
	{/if}
</g>

<style>
	.edge {
		pointer-events: auto;
	}

	.edge-path {
		cursor: pointer;
		transition: opacity 0.2s ease;
	}

	.edge-path:hover {
		opacity: 1;
		filter: brightness(1.2);
	}

	.edge.selected .edge-path {
		stroke-width: 4;
		filter: drop-shadow(0 0 4px currentColor);
	}

	.edge-hit-area {
		cursor: pointer;
	}

	.edge-path-with-arrow {
		pointer-events: none;
	}
</style>
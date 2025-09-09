<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { select } from 'd3-selection';
	import { zoom, zoomIdentity } from 'd3-zoom';
	import { drag } from 'd3-drag';
	import { editorNodes, editorEdges, editorStore } from '$lib/stores/graphStore';
	import NodeComponent from './NodeComponent.svelte';
	import EdgeComponent from './EdgeComponent.svelte';
	import type { EditorNode, EditorEdge } from '$lib/types';

	let svgElement: SVGSVGElement;
	let containerElement: HTMLDivElement;
	let backgroundGroup: SVGGElement;
	let edgesGroup: SVGGElement;
	let nodesGroup: SVGGElement;

	let nodes: EditorNode[] = [];
	let edges: EditorEdge[] = [];

	// Subscribe to store updates
	const unsubscribeNodes = editorNodes.subscribe(value => {
		nodes = value;
	});

	const unsubscribeEdges = editorEdges.subscribe(value => {
		edges = value;
	});

	onMount(() => {
		initializeEditor();
	});

	onDestroy(() => {
		unsubscribeNodes();
		unsubscribeEdges();
	});

	function initializeEditor() {
		if (!svgElement) return;

		const svg = select(svgElement);
		const container = select(containerElement);

		// Set up zoom behavior
		const zoomBehavior = zoom<SVGSVGElement, unknown>()
			.scaleExtent([0.1, 3])
			.on('zoom', (event) => {
				const { x, y, k } = event.transform;
				
				// Update the transform for all groups
				if (backgroundGroup) {
					select(backgroundGroup).attr('transform', `translate(${x},${y}) scale(${k})`);
				}
				if (edgesGroup) {
					select(edgesGroup).attr('transform', `translate(${x},${y}) scale(${k})`);
				}
				if (nodesGroup) {
					select(nodesGroup).attr('transform', `translate(${x},${y}) scale(${k})`);
				}

				// Update editor store
				editorStore.update(state => ({
					...state,
					viewTransform: { x, y, scale: k }
				}));
			});

		svg.call(zoomBehavior);

		// Set initial zoom
		svg.call(zoomBehavior.transform, zoomIdentity.translate(0, 0).scale(1));
	}

	function getNodePosition(node: EditorNode) {
		return {
			x: node.visual.x,
			y: node.visual.y
		};
	}

	function getEdgePath(edge: EditorEdge) {
		const sourceNode = nodes.find(n => n.id === edge.sourceNodeId);
		const targetNode = nodes.find(n => n.id === edge.targetNodeId);
		
		if (!sourceNode || !targetNode) return '';

		const sourcePos = getNodePosition(sourceNode);
		const targetPos = getNodePosition(targetNode);

		// Simple straight line for now - can be enhanced with curves
		return `M ${sourcePos.x + 150} ${sourcePos.y + 40} L ${targetPos.x} ${targetPos.y + 40}`;
	}

	function handleNodeDragStart(nodeId: string, event: MouseEvent) {
		editorStore.update(state => ({
			...state,
			isDragging: true,
			dragStartPos: { x: event.clientX, y: event.clientY }
		}));
	}

	function handleNodeDrag(nodeId: string, event: MouseEvent) {
		// This will be handled by D3 drag behavior when implemented
	}

	function handleNodeDragEnd(nodeId: string) {
		editorStore.update(state => ({
			...state,
			isDragging: false,
			dragStartPos: null
		}));
	}
</script>

<div class="editor-container" bind:this={containerElement}>
	<svg bind:this={svgElement} class="graph-svg">
		<!-- Background grid -->
		<g bind:this={backgroundGroup}>
			<defs>
				<pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
					<path d="M 20 0 L 0 0 0 20" fill="none" stroke="#333" stroke-width="1"/>
				</pattern>
			</defs>
			<rect width="100%" height="100%" fill="url(#grid)" opacity="0.3"/>
		</g>

		<!-- Edges -->
		<g bind:this={edgesGroup} class="edges-group">
			{#each edges as edge (edge.id)}
				<EdgeComponent {edge} path={getEdgePath(edge)} />
			{/each}
		</g>

		<!-- Nodes -->
		<g bind:this={nodesGroup} class="nodes-group">
			{#each nodes as node (node.id)}
				<NodeComponent 
					{node} 
					on:dragstart={(e) => handleNodeDragStart(node.id, e.detail)}
					on:drag={(e) => handleNodeDrag(node.id, e.detail)}
					on:dragend={() => handleNodeDragEnd(node.id)}
				/>
			{/each}
		</g>
	</svg>
</div>

<style>
	.editor-container {
		width: 100%;
		height: 100%;
		position: relative;
		overflow: hidden;
		background: #1a1a1a;
	}

	.graph-svg {
		width: 100%;
		height: 100%;
		cursor: grab;
	}

	.graph-svg:active {
		cursor: grabbing;
	}

	.edges-group {
		pointer-events: none;
	}

	.nodes-group {
		pointer-events: all;
	}
</style>
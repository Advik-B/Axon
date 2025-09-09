<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { EditorNode } from '$lib/types';

	export let node: EditorNode;

	const dispatch = createEventDispatcher<{
		dragstart: MouseEvent;
		drag: MouseEvent;
		dragend: MouseEvent;
		select: { nodeId: string };
	}>();

	let isDragging = false;
	let dragStartPos = { x: 0, y: 0 };
	let nodeStartPos = { x: 0, y: 0 };

	function getNodeColor(type: string) {
		switch (type) {
			case 'START': return '#4CAF50';
			case 'END': return '#F44336';
			case 'CONSTANT': return '#2196F3';
			case 'VARIABLE': return '#FF9800';
			case 'OPERATOR': return '#9C27B0';
			case 'FUNCTION': return '#00BCD4';
			case 'IF': return '#FFEB3B';
			case 'LOOP': return '#8BC34A';
			case 'STRUCT': return '#795548';
			default: return '#757575';
		}
	}

	function handleMouseDown(event: MouseEvent) {
		isDragging = true;
		dragStartPos = { x: event.clientX, y: event.clientY };
		nodeStartPos = { x: node.visual.x, y: node.visual.y };
		
		dispatch('dragstart', event);
		
		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
		
		event.preventDefault();
		event.stopPropagation();
	}

	function handleMouseMove(event: MouseEvent) {
		if (!isDragging) return;
		
		const deltaX = event.clientX - dragStartPos.x;
		const deltaY = event.clientY - dragStartPos.y;
		
		// Update node position
		node.visual.x = nodeStartPos.x + deltaX;
		node.visual.y = nodeStartPos.y + deltaY;
		
		dispatch('drag', event);
	}

	function handleMouseUp(event: MouseEvent) {
		if (isDragging) {
			isDragging = false;
			dispatch('dragend', event);
		}
		
		document.removeEventListener('mousemove', handleMouseMove);
		document.removeEventListener('mouseup', handleMouseUp);
	}

	function handleClick(event: MouseEvent) {
		dispatch('select', { nodeId: node.id });
		event.stopPropagation();
	}
</script>

<g class="node" 
   transform="translate({node.visual.x}, {node.visual.y})"
   on:mousedown={handleMouseDown}
   on:click={handleClick}>
	
	<!-- Node background -->
	<rect 
		width={node.visual.width} 
		height={node.visual.height}
		fill={getNodeColor(node.type)}
		stroke={node.visual.selected ? '#fff' : '#444'}
		stroke-width={node.visual.selected ? 2 : 1}
		rx="8"
		class="node-background"
	/>
	
	<!-- Node type indicator -->
	<rect 
		width={node.visual.width} 
		height="24"
		fill="rgba(0,0,0,0.2)"
		rx="8"
		class="node-header"
	/>
	
	<!-- Node title -->
	<text 
		x={node.visual.width / 2} 
		y="16"
		text-anchor="middle"
		fill="white"
		font-size="12"
		font-weight="bold"
		class="node-title"
	>
		{node.type}
	</text>
	
	<!-- Node label -->
	{#if node.label}
		<text 
			x={node.visual.width / 2} 
			y="42"
			text-anchor="middle"
			fill="white"
			font-size="11"
			class="node-label"
		>
			{node.label}
		</text>
	{/if}
	
	<!-- Input ports -->
	{#if node.inputs}
		{#each node.inputs as input, i}
			<g class="input-port" transform="translate(-4, {28 + i * 16})">
				<circle r="4" fill="#666" stroke="#fff" stroke-width="1"/>
				<text x="-12" y="4" text-anchor="end" fill="white" font-size="10">
					{input.name}
				</text>
			</g>
		{/each}
	{/if}
	
	<!-- Output ports -->
	{#if node.outputs}
		{#each node.outputs as output, i}
			<g class="output-port" transform="translate({node.visual.width + 4}, {28 + i * 16})">
				<circle r="4" fill="#666" stroke="#fff" stroke-width="1"/>
				<text x="12" y="4" text-anchor="start" fill="white" font-size="10">
					{output.name}
				</text>
			</g>
		{/each}
	{/if}
	
	<!-- Execution flow indicators -->
	{#if node.type !== 'CONSTANT' && node.type !== 'VARIABLE'}
		<!-- Execution input (top) -->
		<g class="exec-input" transform="translate({node.visual.width / 2}, -4)">
			<polygon points="-6,0 0,-6 6,0 0,4" fill="#FFD700" stroke="#FFA000" stroke-width="1"/>
		</g>
		
		<!-- Execution output (bottom) -->
		{#if node.type !== 'END'}
			<g class="exec-output" transform="translate({node.visual.width / 2}, {node.visual.height + 4})">
				<polygon points="-6,0 0,6 6,0 0,-4" fill="#FFD700" stroke="#FFA000" stroke-width="1"/>
			</g>
		{/if}
	{/if}
	
	<!-- Config display for constants -->
	{#if node.type === 'CONSTANT' && node.config?.value}
		<text 
			x={node.visual.width / 2} 
			y={node.visual.height - 8}
			text-anchor="middle"
			fill="white"
			font-size="10"
			font-family="monospace"
			class="node-config"
		>
			{node.config.value}
		</text>
	{/if}
</g>

<style>
	.node {
		cursor: move;
		user-select: none;
	}

	.node:hover .node-background {
		stroke: #fff;
		stroke-width: 2;
	}

	.node-title {
		pointer-events: none;
	}

	.node-label {
		pointer-events: none;
		opacity: 0.9;
	}

	.node-config {
		pointer-events: none;
		opacity: 0.8;
	}

	.input-port, .output-port {
		cursor: crosshair;
	}

	.input-port:hover circle,
	.output-port:hover circle {
		fill: #4CAF50;
		r: 5;
	}

	.exec-input, .exec-output {
		cursor: crosshair;
	}

	.exec-input:hover polygon,
	.exec-output:hover polygon {
		fill: #FFC107;
	}
</style>
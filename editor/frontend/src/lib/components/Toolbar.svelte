<script lang="ts">
	import { Save, FolderOpen, Play, Eye, Plus, Trash2, Download, Upload } from 'lucide-svelte';
	import { graphApi, graphActions } from '$lib/stores/graphStore';
	
	let showFileDialog = false;
	let showNodeDialog = false;
	let availableGraphs: string[] = [];
	let transpileResult: string = '';
	let showTranspileModal = false;

	async function handleLoadGraph() {
		const result = await graphApi.listGraphs();
		if (result.success && result.data) {
			availableGraphs = result.data;
			showFileDialog = true;
		} else {
			alert(`Failed to load graphs: ${result.message}`);
		}
	}

	async function loadSelectedGraph(filename: string) {
		const result = await graphApi.loadGraph(filename);
		if (result.success) {
			showFileDialog = false;
			alert('Graph loaded successfully!');
		} else {
			alert(`Failed to load graph: ${result.message}`);
		}
	}

	async function handleSaveGraph() {
		const filename = prompt('Enter filename (with .ax extension):');
		if (!filename) return;
		
		const result = await graphApi.saveGraph(filename);
		if (result.success) {
			alert('Graph saved successfully!');
		} else {
			alert(`Failed to save graph: ${result.message}`);
		}
	}

	async function handleValidateGraph() {
		const result = await graphApi.validateGraph();
		if (result.success) {
			alert(`Graph is valid! Nodes: ${result.data?.nodes}, Edges: ${result.data?.edges}`);
		} else {
			alert(`Graph validation failed: ${result.message}`);
		}
	}

	async function handleTranspileGraph() {
		const result = await graphApi.transpileGraph();
		if (result.success && result.data) {
			transpileResult = result.data.go_code;
			showTranspileModal = true;
		} else {
			alert(`Transpilation failed: ${result.message}`);
		}
	}

	function handleAddNode() {
		showNodeDialog = true;
	}

	function addNodeOfType(type: string) {
		const nodeId = `${type.toLowerCase()}-${Date.now()}`;
		
		graphActions.addNode({
			id: nodeId,
			type: type as any,
			label: type,
			inputs: getDefaultInputsForType(type),
			outputs: getDefaultOutputsForType(type),
			config: getDefaultConfigForType(type)
		});
		
		showNodeDialog = false;
	}

	function getDefaultInputsForType(type: string) {
		switch (type) {
			case 'OPERATOR':
				return [
					{ name: 'a', type_name: 'int' },
					{ name: 'b', type_name: 'int' }
				];
			case 'FUNCTION':
				return [{ name: 'input', type_name: 'any' }];
			default:
				return [];
		}
	}

	function getDefaultOutputsForType(type: string) {
		switch (type) {
			case 'CONSTANT':
			case 'VARIABLE':
			case 'OPERATOR':
				return [{ name: 'out', type_name: 'int' }];
			default:
				return [];
		}
	}

	function getDefaultConfigForType(type: string) {
		switch (type) {
			case 'CONSTANT':
				return { value: '0' };
			case 'OPERATOR':
				return { op: '+' };
			case 'FUNCTION':
				return { impl_reference: 'fmt.Println' };
			default:
				return {};
		}
	}

	function handleNewGraph() {
		if (confirm('Create a new graph? This will clear the current graph.')) {
			graphActions.initializeEmpty();
		}
	}
</script>

<div class="toolbar">
	<div class="toolbar-section">
		<button class="toolbar-btn" on:click={handleNewGraph} title="New Graph">
			<Plus size={18} />
		</button>
		
		<button class="toolbar-btn" on:click={handleLoadGraph} title="Open Graph">
			<FolderOpen size={18} />
		</button>
		
		<button class="toolbar-btn" on:click={handleSaveGraph} title="Save Graph">
			<Save size={18} />
		</button>
	</div>

	<div class="toolbar-divider"></div>

	<div class="toolbar-section">
		<button class="toolbar-btn" on:click={handleAddNode} title="Add Node">
			<Plus size={18} />
			<span>Add Node</span>
		</button>
		
		<button class="toolbar-btn" on:click={handleValidateGraph} title="Validate Graph">
			<Eye size={18} />
			<span>Validate</span>
		</button>
		
		<button class="toolbar-btn primary" on:click={handleTranspileGraph} title="Transpile to Go">
			<Play size={18} />
			<span>Transpile</span>
		</button>
	</div>
</div>

<!-- File Dialog Modal -->
{#if showFileDialog}
	<div class="modal-overlay" on:click={() => showFileDialog = false}>
		<div class="modal" on:click|stopPropagation>
			<h3>Load Graph</h3>
			<div class="file-list">
				{#each availableGraphs as filename}
					<div class="file-item" on:click={() => loadSelectedGraph(filename)}>
						<span>{filename}</span>
					</div>
				{/each}
			</div>
			<div class="modal-actions">
				<button on:click={() => showFileDialog = false}>Cancel</button>
			</div>
		</div>
	</div>
{/if}

<!-- Node Type Dialog Modal -->
{#if showNodeDialog}
	<div class="modal-overlay" on:click={() => showNodeDialog = false}>
		<div class="modal" on:click|stopPropagation>
			<h3>Add Node</h3>
			<div class="node-types">
				{#each ['CONSTANT', 'VARIABLE', 'OPERATOR', 'FUNCTION', 'IF', 'LOOP'] as nodeType}
					<button class="node-type-btn" on:click={() => addNodeOfType(nodeType)}>
						{nodeType}
					</button>
				{/each}
			</div>
			<div class="modal-actions">
				<button on:click={() => showNodeDialog = false}>Cancel</button>
			</div>
		</div>
	</div>
{/if}

<!-- Transpile Result Modal -->
{#if showTranspileModal}
	<div class="modal-overlay" on:click={() => showTranspileModal = false}>
		<div class="modal large" on:click|stopPropagation>
			<h3>Transpiled Go Code</h3>
			<pre class="code-block"><code>{transpileResult}</code></pre>
			<div class="modal-actions">
				<button on:click={() => navigator.clipboard.writeText(transpileResult)}>Copy</button>
				<button on:click={() => showTranspileModal = false}>Close</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.toolbar {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0;
	}

	.toolbar-section {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.toolbar-divider {
		width: 1px;
		height: 24px;
		background: #555;
		margin: 0 0.5rem;
	}

	.toolbar-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		background: #333;
		border: 1px solid #555;
		border-radius: 4px;
		color: white;
		cursor: pointer;
		transition: all 0.2s ease;
		font-size: 0.875rem;
	}

	.toolbar-btn:hover {
		background: #444;
		border-color: #666;
	}

	.toolbar-btn.primary {
		background: #0066cc;
		border-color: #0088ff;
	}

	.toolbar-btn.primary:hover {
		background: #0088ff;
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: #2a2a2a;
		border: 1px solid #555;
		border-radius: 8px;
		padding: 1.5rem;
		min-width: 300px;
		max-width: 90vw;
		max-height: 90vh;
		overflow: auto;
	}

	.modal.large {
		min-width: 600px;
		min-height: 400px;
	}

	.modal h3 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
	}

	.file-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}

	.file-item {
		padding: 0.75rem;
		background: #333;
		border: 1px solid #555;
		border-radius: 4px;
		cursor: pointer;
		transition: background 0.2s ease;
	}

	.file-item:hover {
		background: #444;
	}

	.node-types {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: 0.5rem;
		margin-bottom: 1rem;
	}

	.node-type-btn {
		padding: 1rem;
		background: #333;
		border: 1px solid #555;
		border-radius: 4px;
		color: white;
		cursor: pointer;
		transition: all 0.2s ease;
		font-size: 0.875rem;
	}

	.node-type-btn:hover {
		background: #444;
		border-color: #666;
	}

	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 0.5rem;
	}

	.modal-actions button {
		padding: 0.5rem 1rem;
		background: #333;
		border: 1px solid #555;
		border-radius: 4px;
		color: white;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.modal-actions button:hover {
		background: #444;
	}

	.code-block {
		background: #1a1a1a;
		border: 1px solid #555;
		border-radius: 4px;
		padding: 1rem;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		font-size: 0.875rem;
		line-height: 1.5;
		white-space: pre-wrap;
		overflow: auto;
		max-height: 400px;
		margin-bottom: 1rem;
	}
</style>
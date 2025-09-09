import { writable, derived } from 'svelte/store';
import type { AxonGraph, EditorNode, EditorEdge, GraphEditorState, ApiResponse } from '$lib/types';

// Graph data store
export const graphStore = writable<AxonGraph | null>(null);

// Editor state store
export const editorStore = writable<GraphEditorState>({
	selectedNodes: new Set(),
	selectedEdges: new Set(),
	viewTransform: { x: 0, y: 0, scale: 1 },
	isDragging: false,
	dragStartPos: null
});

// Derived stores
export const editorNodes = derived(graphStore, ($graph) => {
	if (!$graph) return [];
	
	return $graph.nodes.map((node, index) => ({
		...node,
		visual: {
			x: node.visual?.x ?? 100 + (index * 200),
			y: node.visual?.y ?? 100 + (index * 100),
			width: node.visual?.width ?? 150,
			height: node.visual?.height ?? 80,
			selected: false,
			dragging: false
		}
	})) as EditorNode[];
});

export const editorEdges = derived(graphStore, ($graph) => {
	if (!$graph) return [];
	
	const edges: EditorEdge[] = [];
	
	// Add data edges
	$graph.data_edges.forEach((edge, index) => {
		edges.push({
			id: `data-${index}`,
			type: 'data',
			sourceNodeId: edge.from_node_id,
			targetNodeId: edge.to_node_id,
			sourcePort: edge.from_port,
			targetPort: edge.to_port,
			selected: false
		});
	});
	
	// Add execution edges
	$graph.exec_edges.forEach((edge, index) => {
		edges.push({
			id: `exec-${index}`,
			type: 'exec',
			sourceNodeId: edge.from_node_id,
			targetNodeId: edge.to_node_id,
			selected: false
		});
	});
	
	return edges;
});

// API integration
const API_BASE = '/api';

async function apiRequest<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
	try {
		const response = await fetch(`${API_BASE}${endpoint}`, {
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			},
			...options
		});
		
		return await response.json();
	} catch (error) {
		return {
			success: false,
			message: `Request failed: ${error}`
		};
	}
}

export const graphApi = {
	// Load a graph from file
	async loadGraph(filename: string) {
		const result = await apiRequest<AxonGraph>(`/graphs/load?file=${filename}`);
		if (result.success && result.data) {
			graphStore.set(result.data);
		}
		return result;
	},
	
	// Save current graph
	async saveGraph(filename: string) {
		const graph = await new Promise<AxonGraph | null>(resolve => {
			const unsubscribe = graphStore.subscribe(value => {
				unsubscribe();
				resolve(value);
			});
		});
		
		if (!graph) {
			return { success: false, message: 'No graph to save' };
		}
		
		return await apiRequest(`/graphs/save?file=${filename}`, {
			method: 'POST',
			body: JSON.stringify(graph)
		});
	},
	
	// Validate current graph
	async validateGraph() {
		const graph = await new Promise<AxonGraph | null>(resolve => {
			const unsubscribe = graphStore.subscribe(value => {
				unsubscribe();
				resolve(value);
			});
		});
		
		if (!graph) {
			return { success: false, message: 'No graph to validate' };
		}
		
		return await apiRequest('/graphs/validate', {
			method: 'POST',
			body: JSON.stringify(graph)
		});
	},
	
	// Transpile current graph to Go
	async transpileGraph() {
		const graph = await new Promise<AxonGraph | null>(resolve => {
			const unsubscribe = graphStore.subscribe(value => {
				unsubscribe();
				resolve(value);
			});
		});
		
		if (!graph) {
			return { success: false, message: 'No graph to transpile' };
		}
		
		return await apiRequest<{ go_code: string }>('/graphs/transpile', {
			method: 'POST',
			body: JSON.stringify(graph)
		});
	},
	
	// List available graphs
	async listGraphs() {
		return await apiRequest<string[]>('/graphs');
	}
};

export const graphActions = {
	initializeEmpty() {
		const emptyGraph: AxonGraph = {
			id: 'new-graph',
			name: 'New Graph',
			imports: [],
			nodes: [
				{
					id: 'start',
					type: 'START',
					label: 'Start',
					visual: { x: 100, y: 200 }
				},
				{
					id: 'end',
					type: 'END',
					label: 'End',
					visual: { x: 500, y: 200 }
				}
			],
			data_edges: [],
			exec_edges: []
		};
		graphStore.set(emptyGraph);
	},
	
	addNode(node: Omit<EditorNode, 'visual'> & { visual?: Partial<EditorNode['visual']> }) {
		graphStore.update(graph => {
			if (!graph) return graph;
			
			const newNode = {
				...node,
				visual: {
					x: 300,
					y: 300,
					width: 150,
					height: 80,
					...node.visual
				}
			};
			
			return {
				...graph,
				nodes: [...graph.nodes, newNode]
			};
		});
	},
	
	updateNode(nodeId: string, updates: Partial<EditorNode>) {
		graphStore.update(graph => {
			if (!graph) return graph;
			
			return {
				...graph,
				nodes: graph.nodes.map(node =>
					node.id === nodeId ? { ...node, ...updates } : node
				)
			};
		});
	},
	
	deleteNode(nodeId: string) {
		graphStore.update(graph => {
			if (!graph) return graph;
			
			return {
				...graph,
				nodes: graph.nodes.filter(node => node.id !== nodeId),
				data_edges: graph.data_edges.filter(
					edge => edge.from_node_id !== nodeId && edge.to_node_id !== nodeId
				),
				exec_edges: graph.exec_edges.filter(
					edge => edge.from_node_id !== nodeId && edge.to_node_id !== nodeId
				)
			};
		});
	}
};
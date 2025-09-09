export interface AxonGraph {
	id: string;
	name: string;
	imports: string[];
	nodes: AxonNode[];
	data_edges: DataEdge[];
	exec_edges: ExecEdge[];
}

export interface AxonNode {
	id: string;
	type: NodeType;
	label?: string;
	inputs?: Port[];
	outputs?: Port[];
	config?: Record<string, any>;
	impl_reference?: string;
	visual?: NodeVisual;
}

export interface Port {
	name: string;
	type_name: string;
}

export interface DataEdge {
	from_node_id: string;
	from_port: string;
	to_node_id: string;
	to_port: string;
}

export interface ExecEdge {
	from_node_id: string;
	to_node_id: string;
}

export interface NodeVisual {
	x: number;
	y: number;
	width?: number;
	height?: number;
}

export type NodeType = 
	| 'START' 
	| 'END' 
	| 'CONSTANT' 
	| 'VARIABLE' 
	| 'OPERATOR' 
	| 'FUNCTION' 
	| 'IF' 
	| 'LOOP' 
	| 'STRUCT';

export interface GraphEditorState {
	selectedNodes: Set<string>;
	selectedEdges: Set<string>;
	viewTransform: {
		x: number;
		y: number;
		scale: number;
	};
	isDragging: boolean;
	dragStartPos: { x: number; y: number } | null;
}

export interface ApiResponse<T = any> {
	success: boolean;
	message: string;
	data?: T;
}

// Editor specific types
export interface EditorNode extends AxonNode {
	visual: NodeVisual & {
		selected?: boolean;
		dragging?: boolean;
	};
}

export interface EditorEdge {
	id: string;
	type: 'data' | 'exec';
	sourceNodeId: string;
	targetNodeId: string;
	sourcePort?: string;
	targetPort?: string;
	selected?: boolean;
}
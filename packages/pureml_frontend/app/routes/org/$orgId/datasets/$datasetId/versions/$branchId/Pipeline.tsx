import { useCallback } from "react";
import ReactFlow, {
  addEdge,
  ConnectionLineType,
  useNodesState,
  useEdgesState,
  Controls,
  MarkerType,
} from "reactflow";
import dagre from "dagre";
import DetailedNode from "./customNode";

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 190;
const nodeHeight = 110;

const getLayoutedElements = (nodes: any, edges: any, direction = "TB") => {
  const isHorizontal = direction === "LR";
  dagreGraph.setGraph({ rankdir: direction });
  nodes.forEach((node: any) => {
    dagreGraph.setNode(node.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  });
  edges.forEach((edge: any) => {
    edge.type = "smoothstep";
    edge.markerEnd = {
      type: MarkerType.ArrowClosed,
    };
    dagreGraph.setEdge(edge.source, edge.target);
  });
  dagre.layout(dagreGraph);
  nodes.forEach((node: any) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.targetPosition = isHorizontal ? "left" : "top";
    node.sourcePosition = isHorizontal ? "right" : "bottom";

    // We are shifting the dagre node position (anchor=center center) to the top left
    // so it matches the React Flow node anchor point (top left).
    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };
    node.data = {
      ...node.data,
      id: node.id,
      title: node.text,
      type: node.nodeType,
      desc: node.desc,
      time: node.time,
    };
    node.type = "detailedNode";
    node.style =
      node.nodeType === "load_data"
        ? { background: "#CDBBFF", borderRadius: "8px" }
        : node.nodeType === "transformer"
        ? { background: "#D6E0FF", borderRadius: "8px" }
        : node.nodeType === "dataset"
        ? { background: "#FFDFA6", borderRadius: "8px" }
        : { background: "#f1f5f9", borderRadius: "8px" };

    return node;
  });
  return { nodes, edges };
};

const nodeTypes = { detailedNode: DetailedNode };

function Pipeline({ pnode, pedge }: any) {
  const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
    pnode,
    pedge
  );
  const [nodes, setNodes, onNodesChange] = useNodesState(layoutedNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(layoutedEdges);
  const onConnect = useCallback(
    (params) =>
      setEdges((eds) =>
        addEdge(
          {
            ...params,
            type: ConnectionLineType.SmoothStep,
            markerEnd: {
              type: MarkerType.ArrowClosed,
            },
            animated: true,
          },
          eds
        )
      ),
    [setEdges]
  );
  return (
    <div className="h-4/5">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        connectionLineType={ConnectionLineType.SmoothStep}
        // fitView
        nodesDraggable={false}
      >
        <Controls position="top-right" />
        {/* <MiniMap zoomable pannable /> */}
      </ReactFlow>
    </div>
  );
}

export default Pipeline;

// ############################ error boundary ###########################

export function ErrorBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
    </div>
  );
}

export function CatchBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
    </div>
  );
}

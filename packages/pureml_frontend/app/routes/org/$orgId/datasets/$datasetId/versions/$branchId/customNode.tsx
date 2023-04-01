import { Table } from "lucide-react";
import { Handle, Position } from "reactflow";

function DetailedNode({ data, isConnectable }: any) {
  return (
    <div className="border border-slate-200 rounded-lg w-64">
      <Handle
        type="target"
        position={Position.Top}
        isConnectable={isConnectable}
      />
      <div className="flex gap-x-2 bg-slate-100 px-4 py-1.5 rounded-t-lg">
        <Table className="w-4" />
        <div className="truncate w-3/4">{data.title || "Node Title"}</div>
      </div>
      <div className="bg-white px-4 py-1.5 truncate">
        {data.desc || "Description"}
      </div>
      <div className="flex justify-between px-4 py-1.5">
        {data.type || "Node Type"}
        <div>{data.time || "00:00:00"}</div>
      </div>
      <Handle
        type="source"
        position={Position.Bottom}
        id="b"
        isConnectable={isConnectable}
      />
    </div>
  );
}

export default DetailedNode;

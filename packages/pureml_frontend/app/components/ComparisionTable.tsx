import { ChevronDown, ChevronUp } from "lucide-react";
import React, { useEffect, useState } from "react";

interface Props {
  metric: string;
  ver1: string;
  ver2: string;
  dataVer1: {
    [key: string]: string;
  };
  dataVer2: {
    [key: string]: string;
  };
}

function ComparisionTable({ metric, ver1, ver2, dataVer1, dataVer2 }: Props) {
  const [open, setOpen] = useState(true);
  const [commonMetrics, setCommonMetrics] = useState<string[]>([]);

  // ##### fetching & displaying latest version data #####
  useEffect(() => {
    if (!dataVer1) return;
    // if (Object.keys(dataVer1).length === 0)
    setCommonMetrics(Object.keys(dataVer1));
  }, [dataVer1]);
  useEffect(() => {
    if (!dataVer2) return;
    if (dataVer1) {
      setCommonMetrics(Object.keys(dataVer1));
    } else {
      setCommonMetrics(Object.keys(dataVer2));
    }

    Object.keys(dataVer2).forEach((key) => {
      if (!commonMetrics.includes(key)) {
        setCommonMetrics((prev) => [...prev, key]);
      }
    });
  }, [dataVer2]);

  return (
    <section>
      <div
        className="flex items-center justify-between w-full border-b-slate-300 border-b pb-4"
        onClick={() => setOpen(!open)}
      >
        <h1 className="text-slate-800 font-medium text-sm">{metric}</h1>
        {open ? (
          <ChevronUp className="text-slate-400" />
        ) : (
          <ChevronDown className="text-slate-400" />
        )}
      </div>

      {open && (
        <div className="py-6">
          {commonMetrics.length !== 0 ? (
            <>
              <table className="max-w-[1000px] w-full">
                {commonMetrics.length !== 0 && (
                  <>
                    <thead>
                      <tr>
                        <th className="text-slate-600 font-medium text-left border p-4">
                          {" "}
                        </th>
                        <th className="text-slate-600 font-medium text-left border p-4 w-1/5">
                          {ver1}
                        </th>
                        {ver2 !== "" ? (
                          <th className="text-slate-600 font-medium text-left border p-4 w-1/5">
                            {ver2}
                          </th>
                        ) : null}
                      </tr>
                    </thead>
                    {commonMetrics.map((metric, i) => (
                      <tr key={i}>
                        <th className="text-slate-600 font-medium text-left border p-4">
                          {metric}
                        </th>
                        <td className="text-slate-600 font-medium text-left border p-4 w-1/5 truncate">
                          {Object.keys(dataVer1).length > 0
                            ? dataVer1[metric]
                              ? dataVer1[metric].slice(0, 5)
                              : "-"
                            : "No-data"}
                        </td>
                        {ver2 !== "" && (
                          <td className="text-slate-600 font-medium text-left border p-4 w-1/5 truncate">
                            {dataVer2
                              ? dataVer2[metric]
                                ? dataVer2[metric].slice(0, 5)
                                : "-"
                              : "No-data"}
                          </td>
                        )}
                      </tr>
                    ))}
                  </>
                )}
              </table>
            </>
          ) : (
            <div>No {metric} available</div>
          )}
        </div>
      )}
    </section>
  );
}

export default ComparisionTable;

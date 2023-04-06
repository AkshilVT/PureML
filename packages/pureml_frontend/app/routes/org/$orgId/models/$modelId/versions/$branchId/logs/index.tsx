import { useLoaderData } from "@remix-run/react";
import { Suspense, useContext, useEffect, useState } from "react";
import { getSession } from "~/session";
import Loader from "~/components/ui/Loading";
import ComparisionTable from "~/components/ComparisionTable";
import {
  fetchModelBranch,
  fetchModelVersions,
} from "~/routes/api/models.server";
import VersionContext from "../../versionContext";

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const allBranch = await fetchModelBranch(
    session.get("orgId"),
    params.modelId,
    session.get("accessToken")
  );
  // console.log(params.modelId, params.branchId);

  const versions = await fetchModelVersions(
    session.get("orgId"),
    params.modelId,
    params.branchId,
    session.get("accessToken")
  );
  // console.log(versions);

  return {
    versions: versions,
    branches: allBranch,
    params: params,
  };
}

function isJson(item: string | object) {
  let value = typeof item !== "string" ? JSON.stringify(item) : item;
  try {
    value = JSON.parse(value);
  } catch (e) {
    return false;
  }
  return typeof value === "object" && value !== null;
}

export default function ModelMetrics() {
  const versionContext = useContext(VersionContext);
  const ver1 = versionContext.ver1;
  const ver2 = versionContext.ver2;
  const data = useLoaderData();
  // const [ver1, setVer1] = useState('');
  // const [ver2, setVer2] = useState('');
  const [dataVersion, setDataVersion] = useState({});
  const [ver1Logs, setVer1Logs] = useState<{ [key: string]: string }>({});
  const [ver2Logs, setVer2Logs] = useState<{ [key: string]: string }>({});
  const [commonMetrics, setCommonMetrics] = useState<string[]>([]);
  const [versionData, setVersionData] = useState(data.versions);

  useEffect(() => {
    setVersionData(data.versions);
  }, [data.versions]);

  // console.log(versionData);
  // ##### checking version data #####
  useEffect(() => {
    if (!versionData) return;
    if (!versionData[0]) return;
    if (!versionContext) return;

    // setVer1(versionContext.ver1);
    // setVer2('');
  }, [versionData]);

  // ##### fetching & displaying latest version data #####
  useEffect(() => {
    if (!versionData) return;

    const tempDict = {};
    versionData.forEach((version: { version: any }) => {
      // @ts-ignore
      tempDict[version.version] = version;
    });
    setDataVersion(tempDict);
    const tt = dataVersion[ver1];
    // console.log('tt=', tt);
    if (tt) {
      if (tt.logs === null) {
        setVer1Logs({});
        setCommonMetrics([]);
        return;
      } else {
        const tempDictv1 = {};
        tt.logs.forEach((log: { key: string; data: any }) => {
          if (isJson(log.data)) {
            tempDictv1[log.key] = JSON.parse(log.data);
            if (!commonMetrics.includes(log.key)) {
              setCommonMetrics((prev) => [...prev, log.key]);
            }
          }
        });
        setVer1Logs(tempDictv1);
        // console.log(tt.logs);
      }
    }
  }, [versionData]);
  // ##### comparing versions #####
  useEffect(() => {
    if (!versionData) return;
    const t1 = dataVersion[ver1];
    // console.log('tt=', tt);
    if (t1) {
      if (t1.logs === null) {
        setVer1Logs({});
        setCommonMetrics([]);
        return;
      } else {
        const tempDictv1 = {};
        t1.logs.forEach((log: { key: string; data: any }) => {
          if (isJson(log.data)) {
            tempDictv1[log.key] = JSON.parse(log.data);
            if (!commonMetrics.includes(log.key)) {
              setCommonMetrics((prev) => [...prev, log.key]);
            }
          }
        });
        setVer1Logs(tempDictv1);
      }
    }
    if (ver2 === "") {
      setVer2Logs({});
      // console.log('ver2 is empty');

      return;
    }

    const tt = dataVersion[ver2];
    // console.log('tt2=', tt);
    if (tt) {
      if (tt.logs === null) {
        setVer2Logs({});
        return;
      } else {
        const tempDictv2 = {};
        // console.log('tt.logs=', tt.logs);

        tt.logs.forEach((log: { data: any }) => {
          try {
            tempDictv2[log.key] = JSON.parse(log.data);
            if (!commonMetrics.includes(log.key)) {
              setCommonMetrics((prev) => [...prev, log.key]);
            }
          } catch {
            console.log("Invalid log.key=", log.key);
          }
        });
        // console.log(tempDictv2);

        setVer2Logs(tempDictv2);
      }
    }
  }, [ver1, ver2, versionData]);

  return (
    <Suspense fallback={<Loader />}>
      <div className="px-12 pt-2 pb-8 h-[70vh] overflow-auto">
        {commonMetrics.length !== 0 && versionData !== null ? (
          <>
            {commonMetrics.map((key) => {
              return (
                <ComparisionTable
                  key={key}
                  metric={key}
                  ver1={ver1}
                  ver2={ver2}
                  dataVer1={
                    ver1Logs[key] as unknown as {
                      [key: string]: string;
                    }
                  }
                  dataVer2={
                    ver2Logs[key] as unknown as {
                      [key: string]: string;
                    }
                  }
                />
              );
            })}
          </>
        ) : null}
      </div>
    </Suspense>
  );
}

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

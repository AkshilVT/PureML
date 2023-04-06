import type { MetaFunction } from "@remix-run/node";
import { useActionData, useLoaderData } from "@remix-run/react";
import { GitFork, GitPullRequest, User } from "lucide-react";
import { Suspense, useEffect, useState } from "react";
import Breadcrumbs from "~/components/Breadcrumbs";
import ReviewTabbar from "~/components/ReviewTabbar";
import Tabbar from "~/components/Tabbar";
import Loader from "~/components/ui/Loading";
import { edgesSchema, nodesSchema, versionDataSchema } from "~/lib.schema";
import Pipeline from "../../../versions/$branchId/Pipeline";
import AvatarIcon from "~/components/ui/Avatar";
import { getSession } from "~/session";
import {
  fetchDatasetVersions,
  fetchOneDatasetVersion,
} from "~/routes/api/datasets.server";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Review Dataset Commit | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  // fetching review model version from it's branch
  const reviewData = await fetchOneDatasetVersion(
    session.get("orgId"),
    params.datasetId,
    session.get("fromBranch"),
    session.get("version"),
    session.get("accessToken")
  );
  const versions = await fetchDatasetVersions(
    session.get("orgId"),
    params.datasetId,
    "dev",
    session.get("accessToken")
  );
  // "dev" branch should be replaced by to_branch of submitted review
  return {
    reviewData,
    reviewVersion: session.get("version"),
    reviewBranch: session.get("fromBranch"),
    toBranch: session.get("toBranch"),
    versions: versions,
    params: params,
  };
}

export default function Review() {
  const data = useLoaderData();
  const adata = useActionData();
  const [versionData, setVersionData] = useState(data.versions);
  const [ver1, setVer1] = useState("");
  const [node, setNode] = useState(null);
  const [edge, setEdge] = useState(null);
  const reviewData = data.reviewData;
  const reviewVersion = data.reviewVersion;
  const reviewBranch = data.reviewBranch;
  const toBranch = data.toBranch;

  useEffect(() => {
    if (!adata) return;
    adata.versions !== undefined
      ? setVersionData(adata.versions)
      : setVersionData(null);
  }, [adata]);

  // ##### checking version data #####
  useEffect(() => {
    if (!versionData) return;
    if (!versionData[0]) return;
    if (!versionData[0].lineage.lineage) return;
    setVer1(versionData.at(0).version);
  }, [versionData]);

  useEffect(() => {
    if (!versionData) return;
    let lineageData1 = versionData.at(-Number(ver1.slice(1))).lineage.lineage;
    lineageData1 = lineageData1.replace(/'/g, "");
    const validJson = versionDataSchema.safeParse(lineageData1);
    const validNodes = nodesSchema.safeParse(JSON.parse(lineageData1).nodes);
    const validEdges = edgesSchema.safeParse(JSON.parse(lineageData1).edges);
    if (validJson.success && validNodes.success && validEdges.success) {
      const n = JSON.parse(lineageData1).nodes;
      n.forEach((e: any) => {
        e.data = { label: e.text };
      });
      setNode(null);
      setTimeout(() => {
        setNode(n);
      }, 10);
      const ed = JSON.parse(lineageData1).edges;
      ed.forEach((e: any) => {
        e.source = e.from;
        e.target = e.to;
      });
      setEdge(null);
      setTimeout(() => {
        setEdge(ed);
      }, 10);
    } else {
      setNode(null);
      setEdge(null);
    }
  }, [ver1, versionData]);

  return (
    <Suspense fallback={<Loader />}>
      <div className="flex justify-center sticky top-0 bg-slate-50 w-full border-b border-slate-200">
        <div className="flex justify-between px-12 2xl:pr-0 w-full max-w-screen-2xl">
          <Breadcrumbs />
          <Tabbar intent="primaryDatasetTab" tab="review" fullWidth={false} />
        </div>
      </div>
      <div className="flex justify-center w-full">
        <div className="bg-slate-50 flex flex-col h-screen overflow-hidden w-full 2xl:max-w-screen-2xl">
          <div className="flex justify-between h-full">
            <div className="w-4/5">
              <ReviewTabbar
                intent="datasetReviewLineageTab"
                tab="datalineage"
              />
              <div className="px-12 pt-8 pb-8 h-[100vh] overflow-auto">
                <section>
                  {!versionData && (
                    <div className="text-slate-600 pt-6">
                      No data lineage available
                    </div>
                  )}
                  {versionData && (
                    <div className="flex pt-6 gap-x-8">
                      <div className="w-full h-screen max-h-full">
                        {node && <Pipeline pnode={node} pedge={edge} />}
                      </div>
                    </div>
                  )}
                </section>
              </div>
            </div>
            {/* ##### versions list right sidebar ##### */}
            <aside className="bg-slate-50 border-l-2 border-slate-100 h-full w-1/4 max-w-[400px] py-8 px-12 z-10">
              <div className="flex justify-end">
                <select
                  name="branch"
                  className="text-slate-500 border-slate-400 bg-transparent rounded"
                  disabled
                >
                  <option value="value" selected>
                    dev
                  </option>
                </select>
              </div>
              {/* incoming branch for review will be displayed here */}
              {versionData ? (
                <ul className="space-y-2 mt-8 overflow-auto">
                  <li className="flex items-center justify-between py-2">
                    <div className="flex items-center">
                      <input
                        className="checkbox checkbox-primary checkbox-sm"
                        name={"version2"}
                        type="checkbox"
                        defaultChecked
                        disabled
                      />
                      <div className="flex items-center justify-center pl-4 text-slate-600">
                        <AvatarIcon>R</AvatarIcon>
                        <div>
                          <p className="px-4 font-medium">
                            Rejected review {reviewVersion}
                          </p>
                          <p className="px-4">from branch {reviewBranch}</p>
                        </div>
                      </div>
                    </div>
                  </li>
                </ul>
              ) : (
                <div className="py-8 mt-8">No version created yet</div>
              )}
              <div className="h-3/5 mt-8">
                <div className="text-base text-slate-800 font-medium pb-4">
                  Details
                </div>
                <div className="h-4/5 overflow-auto">
                  <div className="text-slate-400">
                    <div className="flex justify-between items-center py-1">
                      <span className="w-1/7 flex items-center">
                        <GitFork className="w-4" />
                      </span>
                      <span className="w-1/2 pl-2">Commit version</span>
                      <span className="w-1/2 pl-2 text-slate-600 font-medium">
                        {reviewVersion || "Version"}
                      </span>
                    </div>
                    <div className="flex justify-between items-center py-1">
                      <span className="w-1/7 flex items-center">
                        <GitPullRequest className="w-4" />
                      </span>
                      <span className="w-1/2 pl-2">Commit branch</span>
                      <span className="w-1/2 pl-2 text-slate-600 font-medium">
                        {reviewBranch || "Branch"}
                      </span>
                    </div>
                    {/* <div className="flex justify-between items-center py-1">
                      <span className="w-1/7 flex items-center">
                        <User className="w-4" />
                      </span>
                      <span className="w-1/2 pl-2">Pushed by</span>
                      <span className="w-1/2 pl-2 text-slate-600 font-medium">
                        Name
                      </span>
                    </div> */}
                  </div>
                </div>
              </div>
            </aside>
          </div>
        </div>
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

import type { MetaFunction } from "@remix-run/node";
import {
  Form,
  useActionData,
  useLoaderData,
  useNavigate,
} from "@remix-run/react";
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
  updateDatasetReview,
} from "~/routes/api/datasets.server";
import Button from "~/components/ui/Button";
import { toast } from "react-toastify";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Review Dataset Commit | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  // fetching review dataset version from it's branch
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
    session.get("toBranch"),
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

export async function action({ params, request }: any) {
  const formData = await request.formData();
  let option = Object.fromEntries(formData);
  const session = await getSession(request.headers.get("Cookie"));
  if (option._action) {
    const updatedReview = await updateDatasetReview(
      session.get("orgId"),
      params.datasetId,
      params.reviewId,
      option.isAccepted,
      session.get("accessToken")
    );
    return {
      _action: option._action,
      updatedReview: updatedReview,
    };
  } else {
    return null;
  }
}

export default function Review() {
  const data = useLoaderData();
  const reviewData = data.reviewData;
  const reviewVersion = data.reviewVersion;
  const reviewBranch = data.reviewBranch;
  const toBranch = data.toBranch;

  const adata = useActionData();
  const navigate = useNavigate();
  const [versionData, setVersionData] = useState(data.versions);
  const [ver1, setVer1] = useState("");
  const [ver2, setVer2] = useState("");
  const [dataVersion, setDataVersion] = useState({});

  const [node, setNode] = useState(null);
  const [edge, setEdge] = useState(null);
  const [node2, setNode2] = useState(null);
  const [edge2, setEdge2] = useState(null);

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
    setVer1(versionData.at(0).version);
    setVer2("");

    const tempDict = {};
    versionData.forEach((version: { version: any }) => {
      // @ts-ignore
      tempDict[version.version] = version;
    });
    setDataVersion(tempDict);
  }, [versionData]);

  // ##### fetching & comparing latest version data #####
  useEffect(() => {
    if (!versionData) return;
    if (Object.keys(dataVersion).length === 0) return;
    if (ver2 === "") {
      setNode2(null);
      setEdge2(null);
      return;
    }

    let lineageData2 = dataVersion[ver2].lineage.lineage;
    lineageData2 = lineageData2.replace(/'/g, "");
    const validJson = versionDataSchema.safeParse(lineageData2);
    const validNodes = nodesSchema.safeParse(JSON.parse(lineageData2).nodes);
    const validEdges = edgesSchema.safeParse(JSON.parse(lineageData2).edges);
    if (validJson.success && validNodes.success && validEdges.success) {
      const n = JSON.parse(lineageData2).nodes;
      n.forEach((e: any) => {
        e.data = { label: e.text };
      });
      setNode2(null);
      setTimeout(() => {
        setNode2(n);
      }, 10);
      const ed = JSON.parse(lineageData2).edges;
      ed.forEach((e: any) => {
        e.source = e.from;
        e.target = e.to;
      });
      setEdge2(null);
      setTimeout(() => {
        setEdge2(ed);
      }, 10);
    } else {
      setNode2(null);
      setEdge2(null);
    }
  }, [dataVersion, ver2, versionData]);

  useEffect(() => {
    if (!versionData) return;
    // if (Object.keys(dataVersion).length === 0) return;
    let lineageData1 = reviewData[0].lineage.lineage;
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
  }, [reviewData, ver1, versionData]);

  // ##### review action functionality #####
  useEffect(() => {
    if (adata === null || adata === undefined) {
      return;
    }
    if (adata._action === "rejected") {
      if (adata.updatedReview.message === "Dataset review updated")
        toast.success("Dataset review version rejected!");
      else toast.error("Something went wrong!");
    } else if (adata._action === "accepted") {
      if (adata.updatedReview.message === "Dataset review updated")
        toast.success("Dataset review version accepted!");
      else toast.error("Something went wrong!");
    } else {
      setVersionData(adata.versions);
    }
    navigate(
      `/org/${data.params.orgId}/datasets/${data.params.datasetId}/versions/datalineage`
    );
  }, [adata]);

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
                      {!node2 && (
                        <div className="w-full h-screen max-h-full">
                          {node && <Pipeline pnode={node} pedge={edge} />}
                        </div>
                      )}
                      {node2 && (
                        <>
                          <div className="w-1/2 h-screen max-h-full">
                            <div className="text-sm text-slate-600 font-medium pb-6">
                              Submitted commit {reviewVersion}
                            </div>
                            {node && <Pipeline pnode={node} pedge={edge} />}
                          </div>
                          <div className="w-1/2 h-screen max-h-full">
                            <div className="text-sm text-slate-600 font-medium pb-6">
                              {ver2}
                            </div>
                            {node2 && <Pipeline pnode={node2} pedge={edge2} />}
                          </div>
                        </>
                      )}
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
                    {toBranch}
                  </option>
                </select>
              </div>
              {/* incoming branch for review will be displayed here */}
              {versionData ? (
                <ul className="h-2/5 space-y-2 mt-8 overflow-auto">
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
                        <AvatarIcon>S</AvatarIcon>
                        <div>
                          <p className="px-4 font-medium">
                            Submitted review {reviewVersion}
                          </p>
                          <p className="px-4">from branch {reviewBranch}</p>
                        </div>
                      </div>
                    </div>
                  </li>
                  {versionData.map((version: any) => (
                    <li
                      className="flex items-center justify-between py-2"
                      key={version.version}
                    >
                      <div className="flex items-center">
                        <input
                          className="checkbox checkbox-primary checkbox-sm"
                          name={"version2"}
                          value={version.version}
                          type="checkbox"
                          checked={version.version === ver2}
                          onChange={(e) => {
                            setVer2(version.version);
                          }}
                        />
                        <div className="flex items-center justify-center pl-4 text-slate-600">
                          <AvatarIcon intent="avatar">
                            {version.created_by.name.charAt(0).toUpperCase()}
                          </AvatarIcon>
                          <div>
                            <p className="px-4 font-medium">
                              {version.version}
                            </p>
                            <p className="px-4">
                              Created by {version.created_by.name}
                            </p>
                          </div>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              ) : (
                <div className="h-2/5 py-8 mt-8">No version created yet</div>
              )}
              <div className="h-2/5 mt-8">
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
                        {`${datasetDetails[0].is_public ? "Yes" : "No"}`}
                      </span>
                    </div> */}
                  </div>
                  <div className="flex gap-x-4 mt-6">
                    <Form method="post">
                      <input name="_action" value="rejected" type="hidden" />
                      <input
                        name="fromBranch"
                        value={data.from_branch}
                        type="hidden"
                      />
                      <input
                        name="fromBranchVersion"
                        value={data.from_branch_version}
                        type="hidden"
                      />
                      <input name="isAccepted" value="false" type="hidden" />
                      <Button intent="secondary" fullWidth={false}>
                        Reject
                      </Button>
                    </Form>
                    <Form method="post">
                      <input name="_action" value="accepted" type="hidden" />
                      <input
                        name="fromBranch"
                        value={data.from_branch}
                        type="hidden"
                      />
                      <input
                        name="fromBranchVersion"
                        value={data.from_branch_version}
                        type="hidden"
                      />
                      <input name="isAccepted" value="true" type="hidden" />
                      <Button intent="primary" fullWidth={false}>
                        Approve
                      </Button>
                    </Form>
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

import Tabbar from "~/components/Tabbar";
import {
  Form,
  useActionData,
  useLoaderData,
  useNavigate,
  useSubmit,
} from "@remix-run/react";
import { getSession } from "~/session";
import {
  fetchDatasetBranch,
  fetchDatasetVersions,
  submitDatasetReview,
} from "~/routes/api/datasets.server";
import Pipeline from "../Pipeline";
import { Suspense, useContext, useEffect, useState } from "react";
import type { MetaFunction } from "@remix-run/node";
import { Check, MoreVertical } from "lucide-react";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import * as SelectPrimitive from "@radix-ui/react-select";
import AvatarIcon from "~/components/ui/Avatar";
import { edgesSchema, nodesSchema, versionDataSchema } from "~/lib.schema";
import Select from "~/components/ui/Select";
import Loader from "~/components/ui/Loading";
import Breadcrumbs from "~/components/Breadcrumbs";
import { toast } from "react-toastify";
import VersionContext from "../../versionContext";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Dataset Versions | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const allBranch = await fetchDatasetBranch(
    session.get("orgId"),
    params.datasetId,
    session.get("accessToken")
  );
  const versions = await fetchDatasetVersions(
    session.get("orgId"),
    params.datasetId,
    params.branchId,
    session.get("accessToken")
  );
  return {
    versions: versions,
    branches: allBranch,
    params: params,
  };
}

export async function action({ params, request }: any) {
  const formData = await request.formData();
  let option = Object.fromEntries(formData);
  const session = await getSession(request.headers.get("Cookie"));
  if (option._action === "changeBranch") {
    let branch = formData.get("branch");
    const versions = await fetchDatasetVersions(
      session.get("orgId"),
      params.datasetId,
      branch,
      session.get("accessToken")
    );
    return {
      versions: versions,
    };
  } else if (option._action === "submitReview") {
    const submitReview = await submitDatasetReview(
      session.get("orgId"),
      params.datasetId,
      option.fromBranch,
      option.version,
      session.get("accessToken")
    );
    return { _action: "submitReview", submitReview: submitReview };
  } else {
    return null;
  }
}

export default function DatasetVersions() {
  const versionContext = useContext(VersionContext);
  const ver1 = versionContext.ver1;
  const ver2 = versionContext.ver2;
  const data = useLoaderData();
  const adata = useActionData();
  const submit = useSubmit();
  const navigate = useNavigate();
  const [versionData, setVersionData] = useState(data.versions);
  const [submitReviewVersion, setSubmitReviewVersion] = useState("");
  const branchData = data.branches;
  // const [ver1, setVer1] = useState('');
  // const [ver2, setVer2] = useState('');
  const [dataVersion, setDataVersion] = useState({});
  const [node, setNode] = useState(null);
  const [edge, setEdge] = useState(null);
  const [node2, setNode2] = useState(null);
  const [edge2, setEdge2] = useState(null);
  const [branch, setBranch] = useState("main");

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
    // setVer1(versionData.at(0).version);
    // setVer2('');

    const tempDict = {};
    versionData.forEach((version: { version: any }) => {
      // @ts-ignore
      tempDict[version.version] = version;
    });
    setDataVersion(tempDict);
  }, [versionData]);

  useEffect(() => {
    if (!versionData) return;
    if (Object.keys(dataVersion).length === 0) return;
    if (dataVersion[ver1].lineage) {
      if (dataVersion[ver1].lineage.lineage) {
        let lineageData1 = dataVersion[ver1].lineage.lineage;
        lineageData1 = lineageData1.replace(/'/g, "");
        const validJson = versionDataSchema.safeParse(lineageData1);
        const validNodes = nodesSchema.safeParse(
          JSON.parse(lineageData1).nodes
        );
        const validEdges = edgesSchema.safeParse(
          JSON.parse(lineageData1).edges
        );
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
      } else {
        setNode(null);
        setEdge(null);
      }
    }
  }, [ver1, dataVersion]);

  // ##### fetching & comparing latest version data #####
  useEffect(() => {
    if (!versionData) return;
    if (Object.keys(dataVersion).length === 0) return;

    if (ver2 === "") {
      setNode2(null);
      setEdge2(null);
      return;
    }
    // console.log(ver2, dataVersion[ver2]);

    if (dataVersion[ver2].lineage) {
      if (dataVersion[ver2].lineage.lineage) {
        let lineageData2 = dataVersion[ver2].lineage.lineage;
        // console.log(typeof lineageData2);
        lineageData2 = lineageData2.replace(/'/g, "");
        const validJson = versionDataSchema.safeParse(lineageData2);
        const validNodes = nodesSchema.safeParse(
          JSON.parse(lineageData2).nodes
        );
        const validEdges = edgesSchema.safeParse(
          JSON.parse(lineageData2).edges
        );
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
          // console.log(n, ed);
        } else {
          setNode2(null);
          setEdge2(null);
        }
      } else {
        setNode2(null);
        setEdge2(null);
      }
    }
  }, [ver2, dataVersion]);

  // ##### submit review functionality #####

  useEffect(() => {
    if (adata === null || adata === undefined) {
      return;
    }
    if (adata._action === "submitReview") {
      toast.success(adata.submitReview.message);
      navigate(
        `/org/${data.params.orgId}/datasets/${data.params.datasetId}/review`
      );
    } else {
      setVersionData(adata.versions);
    }
  }, [adata]);

  function branchChange(event: any) {
    setBranch(event.target.value);
    submit(event.currentTarget, { replace: true });
  }
  function submitReview(event: any) {
    submit(event.currentTarget, { replace: true });
  }

  return (
    <Suspense fallback={<Loader />}>
      <div className="px-12 pt-2 pb-8 h-[100vh] overflow-auto">
        <section>
          {!versionData && (
            <div className="text-slate-600">No data lineage available</div>
          )}
          {versionData && (
            <div className="flex pt-6 gap-x-8">
              {!ver2 && !node2 && (
                <div className="w-full h-[80vh] max-h-full">
                  {node ? (
                    <Pipeline pnode={node} pedge={edge} />
                  ) : (
                    "No datalineage available"
                  )}
                </div>
              )}
              {ver2 && (
                <>
                  <div className="w-1/2 h-[80vh] max-h-full">
                    <div className="text-sm text-slate-600 font-medium pb-6">
                      {ver1}
                    </div>
                    {node ? (
                      <Pipeline pnode={node} pedge={edge} />
                    ) : (
                      "No datalineage available"
                    )}
                  </div>
                  <div className="w-1/2 h-[80vh] max-h-full">
                    <div className="text-sm text-slate-600 font-medium pb-6">
                      {ver2}
                    </div>
                    {node2 ? (
                      <Pipeline pnode={node2} pedge={edge2} />
                    ) : (
                      "No datalineage available"
                    )}
                  </div>
                </>
              )}
            </div>
          )}
        </section>
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

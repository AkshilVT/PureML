import { Form, useActionData, useLoaderData } from "@remix-run/react";
import { Suspense, useState } from "react";
import { ChevronDown, ChevronUp } from "lucide-react";
import Tabbar from "~/components/Tabbar";
import { useSubmit, useTransition } from "@remix-run/react";

import {
  fetchModelMetrics,
  fetchModelVersions,
} from "~/routes/api/models.server";
import { getSession } from "~/session";
import Dropdown from "~/components/ui/Dropdown";
import Loader from "~/components/ui/Loading";
import Breadcrumbs from "~/components/Breadcrumbs";

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const versions = await fetchModelVersions(
    session.get("orgId"),
    params.modelId,
    "dev",
    session.get("accessToken")
  );
  const metrics = await fetchModelMetrics(
    session.get("orgId"),
    params.modelId,
    versions.at(-1),
    session.get("accessToken")
  );
  // return [metrics, projectId, url.href];
  return {
    metrics1: metrics,
    versions: versions,
  };
}

export async function action({ params, request }: any) {
  const formData = await request.formData();
  let version1 = formData.get("version1");
  let version2 = formData.get("version2");
  let v = formData.get("v");
  const session = await getSession(request.headers.get("Cookie"));
  if (v === "true") {
    version1 = version2;
    version2 = null;
  }
  const metrics1 =
    version1 !== null
      ? await fetchModelMetrics(
          session.get("orgId"),
          params.modelId,
          version1,
          session.get("accessToken")
        )
      : null;
  const metrics2 =
    version2 !== null
      ? await fetchModelMetrics(
          session.get("orgId"),
          params.modelId,
          version2,
          session.get("accessToken")
        )
      : null;
  return {
    metrics1: metrics1,
    metrics2: metrics2,
    version1: version1,
    version2: version2,
  };
}

export default function ModelGraphs() {
  // const data = useLoaderData();
  // const adata = useActionData();
  // const submit = useSubmit();
  // const transition = useTransition();
  // const [graphs, setGraphs] = useState(true);
  // const versions = data.versions;
  // let metricsData = JSON.parse(data.metrics1).metrics[1];
  // let v1 = versions.at(-1);
  // let v2 = "";
  // let metricsData2: any;
  // if (adata) {
  //   metricsData = adata.metrics1 !== null ? adata.metrics1 : data.metrics1;
  //   v1 = adata.version1 !== null ? adata.version1 : versions.at(-1);
  //   v2 = adata.version2 !== null ? adata.version2 : "";
  //   metricsData2 = adata.version2 !== null ? adata.metrics2 : [];
  // }
  // function versionChange(event: any) {
  //   submit(event.currentTarget, { replace: true });
  // }
  return (
    <Suspense fallback={<Loader />}>
      <div className="flex justify-center sticky top-0 bg-slate-50 w-full border-b border-slate-200">
        <div className="flex justify-between px-12 2xl:pr-0 w-full max-w-screen-2xl">
          <Breadcrumbs />
          <Tabbar intent="primaryModelTab" tab="versions" fullWidth={false} />
        </div>
      </div>
      <div className="flex justify-center w-full">
        <div className="bg-slate-50 flex flex-col h-screen overflow-hidden w-full 2xl:max-w-screen-2xl">
          <div className="flex h-full">
            <div className="w-full">
              <Tabbar intent="modelTab" tab="graphs" />
              <div className="px-12 py-8">
                <section className="w-full">
                  <div className="px-2">Coming Soon...</div>
                </section>
              </div>
            </div>
            {/* <aside className="bg-slate-50 border-l-2 border-slate-100 w-1/3 max-w-[400px] max-h-[700px] px-4 py-6">
          <Dropdown fullWidth={false} intent="branch">
            dev
          </Dropdown>
          <div className="py-4">Status: {JSON.stringify(transition.state)}</div>
          <ul className="space-y-2">
            {versions.map((version: any) => (
              <li className="flex items-center space-x-2" key={version}>
                <Form method="post" onChange={versionChange}>
                  <input hidden name="version1" value={v1} />
                  <input hidden name="v" value={version === v1} />
                  <input
                    // name={version.version === v1 ? 'version1' : 'version2'}
                    name={"version2"}
                    value={version}
                    type="checkbox"
                    checked={version === v1 || version === v2}
                  />
                </Form>
                <p>{version}</p>
              </li>
            ))}
          </ul>
        </aside> */}
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

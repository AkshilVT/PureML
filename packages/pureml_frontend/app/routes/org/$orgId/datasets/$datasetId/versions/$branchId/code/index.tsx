import Tabbar from "~/components/Tabbar";
import { CopyBlock, paraisoLight } from "react-code-blocks";
import {
  fetchDatasetBranch,
  fetchDatasetVersions,
} from "~/routes/api/datasets.server";
import { Suspense, useEffect, useState } from "react";
import {
  Form,
  useActionData,
  useLoaderData,
  useSubmit,
} from "@remix-run/react";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import * as SelectPrimitive from "@radix-ui/react-select";
import { Check, MoreVertical } from "lucide-react";
import AvatarIcon from "~/components/ui/Avatar";
import { getSession } from "~/session";
import Select from "~/components/ui/Select";
import Loader from "~/components/ui/Loading";
import Breadcrumbs from "~/components/Breadcrumbs";

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
    allBranch[0].value,
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
  let branch = formData.get("branch");
  const session = await getSession(request.headers.get("Cookie"));
  const versions = await fetchDatasetVersions(
    session.get("orgId"),
    params.datasetId,
    branch,
    session.get("accessToken")
  );
  return {
    versions: versions,
  };
}

export default function DatasetCode() {
  const data = useLoaderData();
  const adata = useActionData();
  const submit = useSubmit();
  function branchChange(event: any) {
    setBranch(event.target.value);
    submit(event.currentTarget, { replace: true });
  }
  const [versionData, setVersionData] = useState(data.versions);
  useEffect(() => {
    if (!adata) return;
    adata.versions !== undefined
      ? setVersionData(adata.versions)
      : setVersionData(null);
  }, [adata]);

  const branchData = data.branches;
  const [ver1, setVer1] = useState("");
  const [ver2, setVer2] = useState("");
  const [node, setNode] = useState(null);
  const [edge, setEdge] = useState(null);
  const [node2, setNode2] = useState(null);
  const [edge2, setEdge2] = useState(null);
  const [branch, setBranch] = useState("main");

  // ##### checking version data #####
  useEffect(() => {
    if (!versionData) return;
    if (!versionData[0]) return;
    setVer1(versionData.at(0).version);
    setVer2("");
  }, [versionData]);

  // ##### fetching & comparing latest version data #####
  useEffect(() => {
    if (!versionData) return;
    if (ver2 === "") {
      setNode2(null);
      setEdge2(null);
      return;
    }
  }, [ver2, versionData]);

  useEffect(() => {
    if (!versionData) return;
  }, [ver1, versionData]);

  return (
    <Suspense fallback={<Loader />}>
      <div className="flex justify-center sticky top-0 bg-slate-50 w-full border-b border-slate-200">
        <div className="flex justify-between px-12 2xl:pr-0 w-full max-w-screen-2xl">
          <Breadcrumbs />
          <Tabbar intent="primaryDatasetTab" tab="versions" fullWidth={false} />
        </div>
      </div>
      <div className="flex justify-center w-full">
        <div className="bg-slate-50 flex flex-col h-screen overflow-hidden w-full 2xl:max-w-screen-2xl">
          <div className="flex justify-between h-full">
            <div className="w-4/5">
              <Tabbar intent="datasetTab" tab="code" />
              <div className="px-12 pt-2 pb-8 h-[100vh] w-[72%] overflow-auto">
                {/* <CopyBlock
            text={`@dataset('telecom churn')             
def build_dataset():
    ...
    return df
df = build_dataset()
`}
            language="python"
            theme={paraisoLight}
            wrapLines
          /> */}
                <section>
                  {!versionData && (
                    <div className="text-slate-600">No code available</div>
                  )}
                  {versionData && (
                    <div className="flex pt-6 gap-x-8">
                      <div className="w-full h-screen max-h-[600px]">
                        <CopyBlock
                          text={`@dataset('telecom churn')             
def build_dataset():
    ...
    return df
df = build_dataset()
`}
                          language="python"
                          theme={paraisoLight}
                          wrapLines
                        />
                      </div>
                    </div>
                  )}
                </section>
              </div>
            </div>
            <aside className="bg-slate-50 border-l-2 border-slate-100 h-full w-1/4 max-w-[400px] py-8 px-12 z-10">
              <Form
                method="post"
                onChange={branchChange}
                className="flex justify-end"
              >
                <Select name="branch" title={branch}>
                  {branchData.map((branch: any, index: number) => (
                    <SelectPrimitive.Item
                      key={`${branch}-${index}`}
                      value={branch.value}
                      className="flex items-center justify-between px-4 py-2 text-base text-slate-600 font-medium cursor-pointer hover:bg-slate-100 hover:border-none focus:outline-none"
                    >
                      <SelectPrimitive.ItemText className="text-slate-600 text-base font-medium">
                        {branch.label}
                      </SelectPrimitive.ItemText>
                      <SelectPrimitive.ItemIndicator>
                        <Check className="text-slate-400 w-4" />
                      </SelectPrimitive.ItemIndicator>
                    </SelectPrimitive.Item>
                  ))}
                </Select>
              </Form>
              {versionData ? (
                <ul className="h-3/4 space-y-2 mt-8 overflow-auto">
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
                          checked={
                            version.version === ver1 || version.version === ver2
                          }
                          onChange={(e) => {
                            if (e.target.checked) {
                              setVer2(version.version);
                            } else if (
                              ver1 === version.version &&
                              ver2 === ""
                            ) {
                              new Error(
                                "You can't uncheck the present version"
                              );
                            } else if (ver1 === version.version) {
                              setVer1(ver2);
                              setVer2("");
                            } else {
                              setVer2("");
                            }
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
                      <DropdownMenu.Root>
                        <DropdownMenu.Trigger className="focus:outline-none">
                          <MoreVertical className="text-slate-400" />
                        </DropdownMenu.Trigger>
                        <DropdownMenu.Content sideOffset={7} align="end">
                          <DropdownMenu.Item
                            className="bg-white p-2"
                            // onClick={(e) => {
                            //   submitForReview(e);
                            // }}
                          >
                            {" "}
                            <Form method="post">
                              <input
                                hidden
                                value={version.version}
                                name="version"
                              />
                              <button
                                name="_action"
                                type="submit"
                                value="submitForReview"
                              >
                                Submit for Review
                              </button>
                            </Form>
                          </DropdownMenu.Item>
                        </DropdownMenu.Content>
                      </DropdownMenu.Root>
                    </li>
                  ))}
                </ul>
              ) : (
                <div className="py-8">No version created yet</div>
              )}
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

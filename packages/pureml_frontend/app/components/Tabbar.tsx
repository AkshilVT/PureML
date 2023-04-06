import clsx from "clsx";
import { tv, type VariantProps } from "tailwind-variants";
import { Link, useMatches } from "@remix-run/react";

const tabStyles = tv({
  base: "pt-6 text-zinc-400 flex sticky z-10 pl-12",
  variants: {
    intent: {
      primaryModelTab: "bg-slate-50",
      primaryDatasetTab: "bg-slate-50",
      primarySettingTab: "top-[3.3rem]",
      modelTab: "pb-6",
      datasetTab: "pb-6",
    },
    fullWidth: {
      true: "w-full",
      false: "w-fit",
    },
  },
  defaultVariants: {
    intent: "primaryModelTab",
    fullWidth: true,
  },
});

interface Props extends VariantProps<typeof tabStyles> {
  tab: string;
}

function primaryLinkCss(currentPage: boolean) {
  return clsx(
    currentPage ? "text-slate-600 font-medium" : "text-slate-500",
    "flex justify-center items-center hover:text-slate-600"
  );
}

export default function TabBar(props: Props) {
  const matches = useMatches();
  const path = matches[2].pathname;
  const pathname = decodeURI(path.slice(1));
  const orgId = pathname.split("/")[1];
  const modelId = pathname.split("/")[3];
  const datasetId = pathname.split("/")[3];
  const primaryModelTabs = [
    {
      id: "modelCard",
      name: "Model Card",
      hyperlink: `/org/${orgId}/models/${modelId}`,
    },
    {
      id: "versions",
      name: "Versions",
      hyperlink: `/org/${orgId}/models/${modelId}/versions/main/logs`,
    },
    {
      id: "review",
      name: "Review",
      hyperlink: `/org/${orgId}/models/${modelId}/review`,
    },
  ];
  const primaryDatasetTabs = [
    {
      id: "datasetCard",
      name: "Dataset Card",
      hyperlink: `/org/${orgId}/datasets/${datasetId}`,
    },
    {
      id: "versions",
      name: "Versions",
      hyperlink: `/org/${orgId}/datasets/${datasetId}/versions/main/datalineage`,
    },
    {
      id: "review",
      name: "Review",
      hyperlink: `/org/${orgId}/datasets/${datasetId}/review`,
    },
  ];
  const primarySettingTabs = [
    {
      id: "profile",
      name: "Profile",
      hyperlink: `/settings`,
    },
    {
      id: "organization",
      name: "Organization",
      hyperlink: `/settings/organization`,
    },
    {
      id: "members",
      name: "Members",
      hyperlink: `/settings/members`,
    },
    // {
    //   id: "billing",
    //   name: "Billing",
    //   hyperlink: `/settings/billing`,
    // },
  ];
  const modelTabs = [
    {
      id: "metrics",
      name: "User Logs",
      hyperlink: `/org/${orgId}/models/${modelId}/versions/main/logs`,
    },
    // {
    //   id: "graphs",
    //   name: "Graphs",
    //   hyperlink: `/org/${orgId}/models/${modelId}/versions/logs`,
    // },
  ];
  const datasetTabs = [
    {
      id: "datalineage",
      name: "Data Lineage",
      hyperlink: `/org/${orgId}/datasets/${datasetId}/versions/main/datalineage`,
    },
    // {
    //   id: "code",
    //   name: "Code",
    //   hyperlink: `/org/${orgId}/datasets/${datasetId}/versions/code`,
    // },
  ];
  return (
    <div className={tabStyles(props)}>
      <div className="flex">
        {props.intent === "primaryModelTab" ||
        props.intent === "primaryDatasetTab" ||
        props.intent === "primarySettingTab" ? (
          <>
            {props.intent === "primaryModelTab" ? (
              <>
                {Object.keys(primaryModelTabs).map((key: string) => (
                  <div key={key} className="pl-6">
                    <div
                      className={`${
                        props.tab === primaryModelTabs[key as never].id
                          ? "text-brand-200 border-b-2 border-brand-200 font-medium"
                          : "text-slate-500"
                      } pb-4 hover:text-slate-850`}
                    >
                      <Link to={primaryModelTabs[key as any].hyperlink}>
                        <span>{primaryModelTabs[key as any].name}</span>
                      </Link>
                    </div>
                  </div>
                ))}
              </>
            ) : (
              <>
                {props.intent === "primaryDatasetTab" ? (
                  <>
                    {Object.keys(primaryDatasetTabs).map((key: string) => (
                      <div key={key} className="pl-6">
                        <div
                          className={`${
                            props.tab === primaryDatasetTabs[key as never].id
                              ? "text-brand-200 border-b-2 border-brand-200 font-medium"
                              : "text-slate-500"
                          } pb-4 hover:text-slate-850`}
                        >
                          <Link to={primaryDatasetTabs[key as any].hyperlink}>
                            <span>{primaryDatasetTabs[key as any].name}</span>
                          </Link>
                        </div>
                      </div>
                    ))}
                  </>
                ) : (
                  <>
                    {Object.keys(primarySettingTabs).map((key: string) => (
                      <div key={key} className="pr-6">
                        <div
                          className={`${
                            props.tab === primarySettingTabs[key as never].id
                              ? "text-brand-200 border-b-2 border-brand-200 font-medium"
                              : "text-slate-500"
                          } pb-4 hover:text-slate-850`}
                        >
                          <Link to={primarySettingTabs[key as any].hyperlink}>
                            <span>{primarySettingTabs[key as any].name}</span>
                          </Link>
                        </div>
                      </div>
                    ))}
                  </>
                )}
              </>
            )}
          </>
        ) : (
          <>
            {props.intent === "modelTab" ? (
              <div className="flex">
                {Object.keys(modelTabs).map((key: string) => (
                  <div key={key} className="pr-2">
                    <div
                      className={`${
                        props.tab === modelTabs[key as never].id
                          ? "bg-slate-200 rounded text-slate-600"
                          : ""
                      } px-4 py-2`}
                    >
                      <Link
                        to={modelTabs[key as any].hyperlink}
                        className={`${primaryLinkCss(
                          props.tab === modelTabs[key as any].id
                        )}`}
                      >
                        <span>{modelTabs[key as any].name}</span>
                      </Link>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex">
                {Object.keys(datasetTabs).map((key: string) => (
                  <div key={key} className="pr-2">
                    <div
                      className={`${
                        props.tab === datasetTabs[key as never].id
                          ? "bg-slate-200 rounded text-slate-600"
                          : ""
                      } px-4 py-2`}
                    >
                      <Link
                        to={datasetTabs[key as any].hyperlink}
                        className={`${primaryLinkCss(
                          props.tab === datasetTabs[key as any].id
                        )}`}
                      >
                        <span>{datasetTabs[key as any].name}</span>
                      </Link>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

import type { MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { useState } from "react";
import CTASection from "~/components/landingPage/CTASection";
import Footer from "~/components/landingPage/Footer";
import Navbar from "~/components/landingPage/Navbar";
import { fetchMLTools } from "./api/auth.server";
import Tabs from "react-simply-tabs";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "ML Tools | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader() {
  const toolsList = await fetchMLTools();
  return toolsList;
}

export default function MLTools() {
  const data = useLoaderData();
  const [activeTabIndex, setActiveTabIndex] = useState(0);
  return (
    <div className="bg-slate-50 landingpg-font flex flex-col justify-center">
      <div className="mltoolsbg bg-cover">
        <Navbar />
        <div className="flex flex-col justify-center md:items-center md:text-center gap-y-4 px-4 md:px-0 h-80 md:h-96 lg:h-[32rem]">
          <h1 className="font-bold text-4xl md:text-5xl lg:text-6xl">
            MLOps Tools
          </h1>
          <h1 className="text-slate-600 text-lg md:text-2xl lg:text-3xl md:w-4/5 lg:w-3/5 xl:w-2/5 2xl:w-1/5">
            A selection of the finest MLOps tools to help you construct the
            ideal machine learning stack.
          </h1>
        </div>
      </div>
      <div className="bg-slate-50 flex justify-center">
        <div className="w-full md:max-w-screen-xl px-4 md:px-8">
          <div className="py-8 md:py-16">
            <Tabs
              activeTabIndex={activeTabIndex}
              onRequestChange={setActiveTabIndex}
              controls={[
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab1"
                >
                  ALL TOOLS
                </button>,
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab2"
                >
                  DATA
                </button>,
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab3"
                >
                  MODEL
                </button>,
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab4"
                >
                  MODEL SERVING
                </button>,
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab5"
                >
                  ORCHESTRATION
                </button>,
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab6"
                >
                  MONITORING
                </button>,
                <button
                  type="button"
                  className="text-base text-slate-400"
                  key="tab7"
                >
                  AUTO ML
                </button>,
              ]}
              controlsWrapperProps={{
                style: {
                  paddingBottom: "10px",
                  paddingTop: "10px",
                  borderTop: "1px solid #CBD5E1",
                  borderBottom: "1px solid #CBD5E1",
                  display: "flex",
                  columnGap: "2rem",
                  overflow: "auto",
                },
              }}
              activeControlProps={{
                className: "active",
                style: { color: "#1E293B" },
              }}
            >
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <a
                        href={tool.properties.Link.url}
                        key={key}
                        className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                      >
                        <div className="flex gap-x-4">
                          <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                            <img
                              src={tool.properties.Logo.files[0].external.url}
                              alt="Logo"
                              className="w-8 h-8"
                            />
                          </div>
                          <div>
                            <div className="text-base text-slate-400">
                              {tool.properties.Category.select.name}
                            </div>
                            <div className="text-2xl font-medium">
                              {tool.properties.Name.title[0].plain_text}
                            </div>
                          </div>
                        </div>
                        <div className="flex flex-col gap-y-2">
                          <span className="text-xl text-justify">
                            {
                              tool.properties.Description.rich_text[0]
                                .plain_text
                            }
                          </span>
                        </div>
                      </a>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <>
                        {tool.properties.Category.select.name === "Data" && (
                          <a
                            href={tool.properties.Link.url}
                            key={key}
                            className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                          >
                            <div className="flex gap-x-4">
                              <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                                <img
                                  src={
                                    tool.properties.Logo.files[0].external.url
                                  }
                                  alt="Logo"
                                  className="w-8 h-8"
                                />
                              </div>
                              <div>
                                <div className="text-base text-slate-400">
                                  {tool.properties.Category.select.name}
                                </div>
                                <div className="text-2xl font-medium">
                                  {tool.properties.Name.title[0].plain_text}
                                </div>
                              </div>
                            </div>
                            <div className="flex flex-col gap-y-2">
                              <span className="text-xl text-justify">
                                {
                                  tool.properties.Description.rich_text[0]
                                    .plain_text
                                }
                              </span>
                            </div>
                          </a>
                        )}
                      </>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <>
                        {tool.properties.Category.select.name === "Model" && (
                          <a
                            href={tool.properties.Link.url}
                            key={key}
                            className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                          >
                            <div className="flex gap-x-4">
                              <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                                <img
                                  src={
                                    tool.properties.Logo.files[0].external.url
                                  }
                                  alt="Logo"
                                  className="w-8 h-8"
                                />
                              </div>
                              <div>
                                <div className="text-base text-slate-400">
                                  {tool.properties.Category.select.name}
                                </div>
                                <div className="text-2xl font-medium">
                                  {tool.properties.Name.title[0].plain_text}
                                </div>
                              </div>
                            </div>
                            <div className="flex flex-col gap-y-2">
                              <span className="text-xl text-justify">
                                {
                                  tool.properties.Description.rich_text[0]
                                    .plain_text
                                }
                              </span>
                            </div>
                          </a>
                        )}
                      </>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <>
                        {tool.properties.Category.select.name ===
                          "Model Serving" && (
                          <a
                            href={tool.properties.Link.url}
                            key={key}
                            className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                          >
                            <div className="flex gap-x-4">
                              <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                                <img
                                  src={
                                    tool.properties.Logo.files[0].external.url
                                  }
                                  alt="Logo"
                                  className="w-8 h-8"
                                />
                              </div>
                              <div>
                                <div className="text-base text-slate-400">
                                  {tool.properties.Category.select.name}
                                </div>
                                <div className="text-2xl font-medium">
                                  {tool.properties.Name.title[0].plain_text}
                                </div>
                              </div>
                            </div>
                            <div className="flex flex-col gap-y-2">
                              <span className="text-xl text-justify">
                                {
                                  tool.properties.Description.rich_text[0]
                                    .plain_text
                                }
                              </span>
                            </div>
                          </a>
                        )}
                      </>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <>
                        {tool.properties.Category.select.name ===
                          "Orchestration" && (
                          <a
                            href={tool.properties.Link.url}
                            key={key}
                            className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                          >
                            <div className="flex gap-x-4">
                              <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                                <img
                                  src={
                                    tool.properties.Logo.files[0].external.url
                                  }
                                  alt="Logo"
                                  className="w-8 h-8"
                                />
                              </div>
                              <div>
                                <div className="text-base text-slate-400">
                                  {tool.properties.Category.select.name}
                                </div>
                                <div className="text-2xl font-medium">
                                  {tool.properties.Name.title[0].plain_text}
                                </div>
                              </div>
                            </div>
                            <div className="flex flex-col gap-y-2">
                              <span className="text-xl text-justify">
                                {
                                  tool.properties.Description.rich_text[0]
                                    .plain_text
                                }
                              </span>
                            </div>
                          </a>
                        )}
                      </>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <>
                        {tool.properties.Category.select.name ===
                          "Monitoring" && (
                          <a
                            href={tool.properties.Link.url}
                            key={key}
                            className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                          >
                            <div className="flex gap-x-4">
                              <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                                <img
                                  src={
                                    tool.properties.Logo.files[0].external.url
                                  }
                                  alt="Logo"
                                  className="w-8 h-8"
                                />
                              </div>
                              <div>
                                <div className="text-base text-slate-400">
                                  {tool.properties.Category.select.name}
                                </div>
                                <div className="text-2xl font-medium">
                                  {tool.properties.Name.title[0].plain_text}
                                </div>
                              </div>
                            </div>
                            <div className="flex flex-col gap-y-2">
                              <span className="text-xl text-justify">
                                {
                                  tool.properties.Description.rich_text[0]
                                    .plain_text
                                }
                              </span>
                            </div>
                          </a>
                        )}
                      </>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
              {data ? (
                <div className="pt-16">
                  <div className="grid justify-center grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 min-w-72">
                    {data.map((tool: any, key: number) => (
                      <>
                        {tool.properties.Category.select.name === "Auto ML" && (
                          <a
                            href={tool.properties.Link.url}
                            key={key}
                            className="flex flex-col gap-y-4 p-4 text-slate-600 rounded-lg hover:bg-slate-100"
                          >
                            <div className="flex gap-x-4">
                              <div className="bg-slate-100 w-14 h-14 flex justify-center items-center rounded-lg border border-slate-200">
                                <img
                                  src={
                                    tool.properties.Logo.files[0].external.url
                                  }
                                  alt="Logo"
                                  className="w-8 h-8"
                                />
                              </div>
                              <div>
                                <div className="text-base text-slate-400">
                                  {tool.properties.Category.select.name}
                                </div>
                                <div className="text-2xl font-medium">
                                  {tool.properties.Name.title[0].plain_text}
                                </div>
                              </div>
                            </div>
                            <div className="flex flex-col gap-y-2">
                              <span className="text-xl text-justify">
                                {
                                  tool.properties.Description.rich_text[0]
                                    .plain_text
                                }
                              </span>
                            </div>
                          </a>
                        )}
                      </>
                    ))}
                  </div>
                </div>
              ) : (
                "Tools coming soon..."
              )}
            </Tabs>
          </div>
        </div>
      </div>
      <CTASection />
      <Footer />
    </div>
  );
}

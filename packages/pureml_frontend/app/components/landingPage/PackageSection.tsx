import { CodeBlock, sunburst } from "react-code-blocks";
import LandingPgTab from "./Tabs";

export default function PackageSection() {
  return (
    <div className="h-fit flex flex-col gap-y-6 pt-16 md:py-16">
      <h1 className="flex items-center text-3xl md:text-4xl lg:text-5xl !text-slate-400">
        04
      </h1>
      <div className="flex flex-col gap-y-12 text-slate-600">
        <div className="flex flex-col gap-y-6">
          <div className="md:w-3/4">
            <h1 className="font-medium text-3xl md:text-4xl lg:text-5xl pb-2">
              PureML-package
            </h1>
          </div>
          <div className="flex flex-col md:flex-row gap-x-12 gap-y-6">
            <div className="md:w-1/2 flex flex-col gap-y-12">
              <div className="flex flex-col gap-y-6">
                <h2 className="text-lg md:text-xl lg:text-2xl">
                  PureML is a versatile tool that allows you to package your
                  machine learning models into a standard, production-ready
                  container. Additionally, you can utilize a user-friendly web
                  interface to demonstrate your machine learning model, making
                  it easily accessible to anyone, from anywhere.
                </h2>
              </div>
            </div>
            <div className="md:w-1/2">
              <LandingPgTab
                tab1="DOCKER"
                tab2="FAST API"
                tab3="STREAMLIT"
                tab1Content={
                  <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                    <CodeBlock
                      text={`pureml.docker.create(“pet_classifier:dev:v1”)`}
                      language="python"
                      theme={sunburst}
                      showLineNumbers={false}
                      wrapLines
                    />
                  </div>
                }
                tab2Content={
                  <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                    <CodeBlock
                      text={`pureml.fastapi.create(“pet_classifier:dev:v1”)`}
                      language="python"
                      theme={sunburst}
                      showLineNumbers={false}
                      wrapLines
                    />
                  </div>
                }
                tab3Content={
                  <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                    <CodeBlock
                      text={`pureml.streamlit.create(“pet_classifier:dev:v1”)`}
                      language="python"
                      theme={sunburst}
                      showLineNumbers={false}
                      wrapLines
                    />
                  </div>
                }
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

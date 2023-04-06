import { CodeBlock, sunburst } from "react-code-blocks";

export default function DeploySection() {
  return (
    <div className="h-fit flex flex-col gap-y-6 pt-16 md:py-16">
      <h1 className="flex items-center text-3xl md:text-4xl lg:text-5xl !text-slate-400">
        05
      </h1>
      <div className="flex flex-col gap-y-12 text-slate-600">
        <div className="flex flex-col gap-y-6">
          <div className="md:w-3/4">
            <h1 className="font-medium text-3xl md:text-4xl lg:text-5xl pb-2">
              PureML-deploy
            </h1>
            <h2 className="text-lg md:text-xl lg:text-3xl">
              PureML gives you the ability to deploy machine learning models
              without the need for managing infrastructure or servers.
            </h2>
          </div>
          <div className="codeblock w-full overflow-hidden md:overflow-visible">
            <CodeBlock
              text={`$ pureml deploy pet_classifier:dev:v1`}
              language="python"
              theme={sunburst}
              showLineNumbers={false}
              wrapLines
            />
          </div>
        </div>
      </div>
    </div>
  );
}

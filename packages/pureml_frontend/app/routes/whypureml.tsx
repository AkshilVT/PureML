import type { MetaFunction } from "@remix-run/node";
import { CodeBlock, sunburst } from "react-code-blocks";
import CTASection from "~/components/landingPage/CTASection";
import Footer from "~/components/landingPage/Footer";
import Navbar from "~/components/landingPage/Navbar";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Why PureML | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export default function WhyPureML() {
  return (
    <div className="bg-slate-50 landingpg-font flex flex-col justify-center">
      <div>
        <Navbar />
        <div className="flex justify-center whypuremlbg">
          <img
            src="/imgs/landingPage/WhyPureMLHeroBG.svg"
            alt="WhyPureMLBG"
            className="hidden md:flex w-full"
          />
          <img
            src="/imgs/landingPage/MobWhyPureMLHeroBG.svg"
            alt="WhyPureMLBG"
            className="md:hidden w-full"
          />
        </div>
      </div>
      <div className="bg-slate-50 flex justify-center">
        <div className="flex flex-col gap-y-16 md:gap-y-32 py-16 md:max-w-screen-xl px-4 md:px-8 text-slate-600">
          <div className="flex flex-col gap-y-4 md:gap-y-6">
            <h1 className="text-3xl md:text-4xl lg:text-5xl">Why PureML?</h1>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              As machine learning (ML) becomes more and more pervasive across
              industries, there is an increasing need for version control
              systems that can handle the unique challenges posed by ML.
              Unfortunately, the current versioning system, based on git, falls
              short in this regard.
            </div>
          </div>
          <img src="/imgs/landingPage/PureMLSoln.svg" alt="PureMLSoln" />
          <div className="flex flex-col gap-y-4 md:gap-y-6">
            <h1 className="text-2xl md:text-3xl lg:text-4xl">Solution</h1>
            <div className="flex flex-col gap-y-12 md:gap-y-24">
              <div className="flex flex-col md:flex-row gap-y-4 gap-x-16">
                <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-1/2">
                  <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                    Object based semantic versioning for data and models:
                  </span>
                  <br />
                  PureML has a readable versioning format for storing/
                  retrieving models and datasets. It can store key value pairs
                  for model/dataset.
                </div>
                <img
                  src="/imgs/landingPage/Soln1.svg"
                  alt="PureMLSoln"
                  className="md:w-1/2"
                />
              </div>
              <div className="flex flex-col md:flex-row gap-y-4 gap-x-16">
                <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-1/2">
                  <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                    In-code orchestration:
                  </span>
                  <br />
                  Using PureML decorators, workflow can be orchestrated to
                  generate lineage, request cloud resources, and many more.
                </div>
                <div className="md:w-1/2">
                  <div className="codeblock w-full overflow-hidden md:overflow-visible">
                    <div className="relative">
                      <div className="overflow-hidden">
                        <div className="bg-yellow-400 opacity-30 h-[3.5rem] z-30 w-full absolute mt-5"></div>
                        <div className="bg-yellow-400 opacity-30 h-8 z-30 w-full absolute mt-[11.25rem]"></div>
                        <div className="bg-yellow-400 opacity-30 h-8 z-30 w-full absolute mt-[21.25rem]"></div>
                      </div>
                    </div>
                    <CodeBlock
                      text={`@resources(mem=”1GB”)
@cloud
def remote_function():
    ...
    ...

@dataset(label=”data_1:dev”)
def generate_features():
    ...
    ...


@model(label=”model_1:dev”)
def train_model():
    ...
    ...
`}
                      language="python"
                      theme={sunburst}
                      showLineNumbers={false}
                      wrapLines
                    />
                  </div>
                </div>
              </div>
              <div className="flex flex-col md:flex-row gap-y-4 gap-x-16">
                <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-1/2">
                  <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                    ML Evaluation:
                  </span>
                  <br />
                  PureML has an evaluation system to automate the evaluation of
                  any model with any dataset. This evaluation system can be
                  customized by user.
                </div>
                <img
                  src="/imgs/landingPage/Soln2.svg"
                  alt="PureMLSoln"
                  className="md:w-1/2"
                />
              </div>
            </div>
          </div>
          <div className="flex flex-col gap-y-4 md:gap-y-6">
            <h1 className="text-2xl md:text-3xl lg:text-4xl">
              Standout features
            </h1>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                Version Control:
              </span>
              <br />
              With Pureml, you can easily version your models and datasets,
              ensuring that you have a clear record of all changes made
              throughout the development process.
            </div>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                Commit Process:
              </span>
              <br />
              Pureml incorporates a review commit process that helps you ensure
              that only the best and most reliable models and datasets are
              shipped to your customers.
            </div>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                Packaging:
              </span>
              <br />
              Whether you need to package your models into Docker, Gradio, or
              FastAPI, Pureml has you covered. Our platform allows you to easily
              package your models in the format that works best for your
              specific needs.
            </div>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                Data Lineage:
              </span>
              <br />
              Pureml gives you full visibility into your datasets' lineage,
              making it easy to trace back any issues that may arise during
              development or deployment.
            </div>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                Branches:
              </span>
              <br />
              Pureml provides branches for both models and datasets, enabling
              teams to work on multiple versions of a model or dataset
              simultaneously.
            </div>
            <div className="text-lg md:text-xl lg:text-3xl text-justify md:w-4/5">
              <span className="font-medium text-slate-950 text-lg md:text-xl lg:text-3xl">
                Testing:
              </span>
              <br />
              Pureml includes testing for ML, making it easy to ship models
              reliably to your customers. With Pureml, you can be confident that
              your models will work as expected, every time.
            </div>
          </div>
          <div className="flex flex-col gap-y-4 md:gap-y-6">
            <h1 className="text-2xl md:text-3xl lg:text-4xl">
              How PureML makes it easy
            </h1>
            <img src="/imgs/landingPage/WithPureML.svg" alt="WithPureML" />
          </div>
        </div>
      </div>
      <CTASection />
      <Footer />
    </div>
  );
}

// ############################ error boundary ###########################

export function ErrorBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/FunctionalError.gif" alt="Error" width="500" />
    </div>
  );
}

export function CatchBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/FunctionalError.gif" alt="Error" width="500" />
    </div>
  );
}

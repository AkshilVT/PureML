import LandingPgTab from "./Tabs";

export default function ReviewSection() {
  return (
    <div className="flex flex-col gap-y-6 pt-16 md:py-16">
      <h1 className="flex items-center text-3xl md:text-4xl lg:text-5xl !text-slate-400">
        03
      </h1>
      <div className="flex flex-col gap-y-12 text-slate-600">
        <div className="flex flex-col gap-y-6">
          <div className="md:w-3/4">
            <h1 className="font-medium text-3xl md:text-4xl lg:text-5xl pb-2">
              Review
            </h1>
            <h2 className="text-lg md:text-xl lg:text-3xl">
              By providing a comprehensive set of metrics and visualizations,
              PureML makes it easy to identify and correct any issues with its
              review feature and allows you to evaluate the quality of their
              data and the accuracy of their model.
            </h2>
          </div>
          <div className="text-base md:text-lg lg:text-2xl">
            <LandingPgTab
              tab1="DATASET"
              tab2="MODEL"
              tab1Content={
                <img
                  src="/imgs/landingPage/ReviewDataset.svg"
                  alt="ReviewDataset"
                />
              }
              tab2Content={
                <img
                  src="/imgs/landingPage/ModelReview.svg"
                  alt="ReviewModel"
                />
              }
              tab3=""
              tab3Content=""
            />
          </div>
        </div>
      </div>
    </div>
  );
}

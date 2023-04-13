import { useState } from "react";
import Carousel from "react-simply-carousel";

export default function TrustedByCompaniesSection() {
  const [activeSlide, setActiveSlide] = useState(0);
  return (
    <div className="flex justify-center items-center">
      <div className="lg:h-fit w-full md:max-w-screen-xl px-0 md:px-8">
        <div className="px-4 md:px-0">
          <section className="">
            <Carousel
              containerProps={{
                style: {
                  justifyContent: "space-between",
                  userSelect: "none",
                },
              }}
              responsiveProps={[
                { minWidth: 1024, itemsToShow: 6 },
                { minWidth: 640, maxWidth: 1024, itemsToShow: 5 },
                { maxWidth: 640, itemsToShow: 4 },
              ]}
              swipeTreshold={60}
              activeSlideIndex={activeSlide}
              onRequestChange={setActiveSlide}
              dotsNav={{
                show: false,
              }}
              speed={8000}
              autoplay={true}
              autoplayDelay={20}
              easing="linear"
              infinite
            >
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/OneImmerseComp.svg"
                  alt="OneImmersive"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/GrpLandmarkComp.svg"
                  alt="GroupLandmark"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/PhenomPeopleComp.svg"
                  alt="PhenomPeople"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/LimeChatLogoComp.svg"
                  alt="LimeChat"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/CoutureLogoComp.svg"
                  alt="Couture"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/Pe2proLogoComp.svg"
                  alt="Pe2pro"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/ActalystLogoComp.svg"
                  alt="Actalyst"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/NeuralSyncComp.svg"
                  alt="NeuralSync"
                />
              </div>
              <div className="w-48 md:w-56 flex justify-center items-center">
                <img
                  src="/imgs/landingPage/trustedCompanyLogo/LotusDeuComp.svg"
                  alt="LotusDeu"
                />
              </div>
            </Carousel>
          </section>
        </div>
      </div>
    </div>
  );
}

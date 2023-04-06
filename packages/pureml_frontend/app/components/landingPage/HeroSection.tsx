export default function HeroSection() {
  return (
    <div className="flex justify-center items-center">
      <div className="w-full md:max-w-screen-xl px-0 md:px-8">
        <div className="px-4 md:px-0">
          <div className="flex flex-col gap-y-12 md:w-3/4">
            <div className="flex flex-col gap-y-4">
              <h1 className="font-medium !leading-snug text-4xl md:text-5xl lg:text-6xl">
                The next-gen developer platform for{" "}
                <u className="underline underline-offset-8 font-medium !leading-snug text-4xl md:text-5xl lg:text-6xl font-satoshi">
                  Production ML
                </u>
                .
              </h1>
              <h1 className="font-medium !leading-snug text-lg md:text-xl lg:text-2xl">
                Develop fast, ship with confidence and scale without limits.
              </h1>
            </div>
            <div className="flex flex-col md:flex-row items-center gap-y-4 md:gap-x-6">
              <a
                href="https://tally.so/r/wa96xv"
                className="w-full md:w-fit"
                target="_blank"
                rel="noreferrer"
              >
                <button className="btn btn-primary btn-sm font-normal text-white w-full md:w-fit hover-effect px-4 rounded-lg text-lg letterSpaced">
                  JOIN WAITLIST
                </button>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

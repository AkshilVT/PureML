export default function CTASection() {
  return (
    <div className="bg-slate-50 px-4 py-16 md:p-24 h-fit flex flex-col justify-center items-center">
      <div className="md:w-[48rem] lg:w-[64rem] flex flex-col gap-y-6">
        <h1 className="flex justify-center items-center font-medium text-center !leading-snug !text-3xl md:!text-5xl lg:!text-7xl text-brand-200">
          PureML empowers everyone to work together, and ship with confidence.
        </h1>
        <a
          href="https://tally.so/r/wa96xv"
          target="_blank"
          rel="noreferrer"
          className="flex justify-center items-center"
        >
          <button
            type="submit"
            className="btn btn-primary btn-sm font-normal text-white w-full md:w-fit hover-effect px-8 rounded-lg text-xl letterSpaced"
          >
            JOIN WAITLIST
          </button>
        </a>
      </div>
    </div>
  );
}

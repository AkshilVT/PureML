import Button from "../ui/Button";

export default function JoinCommunitySection() {
  return (
    <div className="h-fit flex flex-col justify-center items-center gap-y-8 lg:gap-y-16">
      <div className="flex flex-col gap-y-4 justify-center">
        <div className="text-center text-brand-200 font-satoshi">
          <span className="text-3xl md:text-4xl lg:text-5xl">Join our </span>
          <span className="text-slate-500 text-3xl md:text-4xl lg:text-5xl text-slate-500">
            open source
          </span>
          <span className="text-3xl md:text-4xl lg:text-5xl"> community</span>
        </div>
        <div className="flex justify-center text-slate-600">
          <div className="xl:w-3/5 text-center text-base md:text-xl lg:text-2xl font-satoshi">
            We're building the largest production ML community on the Internet.
            Check out our{" "}
            <a
              href="https://pureml.height.app/roadmap"
              target="_blank"
              rel="noreferrer"
              className="underline underline-offset-4 text-slate-400 text-base md:text-xl lg:text-2xl"
            >
              Public Roadmap
            </a>{" "}
            and leave comments!
          </div>
        </div>
      </div>
      <div className="grid md:grid-cols-3 gap-y-8 gap-x-8">
        <div className="flex flex-col gap-y-2 justify-center p-8 bg-slate-100 rounded-2xl">
          <div className="flex justify-center text-3xl md:text-4xl lg:text-5xl">
            06+
          </div>
          <div className="flex justify-center text-base md:text-xl lg:text-2xl">
            Contributors
          </div>
          <a
            href="https://github.com/PureML-Inc"
            target="_blank"
            rel="noreferrer"
          >
            <Button intent="secondary">Build PureML</Button>
          </a>
        </div>
        <div className="flex flex-col gap-y-2 justify-center p-8 bg-slate-100 rounded-2xl">
          <div className="flex justify-center text-3xl md:text-4xl lg:text-5xl">
            100+
          </div>
          <div className="flex justify-center text-base md:text-xl lg:text-2xl">
            Github Stars
          </div>
          <a
            href="https://github.com/PureML-Inc/PureML"
            target="_blank"
            rel="noreferrer"
          >
            <Button intent="secondary">Browse Github</Button>
          </a>
        </div>
        <div className="flex flex-col gap-y-2 justify-center p-8 bg-slate-100 rounded-2xl">
          <div className="flex justify-center text-3xl md:text-4xl lg:text-5xl">
            100+
          </div>
          <div className="flex justify-center text-center text-base md:text-xl lg:text-2xl">
            Community members
          </div>
          <a
            href="https://discord.gg/xNUHt9yguJ"
            target="_blank"
            rel="noreferrer"
          >
            <Button intent="secondary">Join Discord</Button>
          </a>
        </div>
      </div>
    </div>
  );
}

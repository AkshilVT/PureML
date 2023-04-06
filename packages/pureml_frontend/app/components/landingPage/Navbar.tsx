import { useState } from "react";
import { ArrowUpRight, Github, Menu, X } from "lucide-react";

export default function Navbar() {
  const [open, setOpen] = useState(true);
  if (!open)
    return (
      <div className="flex justify-center">
        <div className="flex justify-center w-full backdrop-blur-sm bg-white/30">
          <div className="lg:h-fit w-full md:max-w-screen-xl px-0 md:px-8">
            <div className="sm:px-24">
              <div className="rounded-b-2xl border-b border-slate-200">
                <div className="flex px-4 py-5 md:px-12 justify-between items-center">
                  <a href="/">
                    <img
                      src="/PureMLLogoWText.svg"
                      alt="PureMLLogo"
                      width="96"
                      height="96"
                    />
                  </a>
                  <X
                    className="sm:hidden text-slate-950 cursor-pointer w-8 h-8"
                    onClick={() => setOpen(!open)}
                  />
                </div>
                <div className="flex flex-col gap-y-4 p-4 font-medium text-brand-200">
                  <div className="flex items-center">
                    <a href="/whypureml" className="w-max letterSpaced">
                      WHY PUREML
                    </a>
                  </div>
                  <div className="flex items-center">
                    <a href="/mltools" className="w-max letterSpaced">
                      MLOPS TOOLS
                    </a>
                  </div>
                  <div className="flex items-center">
                    <a
                      href="https://changelog.pureml.com"
                      target="_blank"
                      rel="noreferrer"
                      className="w-max letterSpaced"
                    >
                      CHANGELOG
                    </a>
                  </div>
                  <div className="flex items-center">
                    <a
                      href="https://pureml.mintlify.app"
                      target="_blank"
                      rel="noreferrer"
                      className="w-max letterSpaced"
                    >
                      DOCS
                    </a>
                  </div>
                  <div className="flex items-center text-brand-200">
                    <a href="/auth/signin" className="w-max letterSpaced">
                      SIGN IN
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  return (
    <div className="flex justify-center">
      <div className="flex justify-center w-full backdrop-blur-sm bg-white/30">
        <div className="lg:h-fit w-full md:max-w-screen-xl px-0 md:px-8">
          <div className="flex p-4 md:px-0 justify-between items-center w-full">
            <a href="/">
              <img
                src="/PureMLLogoWText.svg"
                alt="PureMLLogo"
                width="96"
                height="96"
              />
            </a>
            <div className="md:hidden flex">
              <div className="flex items-center pr-8">
                <a
                  href="https://github.com/pureml-inc/pureml"
                  className="flex gap-x-2 border border-slate-950 rounded-full py-2 px-4 text-brand-200 text-lg letterSpaced hover:bg-slate-100"
                >
                  <Github className="w-5" />
                  STAR
                </a>
              </div>
              <div className="flex justify-center items-center">
                <Menu
                  className="sm:hidden text-slate-950 cursor-pointer w-8 h-8"
                  onClick={() => setOpen(!open)}
                />
              </div>
            </div>
            <div className="hidden sm:flex flex font-medium">
              <div className="px-4 flex justify-center items-center">
                <a
                  href="/whypureml"
                  className="w-full text-lg letterSpaced text-brand-200 hover:underline hover:underline-offset-2"
                >
                  WHY PUREML
                </a>
              </div>
              <div className="px-4 flex justify-center items-center">
                <a
                  href="/mltools"
                  className="w-full text-lg letterSpaced text-brand-200 hover:underline hover:underline-offset-2"
                >
                  MLOPS TOOLS
                </a>
              </div>
              <div className="px-4 flex justify-center items-center">
                <a
                  href="https://changelog.pureml.com"
                  className="w-full text-lg letterSpaced text-brand-200 hover:underline hover:underline-offset-2"
                  target="_blank"
                  rel="noreferrer"
                >
                  CHANGELOG
                </a>
              </div>
              <div className="px-4 flex justify-center items-center">
                <a
                  href="https://pureml.mintlify.app"
                  className="w-full text-lg letterSpaced text-brand-200 hover:underline hover:underline-offset-2"
                  target="_blank"
                  rel="noreferrer"
                >
                  DOCS
                </a>
              </div>
              <div className="px-4 flex justify-center items-center">
                <a
                  href="/auth/signin"
                  className="flex justify-center items-center gap-y-4 text-brand-200 hover:underline hover:underline-offset-2 text-lg letterSpaced"
                >
                  SIGN IN <ArrowUpRight className="text-brand-200 w-5" />
                </a>
              </div>
              <div className="px-4 flex justify-center items-center">
                <a
                  href="https://github.com/pureml-inc/pureml"
                  className="flex gap-x-2 border border-slate-950 rounded-full py-2 px-4 text-brand-200 text-lg letterSpaced hover:bg-slate-100"
                >
                  <Github className="w-5" />
                  STAR
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

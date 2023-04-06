export default function Footer() {
  return (
    <div className="bg-slate-50">
      <div className="flex justify-center">
        <div className="flex flex-col lg:flex-row justify-between gap-y-8 px-4 py-8 md:px-8 md:py-12 xl:px-0 text-slate-600 !text-lg md:!text-xl w-full max-w-screen-xl">
          <div className="flex flex-col gap-y-4 md:gap-y-6">
            <a href="/" className="text-xl xl:text-2xl">
              © 2022, PureML Inc
            </a>
            <div className="flex gap-x-8 items-center">
              <a
                href="https://discord.com/invite/xNUHt9yguJ"
                target="_blank"
                rel="noreferrer"
              >
                <img
                  src="/imgs/landingPage/icons/DiscordIcon.svg"
                  alt="Discord"
                  className="w-6"
                />
              </a>
              <a
                href="https://github.com/PuremlHQ/PureML"
                target="_blank"
                rel="noreferrer"
              >
                <img
                  src="/imgs/landingPage/icons/GitHubIcon.svg"
                  alt="Github"
                  className="w-6"
                />
              </a>
              <a
                href="https://twitter.com/puremlHQ"
                target="_blank"
                rel="noreferrer"
              >
                <img
                  src="/imgs/landingPage/icons/TwitterIcon.svg"
                  alt="Twitter"
                  className="w-6"
                />
              </a>
              <a
                href="https://www.linkedin.com/company/pureml-inc/"
                target="_blank"
                rel="noreferrer"
              >
                <img
                  src="/imgs/landingPage/icons/LinkedinIcon.svg"
                  alt="LinkedIn"
                  className="w-6"
                />
              </a>
              <a
                href="mailto:contact@pureml.com"
                target="_blank"
                rel="noreferrer"
              >
                <img
                  src="/imgs/landingPage/icons/MailIcon.svg"
                  alt="Mail"
                  className="w-6"
                />
              </a>
            </div>
          </div>
          <div className="flex justify-between md:gap-x-96">
            <div className="flex flex-col justify-between gap-y-2">
              <a
                href="https://pureml.mintlify.app"
                target="_blank"
                rel="noreferrer"
                className="underline text-lg xl:text-xl"
              >
                Docs
              </a>
              <a href="/whypureml" className="underline text-lg xl:text-xl">
                Why PureML
              </a>
              <a href="/mltools" className="underline text-lg xl:text-xl">
                ML Tools
              </a>
              <a
                href="https://pureml.notion.site/7de13568835a4cf18913307503a2cdd4?v=82199f96833a48e5907023c8a8d565c6"
                target="_blank"
                rel="noreferrer"
                className="underline text-lg xl:text-xl"
              >
                Roadmap
              </a>
              <a
                href="https://pureml.notion.site/PureML-Changelog-096f7541dd6245c0a3c244e9f216ad31"
                target="_blank"
                rel="noreferrer"
                className="underline text-lg xl:text-xl"
              >
                Changelog
              </a>
            </div>
            <div className="flex flex-col gap-x-2">
              <a
                href="https://discord.gg/xNUHt9yguJ"
                target="_blank"
                rel="noreferrer"
                className="text-lg xl:text-xl"
              >
                Join Discord
              </a>
              <a
                href="https://pureml.instatus.com"
                target="_blank"
                rel="noreferrer"
                className="text-lg xl:text-xl"
              >
                Status
              </a>
              <div
                // href="/TermsAndCondition.pdf"
                // target="_blank"
                className="text-lg xl:text-xl"
              >
                Terms & Conditions
              </div>
              <div
                // href="/PrivacyPolicy.pdf"
                // target="_blank"
                className="text-lg xl:text-xl"
              >
                Privacy Policy
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

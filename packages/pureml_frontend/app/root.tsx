import type { LinksFunction, MetaFunction } from "@remix-run/node";
import { redirect } from "@remix-run/node";
import {
  Links,
  LiveReload,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useCatch,
  useLoaderData,
  useLocation,
  useMatches,
  useNavigate,
} from "@remix-run/react";
import type { ReactNode } from "react";
import { useEffect } from "react";
import CTASection from "./components/landingPage/CTASection";
import Footer from "./components/landingPage/Footer";
import HeroSection from "./components/landingPage/HeroSection";
import Navbar from "./components/landingPage/Navbar";
import PackageSection from "./components/landingPage/PackageSection";
import { fetchUserSettings } from "./routes/api/auth.server";
import { destroySession, getSession } from "./session";
import * as gtag from "~/analytics/gtags.client";

// Tailwind CSS---------------------------------------------------------------
import styles from "./styles/app.css";
import style from "reactflow/dist/style.css";
import contributionStyle from "react-calendar-heatmap/dist/styles.css";
import toastStyle from "react-toastify/dist/ReactToastify.css";
import { Slide, ToastContainer } from "react-toastify";
import TrustedByCompanies from "./components/landingPage/TrustedByCompanies";
import VersionSection from "./components/landingPage/VersionSection";
import ReviewSection from "./components/landingPage/ReviewSection";
import TestingSection from "./components/landingPage/TestingSection";
import DeploySection from "./components/landingPage/DeploySection";
import posthog from "posthog-js";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: styles },
  { rel: "stylesheet", href: style },
  { rel: "stylesheet", href: toastStyle },
  { rel: "stylesheet", href: contributionStyle },
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono&family=IBM+Plex+Sans:wght@400;500&display=swap",
  },
  {
    rel: "stylesheet",
    href: "https://api.fontshare.com/v2/css?f[]=satoshi@500,400&display=swap",
  },
];
// ---------------------------------------------------------------------------

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "PureML",
  description: "Developer Platform for Production ML.",
  "og:url": "https://pureml.com",
  "og:title": "PureML",
  "og:description": "Developer Platform for Production ML.",
  "og:image": "https://www.pureml.com/PureMLMetaImg.svg",
  "twitter:url": "https://pureml.com",
  "twitter:title": "PureML",
  "twitter:description": "Developer Platform for Production ML.",
  "twitter:image": "https://www.pureml.com/PureMLMetaImg.svg",
  viewport: "width=device-width, height=device-height, initial-scale=1",
});

export async function loader({ request }: any) {
  const cookieHeader = request.headers.get("cookie");
  const cacheControlHeader = request.headers.get("cache-control");
  if (cookieHeader && cacheControlHeader) {
    const session = await getSession(request.headers.get("Cookie"));
    const accesstoken = session.get("accessToken");
    if (accesstoken) {
      const profile = await fetchUserSettings(accesstoken);
      if (!profile)
        return redirect("/", {
          headers: {
            "Set-Cookie": await destroySession(session),
          },
        });
      return { profile, gaTrackingId: process.env.GA_TRACKING_ID };
    }
  }
  return { profile: null, gaTrackingId: process.env.GA_TRACKING_ID };
}

export default function App() {
  const location = useLocation();
  const navigate = useNavigate();
  const matches = useMatches();
  const pathname = matches[1].pathname;
  const profile = useLoaderData();
  const { prof, gaTrackingId } = useLoaderData();
  useEffect(() => {
    posthog.init("phc_u5aRMx559YsrAaBkCIWGoxEXepwpQBxvhUncdb9giP5", {
      api_host: "https://app.posthog.com",
    });
  });

  useEffect(() => {
    if (profile.profile && pathname === "/") return navigate("/models");
  }, [profile]);

  useEffect(() => {
    if (gaTrackingId?.length && location) {
      gtag.pageview(location.pathname, gaTrackingId);
    }
  }, [location, gaTrackingId]);
  return (
    <Document gaTrackingId={gaTrackingId}>
      <ToastContainer
        theme="light"
        position="bottom-right"
        autoClose={5000}
        transition={Slide}
        hideProgressBar
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
      {!prof && pathname === "/" ? (
        <>
          {!prof && (
            <div className="bg-slate-50 landingpg-font flex flex-col justify-center">
              <div className="flex flex-col justify-between h-screen bgvideo bg-cover">
                <Navbar />
                <HeroSection />
                <div className="md:gap-y-48 pb-24 md:pb-24 2xl:pt-32">
                  <TrustedByCompanies />
                </div>
              </div>

              <div className="bg-slate-50 flex justify-center">
                <div className="md:max-w-screen-xl px-4 md:px-8">
                  <VersionSection />
                  <TestingSection />
                  <ReviewSection />
                  <PackageSection />
                  <DeploySection />
                </div>
              </div>
              {/* <div className="xl:flex xl:justify-center overflow-hidden pt-14 md:py-20">
                <div className="md:max-w-screen-xl">
                  <TestimonialsSection />
                  <JoinCommunitySection />
                </div>
              </div> */}
              <CTASection />
              <Footer />
            </div>
          )}
        </>
      ) : (
        ""
      )}
      <Outlet />
    </Document>
  );
}

function Document({
  children,
  gaTrackingId,
}: {
  children: ReactNode;
  gaTrackingId: string;
}) {
  return (
    <html lang="en">
      <head>
        <meta name="title" content="PureML" />
        <meta
          name="description"
          content="Developer platform for Production ML."
        />
        <meta property="og:type" content="website" />
        <meta property="og:url" content="https://pureml.com" />
        <meta property="og:title" content="PureML" />
        <meta
          property="og:description"
          content="Developer platform for Production ML."
        />
        <meta
          property="og:image"
          content="https://www.pureml.com/PureMLMetaImg.svg"
        />
        <meta property="twitter:card" content="summary_large_image" />
        <meta property="twitter:url" content="https://pureml.com" />
        <meta property="twitter:title" content="PureML" />
        <meta
          property="twitter:description"
          content="Developer platform for Production ML."
        />
        <meta
          property="twitter:image"
          content="https://www.pureml.com/PureMLMetaImg.svg"
        />
        <Meta />
        <Links />
      </head>
      <body>
        {!gaTrackingId ? null : (
          <>
            <script
              async
              src={`https://www.googletagmanager.com/gtag/js?id=${gaTrackingId}`}
            />
            <script
              async
              id="gtag-init"
              dangerouslySetInnerHTML={{
                __html: `
                window.dataLayer = window.dataLayer || [];
                function gtag(){dataLayer.push(arguments);}
                gtag('js', new Date());
                gtag('config', '${gaTrackingId}', {
                  page_path: window.location.pathname,
                });
              `,
              }}
            />
          </>
        )}
        <Scripts />
        <div className="bg-slate-50">{children}</div>
        <ScrollRestoration />
        <LiveReload />
      </body>
    </html>
  );
}

// ############################ error boundary ###########################

export function ErrorBoundary({ error }: { error: Error }) {
  return (
    <Document gaTrackingId={""}>
      {/* <div className="p-12">
        <span className="text-3xl font-medium">Error</span>
        <p>{error.message}</p>
        <div className="text-xl pt-8 font-medium">The stack trace is:</div>
        <pre className="whitespace-pre-wrap">{error.stack}</pre> */}
      <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
        <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
        <div className="text-3xl text-slate-600 font-medium">
          Something went wrong :(
        </div>
        <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
      </div>
      {/* </div> */}
    </Document>
  );
}

export function CatchBoundary() {
  const caught = useCatch();

  return (
    <Document gaTrackingId={""}>
      {/* <div className="p-12">
        <span className="text-3xl font-medium">Status: {caught?.status}</span>
        <div className="text-xl pt-8 font-medium">Data: {caught?.data}</div> */}
      <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
        <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
        <div className="text-3xl text-slate-600 font-medium">
          Something went wrong :(
        </div>
        <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
      </div>
      {/* </div> */}
    </Document>
  );
}

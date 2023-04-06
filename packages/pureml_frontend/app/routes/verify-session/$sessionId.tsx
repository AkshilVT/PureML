import type { MetaFunction } from "@remix-run/node";
import { useLoaderData, useNavigate } from "@remix-run/react";
import { Suspense, useEffect } from "react";
import { toast } from "react-toastify";
import Loader from "~/components/ui/Loading";
import { getSession } from "~/session";
import { fetchUserSettings, fetchVerifySession } from "../api/auth.server";
import NavBar from "~/components/Navbar";
import { fetchOrgDetails } from "../api/org.server";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Verify Session to Login | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const accesstoken = session.get("accessToken");
  const sessionId = params.sessionId;
  if (accesstoken) {
    const orgId = session.get("orgId");
    if (!orgId) return null;
    const profile = await fetchUserSettings(accesstoken);
    const orgDetails = await fetchOrgDetails(orgId, session.get("accessToken"));
    const verifySession = await fetchVerifySession(accesstoken, sessionId);
    return { profile, orgDetails, verifySession, sessionId };
  }
  return sessionId;
}

export default function VerifySession() {
  const data = useLoaderData();
  const navigate = useNavigate();

  useEffect(() => {
    if (!data.verifySession) {
      const timer = setTimeout(async () => {
        toast.info("User not logged in. Proceeding to signin!");
        navigate(`/auth/signin?sessionid=${data}`);
      }, 2000);
      return () => clearTimeout(timer);
    }
  }, [data, navigate]);

  if (data.verifySession) {
    if (
      data.verifySession.status === 200 ||
      data.verifySession.message === "Session already approved"
    )
      return (
        <Suspense fallback={<Loader />}>
          <div className="h-screen">
            {data ? (
              <NavBar
                intent="loggedIn"
                user={data.profile[0].name.charAt(0).toUpperCase()}
                orgName={
                  <a
                    href={`/org/${data.orgDetails[0].name}`}
                    className="flex items-center justify-center"
                  >
                    {data.orgDetails[0].name}
                  </a>
                }
                orgAvatarName={data.orgDetails[0].name.charAt(0)}
              />
            ) : (
              <NavBar
                intent="loggedOut"
                user=""
                orgName={
                  <a
                    href="/models"
                    className="flex items-center justify-center pr-8"
                  >
                    <img
                      src="/PureMLLogoWText.svg"
                      alt="Logo"
                      className="w-20"
                    />
                  </a>
                }
                orgAvatarName=""
              />
            )}
            <div className="flex flex-col gap-y-8 justify-center items-center h-4/5 text-slate-600 font-medium">
              <img
                src="/imgs/CLILoginSuccess.svg"
                alt="Success"
                className="w-96"
              />
              <div className="flex flex-col justify-center items-center gap-y-2">
                <div className="text-slate-600 text-2xl">Welcome to PureML</div>
                <div className="text-slate-400 text-lg">
                  {data.verifySession.message}
                </div>
                <div className="text-slate-400 text-lg">
                  You successfully logged in through CLI. You can close this
                  window now :)
                </div>
              </div>
            </div>
          </div>
        </Suspense>
      );
    else
      return (
        <Suspense fallback={<Loader />}>
          <div className="h-screen">
            {data ? (
              <NavBar
                intent="loggedIn"
                user={data.profile[0].name.charAt(0).toUpperCase()}
                orgName={
                  <a
                    href={`/org/${data.orgDetails[0].name}`}
                    className="flex items-center justify-center"
                  >
                    {data.orgDetails[0].name}
                  </a>
                }
                orgAvatarName={data.orgDetails[0].name.charAt(0)}
              />
            ) : (
              <NavBar
                intent="loggedOut"
                user=""
                orgName={
                  <a
                    href="/models"
                    className="flex items-center justify-center pr-8"
                  >
                    <img
                      src="/PureMLLogoWText.svg"
                      alt="Logo"
                      className="w-20"
                    />
                  </a>
                }
                orgAvatarName=""
              />
            )}
            <div className="flex flex-col gap-y-8 justify-center items-center h-4/5 text-slate-600 font-medium">
              <img
                src="/error/CLILoginError.svg"
                alt="Error"
                className="w-96"
              />
              <div className="flex flex-col justify-center items-center gap-y-2">
                <div className="text-slate-600 text-2xl">
                  {data.verifySession.message}
                </div>
                <div className="text-slate-400 text-lg">
                  Close this window and Please try logging in again.
                </div>
              </div>
            </div>
          </div>
        </Suspense>
      );
  }
}

// ############################ error boundary ###########################

export function ErrorBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
    </div>
  );
}

export function CatchBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
    </div>
  );
}

import { json, redirect } from "@remix-run/node";
import {
  Form,
  useActionData,
  useFetcher,
  useLoaderData,
  useMatches,
  useParams,
  useResolvedPath,
  useSearchParams,
} from "@remix-run/react";
import { Eye, EyeOff } from "lucide-react";
import { Suspense, useEffect, useState } from "react";
import { toast } from "react-toastify";
import NavBar from "~/components/Navbar";
import Button from "~/components/ui/Button";
import Input from "~/components/ui/Input";
import Link from "~/components/ui/Link";
import Loader from "~/components/ui/Loading";
import {
  fetchSignIn,
  fetchUserSettings,
  fetchVerifySession,
} from "~/routes/api/auth.server";
import { fetchAllOrgs, fetchOrgDetails } from "~/routes/api/org.server";
import { commitSession, getSession } from "~/session";

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const accesstoken = session.get("accessToken");
  const sessionId = new URL(request.url);
  const verifySessionId = sessionId.searchParams.get("sessionid");
  if (!verifySessionId && accesstoken) {
    return redirect("/models", {});
  } else if (verifySessionId && accesstoken) {
    return null;
  }
  return null;
}

export const action = async ({ request }: any) => {
  const form = await request.formData();
  const email = form.get("email");
  const password = form.get("password");
  const data = await fetchSignIn(email, password);
  const session = await getSession(request.headers.get("Cookie"));
  const sessionId = new URL(request.url);
  const verifySessionId = sessionId.searchParams.get("sessionid");
  const accessToken = data.data[0].accessToken;
  session.set("accessToken", accessToken);
  const org = await fetchAllOrgs(accessToken);
  session.set("orgId", org[0].org.uuid);
  session.set("orgName", org[0].org.name);
  if (!verifySessionId && data.message === "User logged in") {
    return json(data, {
      headers: {
        "Set-Cookie": await commitSession(session),
      },
    });
  } else if (verifySessionId && data.message === "User logged in") {
    if (verifySessionId && accessToken) {
      const orgId = session.get("orgId");
      if (!orgId) return null;
      const profile = await fetchUserSettings(accessToken);
      const orgDetails = await fetchOrgDetails(
        orgId,
        session.get("accessToken")
      );
      const verifySession = await fetchVerifySession(
        accessToken,
        verifySessionId
      );
      return json(
        { profile, orgDetails, verifySession, verifySessionId },
        {
          headers: {
            "Set-Cookie": await commitSession(session),
          },
        }
      );
    }
  }
  return { data, verifySessionId, ok: true };
};

export default function SignIn() {
  useLoaderData();
  const data = useActionData();
  const [show, setShow] = useState(false);
  useEffect(() => {
    if (!data) return;
    if (!data.verifySessionId) {
      if (data.message === "User logged in") {
        toast.success("Successfully signed in");
      } else if (data.data.message === "User email is not verified")
        toast.error("Email not verified! Please verify to proceed.");
      else if (data.data.message === "User not found")
        toast.error("User not found");
      else if (data.data.message === "Invalid username")
        toast.error("Invalid username!");
      else if (data.data.message === "Invalid credentials")
        toast.error("Invalid credentials!");
      else toast.error("Something went wrong!");
    }
  }, [data]);

  if (data) {
    if (data.verifySessionId) {
      if (
        data.verifySession.status === 200 ||
        data.verifySession.message === "Session already approved"
      ) {
        return (
          <Suspense fallback={<Loader />}>
            <div className="h-screen bg-slate-50">
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
                  <div className="text-slate-600 text-2xl">
                    Welcome to PureML
                  </div>
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
      } else {
        return (
          <Suspense fallback={<Loader />}>
            <div className="h-screen bg-slate-50">
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
  }

  return (
    <Suspense fallback={<Loader />}>
      <div className="flex justify-center">
        <div className="w-fit text-center">
          <div className="flex justify-center items-center pb-16">
            <img src="/PureMLLogoWText.svg" alt="Logo" className="w-36" />
          </div>
          <Form
            method="post"
            className="text-slate-600 flex flex-col text-left"
          >
            <div className="flex flex-col gap-y-12">
              <div className="flex flex-col gap-y-6">
                <label htmlFor="email" className="font-medium">
                  Email
                  <Input
                    intent="primary"
                    type="email"
                    name="email"
                    fullWidth={false}
                    aria-label="emailid"
                    data-testid="email-input1"
                    required
                  />
                </label>
                <label htmlFor="password" className="font-medium">
                  Password
                  <div className="input-icons">
                    <input
                      className="input-field rounded"
                      name="password"
                      aria-label="password"
                      data-testid="password-input1"
                      type={show ? "text" : "password"}
                      required
                    />
                    {show ? (
                      <i>
                        <Eye
                          className="text-slate-400 hover:text-slate-600 w-4 cursor-pointer"
                          onClick={() => setShow(!show)}
                        />
                      </i>
                    ) : (
                      <i>
                        <EyeOff
                          className="text-slate-400 hover:text-slate-600 w-4 cursor-pointer"
                          onClick={() => setShow(!show)}
                        />
                      </i>
                    )}
                  </div>
                </label>
              </div>
              <Button intent="primary">Sign in</Button>
            </div>
          </Form>
          <div className="flex items-center text-slate-600 justify-center mt-6">
            <Link intent="secondary" hyperlink="/auth/forgot_password">
              Forgot Password?
            </Link>
            {/* <p className="px-2 text-slate-400">|</p>
        <div className="flex items-center space-x-1 font-medium">
          <Link intent="secondary" hyperlink="/auth/signup">
            Sign Up
          </Link>
        </div> */}
          </div>
        </div>
      </div>
    </Suspense>
  );
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

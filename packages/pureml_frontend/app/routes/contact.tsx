import type { MetaFunction } from "@remix-run/node";
import { redirect } from "@remix-run/node";
import {
  Form,
  useActionData,
  useLoaderData,
  useNavigate,
} from "@remix-run/react";
import { X } from "lucide-react";
import { toast } from "react-toastify";
import Navbar from "~/components/Navbar";
import Button from "~/components/ui/Button";
import { getSession } from "~/session";
import { fetchUserSettings } from "./api/auth.server";
import { fetchOrgDetails } from "./api/org.server";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Contact Us | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader({ request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const accesstoken = session.get("accessToken");
  if (accesstoken) {
    const userProfile = await fetchUserSettings(accesstoken);
    const orgId = session.get("orgId");
    const orgDetails = await fetchOrgDetails(orgId, session.get("accessToken"));
    return { userProfile, orgDetails };
  }
  return null;
}

export const action = async ({ request }: any) => {
  const form = await request.formData();
  const comment = form.get("comment");
  const session = await getSession(request.headers.get("Cookie"));
  const accessToken = session.get("accessToken");
  const userDetails = await fetchUserSettings(accessToken);
  const email = userDetails[0].email;
  // const submit = base("tbl7qXTTBN3Ln6KyE").create(
  //   [
  //     {
  //       fields: {
  //         Email: email,
  //         Comment: comment,
  //       },
  //     },
  //   ],
  //   (err: string) => {
  //     if (err) {
  //       console.error(err);
  //     }
  //   }
  // );
  return redirect("/models", {});
};

export default function Contact() {
  const navigate = useNavigate();
  const data = useActionData();
  const userProfileData = useLoaderData();
  if (!userProfileData) return navigate("/");
  return (
    <div className="flex flex-col h-screen overflow-hidden overlay-bg">
      <div className="opacity-10">
        <Navbar
          intent="loggedIn"
          user={userProfileData.userProfile[0].name.charAt(0).toUpperCase()}
          orgName={
            <a
              href={`/org/${userProfileData.orgDetails[0].name}`}
              className="flex items-center justify-center"
            >
              {userProfileData.orgDetails[0].name}
            </a>
          }
          orgAvatarName={userProfileData.orgDetails[0].name.charAt(0)}
        />
      </div>
      <div className="w-full h-full flex justify-center items-center">
        <div className="bg-slate-0 p-4 rounded-lg w-[28rem] h-fit">
          <div className="flex justify-between text-slate-800 font-medium pb-8">
            Contact Us
            <X
              className="text-slate-400 w-4 cursor-pointer"
              onClick={() => {
                navigate("/models");
              }}
            />
          </div>
          <Form method="post">
            <label htmlFor="comment">
              Comment
              <textarea
                typeof="text"
                name="comment"
                required
                className="whitespace-pre-line w-full bg-transparent text-sm border border-slate-200 rounded-md h-full hover:border-slate-400 focus:outline-none focus:border-blue-250 max-h-[200px] p-4"
              />
            </label>
            <div className="flex justify-end pt-6">
              <div className="flex justify-between w-2/3 gap-x-6">
                <div className="w-1/2">
                  <Button
                    intent="secondary"
                    onClick={() => {
                      navigate("/models");
                    }}
                  >
                    Cancel
                  </Button>
                </div>
                <div className="w-1/2">
                  <Button
                    intent="primary"
                    type="submit"
                    onClick={() => {
                      toast.success("Query sent! We will be back shortly :)");
                    }}
                  >
                    Submit
                  </Button>
                </div>
              </div>
            </div>
          </Form>
        </div>
      </div>
    </div>
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

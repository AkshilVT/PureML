import type { MetaFunction } from "@remix-run/node";
import { Outlet, useSearchParams } from "@remix-run/react";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Join Us | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export default function AuthLayout() {
  let [searchParams] = useSearchParams();
  const session = searchParams.get("sessionid");
  return (
    <div className="w-screen h-screen flex items-center justify-center bg-slate-0">
      {!session && (
        <div className="bg-slate-0 md:flex md:flex-col justify-center items-center md:pt-0 text-white">
          <Outlet />
        </div>
      )}
      {session && (
        <div className="w-full">
          <Outlet />
        </div>
      )}
    </div>
  );
}

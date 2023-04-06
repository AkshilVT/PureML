import { Outlet } from "@remix-run/react";
import posthog from "posthog-js";
import { Suspense, useEffect } from "react";
import Loader from "~/components/ui/Loading";

export default function Index() {
  useEffect(() => {
    posthog.init("phc_u5aRMx559YsrAaBkCIWGoxEXepwpQBxvhUncdb9giP5", {
      api_host: "https://app.posthog.com",
    });
  });
  return (
    <Suspense fallback={<Loader />}>
      <Outlet />
    </Suspense>
  );
}

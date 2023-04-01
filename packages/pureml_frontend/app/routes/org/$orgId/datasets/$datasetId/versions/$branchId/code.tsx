import { Outlet } from "@remix-run/react";
import Tabbar from "~/components/Tabbar";

export default function DatasetCode() {
  return (
    <div className="h-full">
      <Outlet />
    </div>
  );
}

import type { ReactNode } from "react";
import { useState } from "react";
import Tabs from "react-simply-tabs";

interface Props {
  tab1: string;
  tab2: string;
  tab3: string;
  tab1Content: ReactNode;
  tab2Content: ReactNode;
  tab3Content: ReactNode;
}

function LandingPgTab(props: Props) {
  const [activeTabIndex, setActiveTabIndex] = useState(0);
  return (
    <Tabs
      activeTabIndex={activeTabIndex}
      onRequestChange={setActiveTabIndex}
      controls={[
        <button
          type="button"
          className="text-base md:text-lg lg:text-2xl text-slate-400"
          key="tab1"
        >
          {props.tab1}
        </button>,
        <button
          type="button"
          className="text-base md:text-lg lg:text-2xl text-slate-400"
          key="tab2"
        >
          {props.tab2}
        </button>,
        <button
          type="button"
          className="text-base md:text-lg lg:text-2xl text-slate-400"
          key="tab3"
        >
          {props.tab3}
        </button>,
      ]}
      controlsWrapperProps={{
        style: {
          paddingBottom: "10px",
          paddingTop: "10px",
          borderTop: "1px solid #CBD5E1",
          borderBottom: "1px solid #CBD5E1",
          display: "flex",
          columnGap: "2rem",
        },
      }}
      activeControlProps={{
        className: "active",
        style: { color: "#020617" },
      }}
    >
      <div className="pt-6">{props.tab1Content}</div>
      <div className="pt-6">{props.tab2Content}</div>
      <div className="pt-6">{props.tab3Content}</div>
    </Tabs>
  );
}

export default LandingPgTab;

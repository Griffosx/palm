import React from "react";
import LeftMenu from "./LeftMenu";
import { darkBrown } from "../styles/themes";

interface WrapperProps {
  children: React.ReactNode;
  title?: string;
}

const Wrapper: React.FC<WrapperProps> = ({ children, title }) => {
  return (
    <div
      className="flex h-full text-darkBrown"
      style={{
        color: darkBrown,
      }}
    >
      {/* Left panel with menu */}
      <LeftMenu title={title} />

      {/* Right content area with left padding to accommodate fixed menu */}
      <div className="pl-64 w-full py-6 px-6">{children}</div>
    </div>
  );
};

export default Wrapper;

import React from "react";
import ContentBox from "./ContentBox";
import LeftMenu from "./LeftMenu";

interface WrapperProps {
  children: React.ReactNode;
  title?: string;
}

const Wrapper: React.FC<WrapperProps> = ({ children, title }) => {
  return (
    <div className="flex h-full">
      {/* Left panel with menu */}
      <LeftMenu title={title} />

      {/* Right content area */}
      <ContentBox>{children}</ContentBox>
    </div>
  );
};

export default Wrapper;

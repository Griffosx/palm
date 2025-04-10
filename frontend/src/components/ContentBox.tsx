import React from "react";
import BackgroundOverlay from "./BackgroundOverlay";
import { hexToRgba, brown, cream, borderCream } from "../styles/themes";

interface ContentBoxProps {
  children: React.ReactNode;
  className?: string;
  backgroundColor?: string;
  boxShadow?: string;
}

const ContentBox: React.FC<ContentBoxProps> = ({
  children,
  className = "",
  backgroundColor = cream,
  boxShadow = `${hexToRgba(borderCream)} 0px 0px 3px`,
}) => {
  return (
    <div
      className={`flex flex-col relative h-full rounded-xl ${className}`}
      style={{
        boxShadow,
      }}
    >
      {/* Background with customizable color */}
      <BackgroundOverlay
        solidColor={backgroundColor}
        rounded="rounded-xl"
        zIndex={0}
      />

      {/* Main content */}
      <div
        className={`relative z-10 flex-1 px-6 pb-6 pt-0 rounded-xl overflow-auto border-1`}
        style={{
          borderColor: borderCream,
        }}
      >
        <div className="rounded-xl p-6 md:p-8 transition-all duration-300 min-h-[200px]">
          <div className="text-gray-700">{children}</div>
        </div>
      </div>
    </div>
  );
};

export default ContentBox;

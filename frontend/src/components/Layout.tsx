import React from "react";
import BackgroundOverlay from "./BackgroundOverlay";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen relative pt-15 pb-15 pr-8">
      {/* Background with gradient */}
      <BackgroundOverlay gradient="main" zIndex={0} position="fixed" />

      {/* Content container */}
      <div className="relative z-10 h-full w-full">{children}</div>
    </div>
  );
};

export default Layout;

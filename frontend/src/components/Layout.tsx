import React from "react";
import BackgroundOverlay from "./BackgroundOverlay";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="h-screen overflow-hidden relative pt-15 pb-5 pr-8">
      {/* Background with gradient */}
      <BackgroundOverlay gradient="main" zIndex={0} position="fixed" />

      {/* Content container needs h-full to pass height down */}
      <div className="relative z-10 h-full w-full">{children}</div>
    </div>
  );
};

export default Layout;

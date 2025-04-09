import React from "react";
import BackgroundOverlay from "./BackgroundOverlay";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen relative pt-15">
      {/* Background with solid color and noise */}
      <BackgroundOverlay gradient="main" />

      {/* Content container */}
      <div className="relative z-10 min-h-screen w-full">{children}</div>
    </div>
  );
};

export default Layout;

import React from "react";
import { Link, useLocation } from "react-router-dom";
import BackgroundOverlay from "./BackgroundOverlay";

interface WrapperProps {
  children: React.ReactNode;
  title?: string;
}

const Wrapper: React.FC<WrapperProps> = ({ children, title }) => {
  const location = useLocation();

  const menuItems = [
    { name: "Inbox", path: "/inbox" },
    { name: "Sent", path: "/sent" },
    { name: "Draft", path: "/draft" },
  ];

  return (
    <div className="flex h-screen">
      {/* Left panel with menu */}
      <div className="w-64 p-6 flex-shrink-0">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Palm</h2>
        <nav>
          <ul className="space-y-2">
            {menuItems.map((item) => (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={`block px-4 py-2 rounded-lg transition-colors ${
                    location.pathname === item.path
                      ? "bg-orange-200/50 text-orange-800 font-medium"
                      : "text-gray-700 hover:bg-orange-100/30"
                  }`}
                >
                  {item.name}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </div>

      {/* Right content area */}
      <div className="flex-1 flex flex-col relative">
        {/* Background with cream color */}
        <BackgroundOverlay solidColor="cream" className="rounded-l-xl" />

        {/* Top bar with same style as content panel */}
        <div className="relative z-10 m-6 mb-0 backdrop-blur-sm rounded-t-xl p-4 flex items-center">
          {title && (
            <h1 className="text-2xl font-bold text-gray-800">{title}</h1>
          )}
        </div>

        {/* Main content */}
        <div className="relative z-10 flex-1 px-6 pb-6 pt-0 overflow-auto">
          <div className="backdrop-blur-sm rounded-b-xl p-6 md:p-8 transition-all duration-300 min-h-[200px]">
            <div className="text-gray-700">{children}</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Wrapper;

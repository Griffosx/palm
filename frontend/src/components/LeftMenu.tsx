import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { darkBrown, peach, orange } from "../styles/themes";
import {
  DraftIcon,
  InboxIcon,
  KnowledgeIcon,
  SentIcon,
  TrashIcon,
} from "./Icons";

interface LeftMenuProps {
  title?: string;
}

const LeftMenu: React.FC<LeftMenuProps> = () => {
  const location = useLocation();
  const [hoveredPath, setHoveredPath] = useState<string | null>(null);

  const menuItems = [
    { name: "Knowledge", path: "/knowledge", icon: KnowledgeIcon },
    { name: "Inbox", path: "/inbox", icon: InboxIcon },
    { name: "Sent", path: "/sent", icon: SentIcon },
    { name: "Draft", path: "/draft", icon: DraftIcon },
    { name: "Trash", path: "/trash", icon: TrashIcon },
  ];

  return (
    <div className="flex flex-col w-64 p-6">
      <nav>
        <ul>
          {menuItems.map((item) => {
            const isActive = location.pathname === item.path;
            const isHovered = hoveredPath === item.path;
            const useActiveColor = isActive || isHovered;
            const IconComponent = item.icon;

            return (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className="flex items-center my-1 px-4 py-2 rounded-4xl transition-colors text-xl"
                  style={{
                    backgroundColor: useActiveColor ? peach : "transparent",
                    color: useActiveColor ? orange : darkBrown,
                    fontWeight: isActive ? "500" : "normal",
                  }}
                  onMouseEnter={(e: React.MouseEvent<HTMLAnchorElement>) => {
                    setHoveredPath(item.path);
                  }}
                  onMouseLeave={(e: React.MouseEvent<HTMLAnchorElement>) => {
                    setHoveredPath(null);
                  }}
                >
                  <div className="w-9 h-9 mr-4 flex items-center justify-center">
                    <IconComponent fill={useActiveColor ? orange : darkBrown} />
                  </div>
                  {item.name}
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>
    </div>
  );
};

export default LeftMenu;

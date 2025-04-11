import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { darkBrown, peach, orange } from "../styles/themes";
import {
  DraftIcon,
  InboxIcon,
  KnowledgeIcon,
  SentIcon,
  TrashIcon,
  SettingsIcon,
} from "./Icons";
import NoiseOverlay from "./NoiseOverlay";

interface LeftMenuProps {
  title?: string;
}

// Define props for the MenuItem component
interface MenuItemProps {
  name: string;
  path: string;
  icon: React.FC<{ fill: string; width?: string; height?: string }>;
  isActive: boolean;
  isHovered: boolean;
  onMouseEnter: (e: React.MouseEvent<HTMLAnchorElement>) => void;
  onMouseLeave: (e: React.MouseEvent<HTMLAnchorElement>) => void;
}

// MenuItem component
const MenuItem: React.FC<MenuItemProps> = ({
  name,
  path,
  icon: IconComponent,
  isActive,
  isHovered,
  onMouseEnter,
  onMouseLeave,
}) => {
  const useActiveColor = isActive || isHovered;

  return (
    <div key={path}>
      <Link
        to={path}
        className="flex items-center mb-1 px-4 py-2 rounded-3xl transition-colors text-xl relative overflow-hidden"
        style={{
          backgroundColor: useActiveColor ? peach : "transparent",
          color: useActiveColor ? orange : darkBrown,
          fontWeight: isActive ? "500" : "normal",
        }}
        onMouseEnter={onMouseEnter}
        onMouseLeave={onMouseLeave}
      >
        {useActiveColor && <NoiseOverlay />}
        <div className="w-9 h-9 mr-4 flex items-center justify-center relative z-10">
          <IconComponent
            fill={useActiveColor ? orange : darkBrown}
            width="28"
            height="28"
          />
        </div>
        <span className="relative z-10">{name}</span>
      </Link>
    </div>
  );
};

const LeftMenu: React.FC<LeftMenuProps> = () => {
  const location = useLocation();
  const [hoveredPath, setHoveredPath] = useState<string | null>(null);

  // Define individual menu items
  const menuItemsData = [
    { name: "Knowledge", path: "/knowledge", icon: KnowledgeIcon },
    { name: "Inbox", path: "/inbox", icon: InboxIcon },
    { name: "Sent", path: "/sent", icon: SentIcon },
    { name: "Draft", path: "/draft", icon: DraftIcon },
    { name: "Trash", path: "/trash", icon: TrashIcon },
  ];

  const settingsPath = "/settings";

  return (
    <div className="flex flex-col w-64 p-6 pt-3 fixed top-15 h-[calc(100vh-80px)] justify-between">
      <div>
        {menuItemsData.map((item) => (
          <MenuItem
            key={item.path}
            name={item.name}
            path={item.path}
            icon={item.icon}
            isActive={location.pathname === item.path}
            isHovered={hoveredPath === item.path}
            onMouseEnter={() => setHoveredPath(item.path)}
            onMouseLeave={() => setHoveredPath(null)}
          />
        ))}
      </div>

      <div className="mb-4">
        <MenuItem
          key={settingsPath}
          name="Settings"
          path={settingsPath}
          icon={SettingsIcon}
          isActive={location.pathname === settingsPath}
          isHovered={hoveredPath === settingsPath}
          onMouseEnter={() => setHoveredPath(settingsPath)}
          onMouseLeave={() => setHoveredPath(null)}
        />
      </div>
    </div>
  );
};

export default LeftMenu;

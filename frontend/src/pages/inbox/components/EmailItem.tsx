import React, { useMemo } from "react";
import { controllers } from "../../../../wailsjs/go/models";
import NoiseOverlay from "../../../components/NoiseOverlay";
import {
  darkAquamarine,
  orange,
  peach,
  aquamarine,
  lightAquamarine,
} from "../../../styles/themes";

interface EmailItemProps {
  email: controllers.EmailResponse;
  onClick: (id: number) => void;
  isSelected: boolean;
}

const EmailItem: React.FC<EmailItemProps> = ({
  email,
  onClick,
  isSelected,
}) => {
  const formattedDate = new Date(email.receivedAt).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });

  const senderDisplay = email.senderName || email.senderEmail;
  const firstLetter = senderDisplay.charAt(0).toUpperCase();

  // Randomly select a background color
  const backgroundColor = useMemo(() => {
    const colors = [darkAquamarine, orange, peach];
    return colors[Math.floor(Math.random() * colors.length)];
  }, [email.id]); // Ensure same color for same email

  // Track hover state
  const [isHovered, setIsHovered] = React.useState(false);

  const leftSide = () => {
    return (
      <div className="mr-3">
        <div
          className="w-10 h-10 rounded-full flex items-center justify-center font-medium relative overflow-hidden"
          style={{ backgroundColor }}
        >
          <NoiseOverlay />
          <span className="z-10">{firstLetter}</span>
        </div>
        {!email.isRead && (
          <div className="flex items-center justify-center h-7">
            <div
              className={`w-3 h-3 rounded-full absolute`}
              style={{ backgroundColor: aquamarine }}
            >
              <NoiseOverlay />
            </div>
          </div>
        )}
      </div>
    );
  };

  const rightSide = () => {
    return (
      <div className="flex-grow flex items-start relative z-10 min-w-0">
        <div className="flex-grow min-w-0">
          <div className="flex justify-between items-start mb-1 text-left">
            <div className="font-medium truncate min-w-0">{senderDisplay}</div>
            <div className="text-xs ml-2 whitespace-nowrap flex-shrink-0">
              {formattedDate}
            </div>
          </div>
          <div className="font-medium mb-1 truncate text-left min-w-0">
            {email.subject}
          </div>
          <div className="text-sm truncate text-left min-w-0">
            {email.body.replace(/<[^>]*>/g, "").substring(0, 100)}
            {email.body.length > 100 ? "..." : ""}
          </div>
          {email.attachments && email.attachments.length > 0 && (
            <div className="text-left text-xs text-gray-500 mt-1">
              <span className="inline-flex items-left">
                <svg
                  className="w-4 h-4 mr-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.414a4 4 0 00-5.656-5.656l-6.415 6.415a6 6 0 108.486 8.486L20.5 13"
                  />
                </svg>
                {email.attachments.length}{" "}
                {email.attachments.length === 1 ? "attachment" : "attachments"}
              </span>
            </div>
          )}
        </div>
      </div>
    );
  };

  return (
    <div
      className="flex p-5 cursor-pointer transition rounded-xl relative mb-2"
      style={{
        backgroundColor:
          isHovered || isSelected ? lightAquamarine : "transparent",
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onClick={() => onClick(email.id)}
    >
      {(isHovered || isSelected) && (
        <NoiseOverlay className="rounded-xl absolute inset-0" />
      )}
      {leftSide()}
      {rightSide()}
    </div>
  );
};

export default EmailItem;

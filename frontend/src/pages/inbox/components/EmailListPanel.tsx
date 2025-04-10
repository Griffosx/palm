import React from "react";
import SearchInput from "./SearchInput";
import EmailList from "./EmailList";
import BackgroundOverlay from "../../../components/BackgroundOverlay";
import { hexToRgba, cream, borderCream } from "../../../styles/themes";

interface EmailListPanelProps {
  searchQuery: string;
  onSearch: (query: string) => void;
  onSelectEmail: (id: number) => void;
  selectedEmailId: number | null;
}

const EmailListPanel: React.FC<EmailListPanelProps> = ({
  searchQuery,
  onSearch,
  onSelectEmail,
  selectedEmailId,
}) => {
  const backgroundColor = cream;
  const boxShadow = `${hexToRgba(borderCream)} 0px 0px 3px`;

  return (
    // Fixed width, non-shrinking, full height container
    <div className="w-96 flex-shrink-0 h-full flex flex-col relative rounded-xl">
      {/* Background Overlay */}
      <BackgroundOverlay
        solidColor={backgroundColor}
        rounded="rounded-xl"
        zIndex={0}
      />

      {/* Inner container for border, padding, and flex layout */}
      <div
        className="relative z-10 flex-1 flex flex-col rounded-xl overflow-hidden border-1"
        style={{
          borderColor: borderCream,
          boxShadow: boxShadow, // Apply shadow here as well
        }}
      >
        {/* Padding Container */}
        <div className="flex flex-col flex-1 p-6 md:p-8 overflow-hidden">
          {/* Search Input (non-scrollable part) */}
          <div className="mb-4 flex-shrink-0">
            <SearchInput onSearch={onSearch} />
          </div>
          {/* Scrollable Email List */}
          <div className="flex-grow overflow-y-auto">
            <EmailList
              searchQuery={searchQuery}
              onSelectEmail={onSelectEmail}
              selectedEmailId={selectedEmailId}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default EmailListPanel;

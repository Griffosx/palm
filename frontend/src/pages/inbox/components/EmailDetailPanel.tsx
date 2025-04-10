import React from "react";
import EmailDetail from "./EmailDetail";
import BackgroundOverlay from "../../../components/BackgroundOverlay";
import { hexToRgba, cream, borderCream } from "../../../styles/themes";

interface EmailDetailPanelProps {
  selectedEmailId: number | null;
}

const EmailDetailPanel: React.FC<EmailDetailPanelProps> = ({
  selectedEmailId,
}) => {
  const backgroundColor = cream;
  const boxShadow = `${hexToRgba(borderCream)} 0px 0px 3px`;

  return (
    // Grow to fill space, full height container
    <div className="flex-grow h-full flex flex-col relative rounded-xl">
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
          boxShadow: boxShadow,
        }}
      >
        {/* Padding Container - Crucially, THIS handles the final padding */}
        {/* EmailDetail itself should NOT have padding if this container does */}
        {/* EmailDetail also needs h-full to fill this padded container */}
        <div className="flex-1 p-6 md:p-8 overflow-hidden">
          <EmailDetail emailId={selectedEmailId} />
        </div>
      </div>
    </div>
  );
};

export default EmailDetailPanel;

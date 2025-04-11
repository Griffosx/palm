import React, { useState, useEffect } from "react";
import BackgroundOverlay from "../../../components/BackgroundOverlay";
import { hexToRgba, cream, borderCream } from "../../../styles/themes";
import { controllers } from "../../../../wailsjs/go/models";
import { GetEmail } from "../../../../wailsjs/go/main/App";
import {
  ReplyIcon,
  ReplyAllIcon,
  ForwardIcon,
  AIIcon,
  KnowledgeIcon,
} from "../../../components/Icons";
import NoiseOverlay from "../../../components/NoiseOverlay";
import { getDeterministicColor } from "../utils/colorUtils";

interface EmailDetailPanelProps {
  selectedEmailId: number | null;
}

const EmailDetailPanel: React.FC<EmailDetailPanelProps> = ({
  selectedEmailId,
}) => {
  const [email, setEmail] = useState<controllers.EmailResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const backgroundColor = cream;
  const boxShadow = `${hexToRgba(borderCream)} 0px 0px 3px`;

  useEffect(() => {
    const fetchEmailDetail = async () => {
      if (!selectedEmailId) {
        setEmail(null);
        setError(null);
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);
      setEmail(null);

      try {
        const emailData = await GetEmail(selectedEmailId);
        setEmail(emailData);
      } catch (err) {
        console.error("Error fetching email details:", err);
        setError("Failed to load email. Please try again.");
      } finally {
        setLoading(false);
      }
    };

    fetchEmailDetail();
  }, [selectedEmailId]);

  const formatDate = (dateString: string | undefined): string => {
    if (!dateString) return "";
    return new Date(dateString).toLocaleString(undefined, {
      weekday: "short",
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "2-digit",
    });
  };

  const getInitials = (emailData: controllers.EmailResponse | null): string => {
    if (!emailData) return "?";
    const name = emailData.senderName;
    const emailAddr = emailData.senderEmail;
    if (name) {
      return name.charAt(0).toUpperCase();
    }
    if (emailAddr) {
      return emailAddr.charAt(0).toUpperCase();
    }
    return "?";
  };

  const formatRecipients = (
    emailData: controllers.EmailResponse | null
  ): string => {
    if (
      !emailData ||
      !emailData.recipients ||
      emailData.recipients.length === 0
    ) {
      return "No recipients";
    }
    const displayRecipients = emailData.recipients
      .slice(0, 3)
      .map((r) => r.name || r.email)
      .join(", ");

    return `To: ${displayRecipients}${
      emailData.recipients.length > 3 ? "..." : ""
    }`;
  };

  return (
    <div className="flex-grow h-full flex flex-col relative rounded-xl min-w-0">
      <BackgroundOverlay
        solidColor={backgroundColor}
        rounded="rounded-xl"
        zIndex={0}
      />

      <div
        className="relative z-10 flex-1 flex flex-col rounded-xl overflow-hidden border-1"
        style={{
          borderColor: borderCream,
          boxShadow: boxShadow,
        }}
      >
        {!selectedEmailId ? (
          <div className="h-full flex items-center justify-center p-6 md:p-8">
            <p>Select an email to view its content</p>
          </div>
        ) : loading ? (
          <div className="h-full flex items-center justify-center p-6 md:p-8">
            <div className="inline-block animate-spin h-8 w-8 border-2 border-gray-500 border-t-transparent rounded-full"></div>
          </div>
        ) : error ? (
          <div className="h-full flex items-center justify-center text-red-500 p-6 md:p-8">
            <p>{error}</p>
          </div>
        ) : email ? (
          <div className="flex flex-col h-full">
            <div
              className="flex-shrink-0 border-b"
              style={{ borderColor: borderCream }}
            >
              <div className="p-4 md:p-6 flex justify-between items-start text-left">
                <div className="flex-grow flex-shrink min-w-0 mr-4">
                  <h1
                    className="text-l md:text-l font-semibold mb-2 truncate"
                    title={email.subject || "(no subject)"}
                  >
                    {email.subject || "(no subject)"}
                  </h1>
                  <div className="flex items-center mb-1">
                    <div
                      className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center mr-3 relative overflow-hidden flex-shrink-0"
                      style={{
                        backgroundColor: getDeterministicColor(
                          email.senderEmail
                        ),
                      }}
                    >
                      <NoiseOverlay />
                      <span className="z-10 text-sm">{getInitials(email)}</span>
                    </div>
                    <div className="text-sm overflow-hidden">
                      <div
                        className="font-medium truncate"
                        title={email.senderName || email.senderEmail}
                      >
                        {email.senderName || email.senderEmail}
                      </div>
                      <div
                        className="text-xstruncate"
                        title={email.senderEmail}
                      >
                        &lt;{email.senderEmail}&gt;
                      </div>
                    </div>
                  </div>
                  <div
                    className="text-xs mb-2 truncate"
                    title={formatRecipients(email)}
                  >
                    {formatRecipients(email)}
                  </div>
                </div>

                <div className="flex flex-col items-end flex-shrink-0 ml-4">
                  <div className="flex space-x-2 mb-2">
                    <button className="p-1.5 hover:text-orange-500 transition-colors duration-150 cursor-pointer">
                      <ReplyIcon fill="currentColor" width="20" height="20" />
                    </button>
                    <button className="p-1.5 hover:text-orange-500 transition-colors duration-150 cursor-pointer">
                      <ReplyAllIcon
                        fill="currentColor"
                        width="20"
                        height="20"
                      />
                    </button>
                    <button className="p-1.5 hover:text-orange-500 transition-colors duration-150 cursor-pointer">
                      <ForwardIcon fill="currentColor" width="20" height="20" />
                    </button>
                    <button className="p-1.5 hover:text-orange-500 transition-colors duration-150 cursor-pointer">
                      <AIIcon fill="currentColor" width="20" height="20" />
                    </button>
                  </div>
                  <div className="text-xs whitespace-nowrap">
                    {formatDate(email.receivedAt)}
                  </div>
                </div>
              </div>
            </div>

            <div className="flex-grow overflow-y-auto p-4 md:p-6">
              <div
                className="overflow-x-auto text-left"
                dangerouslySetInnerHTML={{
                  __html: email.body || "<p>No content</p>",
                }}
              />

              {email.attachments && email.attachments.length > 0 && (
                <div className="mt-6 pt-4 border-t border-gray-200">
                  <h2 className="text-base font-medium mb-2">Attachments</h2>
                  <div className="flex flex-wrap gap-2">
                    {email.attachments.map(
                      (attachment: { filename: string }, index: number) => (
                        <div
                          key={index}
                          className="px-3 py-1.5 bg-gray-100 hover:bg-gray-200 text-sm rounded-lg flex items-center cursor-pointer"
                        >
                          <svg
                            className="w-4 h-4 mr-2 flex-shrink-0"
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
                          <span
                            className="truncate"
                            title={attachment.filename}
                          >
                            {attachment.filename}
                          </span>
                        </div>
                      )
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        ) : (
          <div className="h-full flex items-center justify-center p-6 md:p-8">
            <p>Email data not available.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default EmailDetailPanel;

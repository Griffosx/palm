import React, { useEffect, useState } from "react";
import { controllers } from "../../../../wailsjs/go/models";
import { GetEmail } from "../../../../wailsjs/go/main/App";
import NoiseOverlay from "../../../components/NoiseOverlay";
import { aquamarine } from "../../../styles/themes";

interface EmailDetailProps {
  emailId: number | null;
}

const EmailDetail: React.FC<EmailDetailProps> = ({ emailId }) => {
  const [email, setEmail] = useState<controllers.EmailResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchEmailDetail = async () => {
      if (!emailId) {
        setEmail(null);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const emailData = await GetEmail(emailId);
        setEmail(emailData);
      } catch (err) {
        console.error("Error fetching email details:", err);
        setError("Failed to load email. Please try again.");
      } finally {
        setLoading(false);
      }
    };

    fetchEmailDetail();
  }, [emailId]);

  if (!emailId) {
    return (
      <div className="h-full flex items-center justify-center text-gray-500">
        <p>Select an email to view its content</p>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="inline-block animate-spin h-8 w-8 border-2 border-gray-500 border-t-transparent rounded-full"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-full flex items-center justify-center text-red-500">
        <p>{error}</p>
      </div>
    );
  }

  if (!email) {
    return null;
  }

  const formattedDate = new Date(email.receivedAt).toLocaleString(undefined, {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });

  return (
    <div className="flex flex-col h-full">
      <div className="mb-4 flex-shrink-0">
        <h1 className="text-2xl font-bold mb-4">{email.subject}</h1>

        <div className="flex items-center">
          <div className="w-10 h-10 rounded-full bg-aquamarine flex items-center justify-center mr-4 relative overflow-hidden">
            <NoiseOverlay />
            <span className="z-10 font-bold">
              {(email.senderName || email.senderEmail).charAt(0).toUpperCase()}
            </span>
          </div>

          <div>
            <div className="font-medium">
              {email.senderName || email.senderEmail}
            </div>
            <div className="text-sm text-gray-600">{formattedDate}</div>
          </div>
        </div>
      </div>

      <div className="flex-grow overflow-y-auto pr-2">
        <div dangerouslySetInnerHTML={{ __html: email.body }} />

        {email.attachments && email.attachments.length > 0 && (
          <div className="mt-6 pt-4 border-t border-gray-200">
            <h2 className="text-lg font-medium mb-2">Attachments</h2>
            <div className="flex flex-wrap gap-2">
              {email.attachments.map(
                (attachment: { filename: string }, index: number) => (
                  <div
                    key={index}
                    className="px-3 py-2 bg-gray-100 rounded-lg flex items-center"
                  >
                    <svg
                      className="w-5 h-5 mr-2"
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
                    <span>{attachment.filename}</span>
                  </div>
                )
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default EmailDetail;

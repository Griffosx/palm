import React, { useState, useEffect, useRef, useCallback } from "react";
import { controllers } from "../../../../wailsjs/go/models";
import { ListEmails } from "../../../../wailsjs/go/main/App";
import EmailItem from "./EmailItem";

const PAGE_SIZE = 20;
const ACCOUNT_ID = 1;

interface EmailListProps {
  searchQuery?: string;
  onSelectEmail: (id: number) => void;
  selectedEmailId: number | null;
}

const EmailList: React.FC<EmailListProps> = ({
  searchQuery,
  onSelectEmail,
  selectedEmailId,
}) => {
  const [emails, setEmails] = useState<controllers.EmailResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [page, setPage] = useState(1);
  const [error, setError] = useState<string | null>(null);

  const observer = useRef<IntersectionObserver | null>(null);
  const lastEmailElementRef = useCallback(
    (node: HTMLDivElement | null) => {
      if (loading) return;
      if (observer.current) observer.current.disconnect();

      observer.current = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && hasMore) {
          setPage((prevPage) => prevPage + 1);
        }
      });

      if (node) observer.current.observe(node);
    },
    [loading, hasMore]
  );

  // Initial load and when search changes
  useEffect(() => {
    setEmails([]);
    setPage(1);
    setHasMore(true);
  }, [searchQuery]);

  // Load emails when page changes
  useEffect(() => {
    const fetchEmails = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await ListEmails(ACCOUNT_ID, page, PAGE_SIZE);

        setEmails((prevEmails) => {
          // Add only unique emails
          const newEmails = response.emails.filter(
            (newEmail) =>
              !prevEmails.some(
                (existingEmail) => existingEmail.id === newEmail.id
              )
          );
          return [...prevEmails, ...newEmails];
        });

        setHasMore(page < response.totalPages);
      } catch (err) {
        setError("Failed to load emails. Please try again.");
        console.error("Error fetching emails:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchEmails();
  }, [page, searchQuery]);

  if (error) {
    return <div className="text-red-500 p-4">{error}</div>;
  }

  return (
    <div>
      {emails.length === 0 && !loading ? (
        <div className="p-4 text-center">No emails found</div>
      ) : (
        <>
          {emails.map((email, index) => {
            const isSelected = email.id === selectedEmailId;
            if (emails.length === index + 1) {
              return (
                <div ref={lastEmailElementRef} key={email.id}>
                  <EmailItem
                    email={email}
                    onClick={onSelectEmail}
                    isSelected={isSelected}
                  />
                </div>
              );
            } else {
              return (
                <EmailItem
                  key={email.id}
                  email={email}
                  onClick={onSelectEmail}
                  isSelected={isSelected}
                />
              );
            }
          })}
          {loading && (
            <div className="p-4 text-center">
              <div className="inline-block animate-spin h-6 w-6 border-2 border-gray-500 border-t-transparent rounded-full"></div>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default EmailList;

import React, { useState } from "react";
import Layout from "../../components/Layout";
import Wrapper from "../../components/Wrapper";
import EmailListPanel from "./components/EmailListPanel";
import EmailDetailPanel from "./components/EmailDetailPanel";

const InboxPage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedEmailId, setSelectedEmailId] = useState<number | null>(null);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    // Optional: Deselect email when searching?
    // setSelectedEmailId(null);
  };

  const handleSelectEmail = (id: number) => {
    setSelectedEmailId(id);
  };

  return (
    <Layout>
      <Wrapper title="Inbox">
        {/* Height is handled by Layout and Wrapper */}
        <div className="h-full w-full flex gap-6">
          <EmailListPanel
            searchQuery={searchQuery}
            onSearch={handleSearch}
            onSelectEmail={handleSelectEmail}
            selectedEmailId={selectedEmailId}
          />
          <EmailDetailPanel selectedEmailId={selectedEmailId} />
        </div>
      </Wrapper>
    </Layout>
  );
};

export default InboxPage;

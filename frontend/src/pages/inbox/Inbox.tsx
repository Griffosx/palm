import React, { useState } from "react";
import Layout from "../../components/Layout";
import Wrapper from "../../components/Wrapper";
import ContentBox from "../../components/ContentBox";
import SearchInput from "./components/SearchInput";
import EmailList from "./components/EmailList";
import EmailDetail from "./components/EmailDetail";

const InboxPage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedEmailId, setSelectedEmailId] = useState<number | null>(null);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const handleSelectEmail = (id: number) => {
    setSelectedEmailId(id);
  };

  return (
    <Layout>
      <Wrapper title="Inbox">
        <div className="h-full w-full flex gap-6">
          {/* Left panel - Email list */}
          <div className="w-2/5 h-full">
            <ContentBox className="h-full flex flex-col">
              <div className="mb-4">
                <SearchInput onSearch={handleSearch} />
              </div>
              <div className="flex-grow overflow-hidden">
                <EmailList
                  searchQuery={searchQuery}
                  onSelectEmail={handleSelectEmail}
                />
              </div>
            </ContentBox>
          </div>

          {/* Right panel - Email detail */}
          <div className="w-3/5 h-full">
            <ContentBox className="h-full overflow-hidden">
              <EmailDetail emailId={selectedEmailId} />
            </ContentBox>
          </div>
        </div>
      </Wrapper>
    </Layout>
  );
};

export default InboxPage;

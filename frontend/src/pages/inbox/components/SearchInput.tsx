import React, { useState } from "react";
import { beige } from "../../../styles/themes";
import NoiseOverlay from "../../../components/NoiseOverlay";

interface SearchInputProps {
  onSearch: (query: string) => void;
}

const SearchInput: React.FC<SearchInputProps> = ({ onSearch }) => {
  const [query, setQuery] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(query);
  };

  return (
    <form onSubmit={handleSearch} className="mb-4">
      <div className="relative">
        <div
          className="absolute inset-0 rounded-full overflow-hidden"
          style={{ backgroundColor: beige }}
        >
          <NoiseOverlay />
        </div>
        <input
          type="text"
          className="w-full p-2 pl-10 rounded-full focus:outline-none focus:ring-2 focus:ring-[#BCA482] relative z-10 bg-transparent"
          placeholder="Search emails..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none z-20">
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
      </div>
    </form>
  );
};

export default SearchInput;

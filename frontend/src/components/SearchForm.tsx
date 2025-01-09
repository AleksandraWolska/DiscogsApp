import React, { useState } from "react";

interface SearchFormProps {
  onSearch: (query: string) => void;
}

const SearchForm: React.FC<SearchFormProps> = ({ onSearch }) => {
  const [artist, setArtist] = useState("");
  const [format, setFormat] = useState("");
  const [year, setYear] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const query = new URLSearchParams();
    if (artist) query.append("artist", artist);
    if (format) query.append("format", format);
    if (year) query.append("year", year);
    onSearch(query.toString());
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>
          Artist:
          <input type="text" value={artist} onChange={(e) => setArtist(e.target.value)} />
        </label>
      </div>
      <div>
        <label>
          Format:
          <input type="text" value={format} onChange={(e) => setFormat(e.target.value)} />
        </label>
      </div>
      <div>
        <label>
          Year:
          <input type="text" value={year} onChange={(e) => setYear(e.target.value)} />
        </label>
      </div>
      <button type="submit">Search</button>
    </form>
  );
};

export default SearchForm;
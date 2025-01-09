import React, { useState, useEffect } from "react";
import "./App.css";
import ReleaseList from "./components/ReleaseList";
import SearchForm from "./components/SearchForm";
import { fetchAllReleases, searchReleases } from "./services/api";

function App() {
  const [releases, setReleases] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const loadReleases = async () => {
      try {
        const data = await fetchAllReleases();
        setReleases(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    loadReleases();
  }, []);

  const handleSearch = async (query: string) => {
    setLoading(true);
    try {
      const data = await searchReleases(query);
      setReleases(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Discogs Releases</h1>
        <SearchForm onSearch={handleSearch} />
        {loading && <p>Loading...</p>}
        {error && <p>Error: {error}</p>}
        <ReleaseList releases={releases} />
      </header>
    </div>
  );
}

export default App;
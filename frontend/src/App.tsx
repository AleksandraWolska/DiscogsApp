import React, { useState, useEffect } from "react";
import { Container, Typography, CircularProgress, Box } from "@mui/material";
import "./App.css";
import ReleaseList from "./components/ReleaseList";
import FilterOnFrontend from "./components/FilterOnFrontend";
import { fetchAllReleases, fetchArtists, fetchFormats } from "./services/api";

function App() {
  const [releases, setReleases] = useState([]);
  const [filteredReleases, setFilteredReleases] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [artists, setArtists] = useState([]);
  const [formats, setFormats] = useState([]);

  const loadReleases = async () => {
    setLoading(true);
    setError("");
    try {
      const data = await fetchAllReleases();
      setReleases(data);
      setFilteredReleases(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const loadArtistsAndFormats = async () => {
    try {
      const [artistsData, formatsData] = await Promise.all([fetchArtists(), fetchFormats()]);
      setArtists(artistsData);
      setFormats(formatsData);
    } catch (err) {
      console.error("Failed to fetch artists or formats", err);
    }
  };

  useEffect(() => {
    const initialize = async () => {
      await loadReleases();
      await loadArtistsAndFormats();
    };
    initialize();
  }, []);

  const handleFrontendFilter = (filteredReleases: any) => {
    console.log("Filtered Releases");
    console.log(filteredReleases);
    setFilteredReleases(filteredReleases);
  };

  return (
    <Container>
      <Box my={4}>
        
        <Typography variant="h4" component="h1" gutterBottom>
          Discogs Releases
        </Typography>
        <Typography variant="h5" component="h2">
          Filter (frontend ones)
        </Typography>
        <FilterOnFrontend releases={releases} artists={artists} formats={formats} onFilter={handleFrontendFilter} />
        {loading && <CircularProgress />}
        {error && <Typography color="error">{error}</Typography>}
        <ReleaseList releases={filteredReleases} />
        
      </Box>
    </Container>
  );
}

export default App;
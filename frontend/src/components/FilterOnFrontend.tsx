import React, { useState } from "react";
import { TextField, Button, MenuItem, Box } from "@mui/material";

interface Artist {
  id: number;
  name: string;
}

interface Format {
  id: number;
  name: string;
}

interface Release {
  id: number;
  title: string;
  year: number;
  artists: Artist[];
  catalogNo: string;
  thumb: string;
  resourceURL: string;
  formats: Format[];
  status: string;
}

interface FilterOnFrontendProps {
  releases: Release[];
  artists: Artist[];
  formats: Format[];
  onFilter: (filteredReleases: Release[]) => void;
}

const FilterOnFrontend: React.FC<FilterOnFrontendProps> = ({ releases, artists, formats, onFilter }) => {
  const [artist, setArtist] = useState("");
  const [format, setFormat] = useState("");
  const [year, setYear] = useState("");
  const [title, setTitle] = useState("");
  const [status, setStatus] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    let filteredReleases = releases;

    if (artist) {
      filteredReleases = filteredReleases.filter((release) =>
        release.artists.some((a) => a.name.toLowerCase().includes(artist.toLowerCase()))
      );
    }

    if (format) {
      filteredReleases = filteredReleases.filter((release) =>
        release.formats.some((f) => f.name.toLowerCase().includes(format.toLowerCase()))
      );
    }

    if (year) {
      filteredReleases = filteredReleases.filter((release) => release.year.toString() === year);
    }

    if (title) {
      filteredReleases = filteredReleases.filter((release) =>
        release.title.toLowerCase().includes(title.toLowerCase())
      );
    }

    if (status) {
      filteredReleases = filteredReleases.filter((release) =>
        release.status.toLowerCase().includes(status.toLowerCase())
      );
    }

    onFilter(filteredReleases);
  };

  const handleReset = () => {
    setArtist("");
    setFormat("");
    setYear("");
    setTitle("");
    setStatus("");
    onFilter(releases);
  };

  return (
    <Box component="form" onSubmit={handleSearch} sx={{ mt: 3 }}>
      <TextField
        select
        label="Artist"
        value={artist}
        onChange={(e) => setArtist(e.target.value)}
        fullWidth
        sx={{ mb: 2 }}
      >
        <MenuItem value="">
          <em>Select an artist</em>
        </MenuItem>
        {artists.map((artist) => (
          <MenuItem key={artist.id} value={artist.name}>
            {artist.name}
          </MenuItem>
        ))}
      </TextField>

      <TextField
        select
        label="Format"
        value={format}
        onChange={(e) => setFormat(e.target.value)}
        fullWidth
        sx={{ mb: 2 }}
      >
        <MenuItem value="">
          <em>Select a format</em>
        </MenuItem>
        {formats.map((format) => (
          <MenuItem key={format.id} value={format.name}>
            {format.name}
          </MenuItem>
        ))}
      </TextField>

      <TextField
        label="Year"
        value={year}
        onChange={(e) => setYear(e.target.value)}
        fullWidth
        sx={{ mb: 2 }}
      />

      <TextField
        label="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        fullWidth
        sx={{ mb: 2 }}
      />

      <TextField
        label="Status"
        value={status}
        onChange={(e) => setStatus(e.target.value)}
        fullWidth
        sx={{ mb: 2 }}
      />

      <Button type="submit" variant="contained" color="primary" fullWidth sx={{ mb: 2 }}>
        Search
      </Button>

      <Button type="button" variant="outlined" color="secondary" fullWidth onClick={handleReset}>
        Reset Filters
      </Button>
    </Box>
  );
};

export default FilterOnFrontend;
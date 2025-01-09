const API_BASE_URL = "http://localhost:8080";

export const fetchAllReleases = async () => {
  const response = await fetch(`${API_BASE_URL}/all`);
  if (!response.ok) {
    throw new Error("Failed to fetch releases");
  }
  return response.json();
};

export const searchReleases = async (query: string) => {
  const response = await fetch(`${API_BASE_URL}/search?${query}`);
  if (!response.ok) {
    throw new Error("Failed to search releases");
  }
  return response.json();
};

export const fetchArtists = async () => {
  const response = await fetch(`${API_BASE_URL}/artists`);
  if (!response.ok) {
    throw new Error("Failed to fetch artists");
  }
  return response.json();
};

export const fetchFormats = async () => {
  const response = await fetch(`${API_BASE_URL}/formats`);
  if (!response.ok) {
    throw new Error("Failed to fetch formats");
  }
  return response.json();
};



export const fetchRelations = async () => {
  const response = await fetch(`${API_BASE_URL}/relations`);
  if (!response.ok) {
    throw new Error("Failed to fetch relations");
  }
  return response.json();
};
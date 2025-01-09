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
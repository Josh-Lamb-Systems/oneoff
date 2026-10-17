// Relative by default so the UI works on whatever host/port (or reverse proxy)
// serves it. Set VITE_API_BASE_URL only when the API lives elsewhere.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "/api";

export { API_BASE_URL };

import { Role } from "@/types/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export type OAuthProvider = "google" | "github";

/**
 * Initiates OAuth flow by redirecting to the backend OAuth endpoint
 */
export const initiateOAuth = (provider: OAuthProvider, role: Role) => {
  // The backend expects role as "volunteer" or "user" (lowercase)
  const roleParam = role.toLowerCase();
  
  // Redirect to backend OAuth signin endpoint
  // The backend will redirect to the OAuth provider
  window.location.href = `${API_URL}/api/auth/oauth/signin?provider=${provider}&role=${roleParam}`;
};

/**
 * Extracts token from OAuth success redirect URL
 */
export const getOAuthToken = (): string | null => {
  if (typeof window === "undefined") return null;
  
  const params = new URLSearchParams(window.location.search);
  return params.get("token");
};

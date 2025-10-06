import { ApiError } from "@/types/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const getCurrentUser = async () => {
    const response = await fetch(`${API_URL}/api/users/me`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });
    if (!response.ok) {
        throw new Error("Failed to fetch current user");
    }
    return response.json()
}

export const login = async (email: string, password: string) => {
    const response = await fetch(`${API_URL}/api/auth/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
        credentials: "include",
    });
    if (!response.ok) {
        const errorData: ApiError  = await response.json();
        throw errorData.error || "Login failed";
    }
    return response.json();
}

export const signup = async (name: string, email: string, password: string, role: string) => {
    const response = await fetch(`${API_URL}/api/auth/signup`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, email, password, role }),
        credentials: "include",
    });
    if (!response.ok) {
        const errorData: ApiError  = await response.json();
        throw errorData.error || "Sign up failed";
    }
    return response.json();
}

export const logout = async () => {
    const response = await fetch(`${API_URL}/api/auth/logout`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });
    if (!response.ok) {
        const errorData: ApiError  = await response.json();
        throw errorData.error || "Logout failed";
    }
    return response.json();
}
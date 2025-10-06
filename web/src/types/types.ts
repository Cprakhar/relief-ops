export enum Role {
    ADMIN = "admin",
    USER = "user",
    VOLUNTEER = "volunteer",
}

export enum AuthMode {
    LOGIN = "login",
    SIGNUP = "signup",
}

export type ApiError = {
    error: string;
}

export type User = {
    id: string;
    name: string;
    email: string;
    role: Role;
    avatar_url?: string;
}

export type FormErrors = {
    [key: string]: string[];
}

export type SignupFormData = {
    name: string;
    email: string;
    password: string;
    confirmPassword: string;
    role?: string;
}

export type LoginFormData = {
    email: string;
    password: string;
    rememberMe?: boolean;
}
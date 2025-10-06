import { login, signup } from "@/api/auth"
import { AuthContext } from "@/context/AuthContext"
import { ApiError, AuthMode, FormErrors, LoginFormData, SignupFormData } from "@/types/types"
import { useRouter } from "next/navigation"
import { FormEvent, useContext, useState } from "react"
import z4 from "zod/v4"

const signupSchema = z4.object({
    name: z4.string().min(1, {error: "Name is required"}),
    email: z4.email({error: "Invalid email address"}),
    password: z4.string().min(8, {error: "Password must be at least 8 characters"}),
    confirmPassword: z4.string()
}).refine((data) => data.password === data.confirmPassword, {
    error: "Passwords do not match",
    path: ["confirmPassword"]
})

const loginSchema = z4.object({
    email: z4.email({error: "Invalid email address"}),
    password: z4.string().min(1, {error: "Password is required"}),
    rememberMe: z4.boolean().optional()
})

const useAuth = (authMode: string) => {
    const router = useRouter()
    const { setUser } = useContext(AuthContext)
    const [loading, setLoading] = useState<boolean>(false)
    const [formErrors, setFormErrors] = useState<FormErrors>({})
    const [error, setError] = useState<string>("")

    const extractFormData = (form: HTMLFormElement): SignupFormData | LoginFormData => {
        const formData = new FormData(form)
        
        if (authMode === AuthMode.SIGNUP) {
            return {
                name: formData.get("name") as string,
                email: formData.get("email") as string,
                password: formData.get("password") as string,
                confirmPassword: formData.get("confirmPassword") as string,
                role: formData.get("role") as string,
            }
        } else {
            return {
                email: formData.get("email") as string,
                password: formData.get("password") as string,
                rememberMe: formData.get("remember-me") === "on",
            }
        }
    }

    const validateForm = (data: SignupFormData | LoginFormData): boolean => {
        const schema = authMode === AuthMode.SIGNUP ? signupSchema : loginSchema
        const result = schema.safeParse(data)
        
        if (!result.success) {
            const fieldErrors: FormErrors = {}
            const flattenedErrors = z4.flattenError<SignupFormData | LoginFormData>(result.error)
            Object.keys(flattenedErrors.fieldErrors).forEach((key) => {
                const messages = flattenedErrors.fieldErrors[key as keyof (SignupFormData | LoginFormData)]
                if (messages) fieldErrors[key] = messages
            })
            setFormErrors(fieldErrors)
            return false
        }
        
        setFormErrors({})
        return true
    }

    const handleFormSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        
        // Clear previous errors
        setError("")
        setFormErrors({})
        
        const form = e.currentTarget
        const formData = extractFormData(form)
        
        if (!validateForm(formData)) {
            return
        }
        
        setLoading(true)
        try {
            if (authMode === AuthMode.LOGIN) {
                const loginData = formData as LoginFormData
                const result = await login(loginData.email, loginData.password)
                setUser(result.data.user)
                router.push("/dashboard")
            } else if (authMode === AuthMode.SIGNUP) {
                const signupData = formData as SignupFormData
                const role = signupData.role || "user"
                await signup(signupData.name, signupData.email, signupData.password, role)
                router.push("/auth")
            }
        } catch (err) {
            if (err && typeof err === "object" && "error" in err) {
                setError((err as ApiError).error)
            } else if (err instanceof Error) {
                setError(err.message)
            } else {
                setError("An unknown error occurred")
            }
        } finally {
            setLoading(false)
        }
    }

    return {
        loading,
        formErrors,
        error,
        handleFormSubmit,
    }
}

export default useAuth
'use client'

import Login from "@/components/Login"
import SignUp from "@/components/SignUp"
import useAuth from "@/hooks/useAuth"
import { AuthMode } from "@/types/types"
import { useState } from "react"

const AuthPage = () => {
  const [mode, setMode] = useState<AuthMode>(AuthMode.LOGIN)
  const { handleFormSubmit, formErrors, error, loading } = useAuth(mode)

  return (
    <div>
    {mode === AuthMode.LOGIN ?
      <Login 
        setAuthMode={setMode} 
        fieldErrors={formErrors} 
        formSubmit={handleFormSubmit}
        error={error}
        loading={loading}
        /> :
      <SignUp 
        setAuthMode={setMode} 
        formSubmit={handleFormSubmit}
        fieldErrors={formErrors}
        error={error}
        loading={loading}
        />
    }
    </div>
  )
}

export default AuthPage
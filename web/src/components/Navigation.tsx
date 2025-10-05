'use client'

import { Role, User } from "@/types/types"
import { 
  Home, 
  AlertTriangle, 
  MapPin, 
  Shield, 
  LogOut, 
  FileText, 
  ChevronDown,
  Menu,
  X
} from "lucide-react"
import Link from "next/link"
import Image from "next/image"
import { useState, useRef, useEffect } from "react"

interface NavigationProps {
  user?: User | null;
  onLogout?: () => void;
}

const Navigation = ({ user, onLogout }: NavigationProps) => {
  const [dropdownOpen, setDropdownOpen] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setDropdownOpen(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  const navLinks = [
    { href: "/dashboard", label: "Home", icon: Home, show: true },
    { 
      href: "/report-disaster", 
      label: "Report a Disaster", 
      icon: AlertTriangle, 
      show: user?.role === Role.ADMIN || user?.role === Role.VOLUNTEER
    },
    { href: "/resources", label: "Resources", icon: MapPin, show: true },
    { 
      href: "/admin", 
      label: "Admin Dashboard", 
      icon: Shield, 
      show: user?.role === Role.ADMIN 
    },
  ]

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  return (
    <nav className="sticky top-0 z-50 bg-white/80 dark:bg-gray-900/80 backdrop-blur-xl border-b border-gray-200 dark:border-gray-800 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/dashboard" className="flex items-center gap-3 group">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/30 group-hover:shadow-blue-500/50 transition-all">
              <AlertTriangle className="w-6 h-6 text-white" />
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent hidden sm:block">
              Relief Ops
            </span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-1">
            {navLinks.map((link) => 
              link.show && (
                <Link
                  key={link.href}
                  href={link.href}
                  className="flex items-center gap-2 px-4 py-2 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all font-medium group"
                >
                  <link.icon className="w-4 h-4 group-hover:scale-110 transition-transform" />
                  <span>{link.label}</span>
                </Link>
              )
            )}
          </div>

          {/* Right Side - Auth/Avatar */}
          <div className="flex items-center gap-3">
            {user ? (
              <div className="relative" ref={dropdownRef}>
                <button
                  onClick={() => setDropdownOpen(!dropdownOpen)}
                  className="flex items-center gap-2 px-3 py-2 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all group"
                >
                  {/* Avatar */}
                  {user.avatar_url ? (
                    <div className="w-9 h-9 rounded-full border-2 border-gray-200 dark:border-gray-700 group-hover:border-blue-500 transition-colors overflow-hidden">
                      <Image
                        src={user.avatar_url}
                        alt={user.username}
                        width={36}
                        height={36}
                        className="w-full h-full object-cover"
                      />
                    </div>
                  ) : (
                    <div className="w-9 h-9 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold text-sm border-2 border-transparent group-hover:border-blue-500 transition-colors">
                      {getInitials(user.username)}
                    </div>
                  )}
                  
                  {/* Username (desktop only) */}
                  <span className="hidden lg:block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {user.username}
                  </span>
                  
                  <ChevronDown 
                    className={`w-4 h-4 text-gray-500 dark:text-gray-400 transition-transform ${
                      dropdownOpen ? 'rotate-180' : ''
                    }`} 
                  />
                </button>

                {/* Dropdown Menu */}
                {dropdownOpen && (
                  <div className="absolute right-0 mt-2 w-56 bg-white dark:bg-gray-800 rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 py-2 animate-in fade-in slide-in-from-top-2 duration-200">
                    {/* User Info */}
                    <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                      <p className="text-sm font-semibold text-gray-900 dark:text-white">
                        {user.username}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                        {user.email}
                      </p>
                      <div className="mt-2">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          user.role === Role.ADMIN 
                            ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300'
                            : user.role === Role.VOLUNTEER
                            ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
                            : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                        }`}>
                          {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                        </span>
                      </div>
                    </div>

                    {/* Menu Items */}
                    <div className="py-1">
                      {user.role === Role.VOLUNTEER && (
                        <Link
                          href="/my-reports"
                          onClick={() => setDropdownOpen(false)}
                          className="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                        >
                          <FileText className="w-4 h-4" />
                          My Reports
                        </Link>
                      )}
                      
                      <button
                        onClick={() => {
                          setDropdownOpen(false)
                          onLogout?.()
                        }}
                        className="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                      >
                        <LogOut className="w-4 h-4" />
                        Logout
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="hidden md:flex items-center gap-2">
                <Link
                  href="/auth"
                  className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 font-medium transition-colors"
                >
                  Sign In
                </Link>
                <Link
                  href="/auth"
                  className="px-4 py-2 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white font-medium rounded-lg shadow-lg shadow-blue-500/30 hover:shadow-blue-500/50 transition-all"
                >
                  Get Started
                </Link>
              </div>
            )}

            {/* Mobile Menu Button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="md:hidden p-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            >
              {mobileMenuOpen ? (
                <X className="w-6 h-6" />
              ) : (
                <Menu className="w-6 h-6" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200 dark:border-gray-800 animate-in slide-in-from-top duration-200">
            <div className="flex flex-col gap-1">
              {navLinks.map((link) => 
                link.show && (
                  <Link
                    key={link.href}
                    href={link.href}
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all font-medium"
                  >
                    <link.icon className="w-5 h-5" />
                    <span>{link.label}</span>
                  </Link>
                )
              )}
              
              {!user && (
                <>
                  <Link
                    href="/auth"
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center justify-center px-4 py-3 mt-4 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg transition-all font-medium"
                  >
                    Sign In
                  </Link>
                  <Link
                    href="/auth"
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center justify-center px-4 py-3 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white rounded-lg shadow-lg shadow-blue-500/30 transition-all font-medium"
                  >
                    Get Started
                  </Link>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  )
}

export default Navigation;
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types';
import { AuthAPI } from '@/services/api';
import { authEventManager } from '@/handlers/api-handlers';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  clearAuth: () => void;
  checkAuthRequired: () => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if user is already logged in on app start
    const initializeAuth = () => {
      try {
        const storedUser = AuthAPI.getCurrentUser();
        const token = AuthAPI.getAccessToken();
        
        if (storedUser && token) {
          setUser(storedUser);
        }
      } catch (error) {
        console.error('Error initializing auth:', error);
        // Clear invalid data
        clearAuth();
      } finally {
        setIsLoading(false);
      }
    };

    initializeAuth();

    // Listen for token invalidation events from API handler
    const handleTokenInvalidation = () => {
      setUser(null);
    };

    authEventManager.addListener(handleTokenInvalidation);

    // Cleanup listener on unmount
    return () => {
      authEventManager.removeListener(handleTokenInvalidation);
    };
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    try {
      const authResponse: AuthResponse = await AuthAPI.login(credentials);
      setUser(authResponse.user);
    } catch (error) {
      throw error;
    }
  };

  const register = async (userData: RegisterRequest): Promise<void> => {
    try {
      const authResponse: AuthResponse = await AuthAPI.register(userData);
      setUser(authResponse.user);
    } catch (error) {
      throw error;
    }
  };

  const logout = async (): Promise<void> => {
    try {
      await AuthAPI.logout();
      setUser(null);
    } catch (error) {
      console.error('Logout error:', error);
      // Still clear user state even if API call fails
      setUser(null);
    }
  };

  const clearAuth = (): void => {
    // Clear all authentication data
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    setUser(null);
  };

  const checkAuthRequired = (): boolean => {
    return !AuthAPI.isAuthenticated();
  };

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
    clearAuth,
    checkAuthRequired,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}; 
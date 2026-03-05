import { defineStore } from 'pinia';

const STORAGE_KEY = 'dnd_auth_token';
const USER_KEY = 'dnd_user';

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  address?: Address | null;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  zip: string;
  country: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  name: string;
}

// Determine API base at runtime
const runtimeHost = typeof window !== 'undefined' ? window.location.hostname : '';
const isLocal = runtimeHost === 'localhost' || runtimeHost === '127.0.0.1' || runtimeHost === 'host.docker.internal';
const API_BASE = import.meta.env.VITE_API_BASE || (isLocal ? 'http://localhost:8001' : 'http://mock-api:8001');

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: null as string | null,
    loading: false,
    error: null as string | null
  }),

  getters: {
    isAuthenticated: (state): boolean => !!state.token && !!state.user,
    currentUser: (state): User | null => state.user
  },

  actions: {
    // Hydrate from localStorage
    hydrate() {
      try {
        const token = localStorage.getItem(STORAGE_KEY);
        const userStr = localStorage.getItem(USER_KEY);
        if (token && userStr) {
          this.token = token;
          this.user = JSON.parse(userStr);
        }
      } catch (err) {
        console.error('Failed to hydrate auth state', err);
        this.clearAuth();
      }
    },

    // Persist to localStorage
    persist() {
      try {
        if (this.token && this.user) {
          localStorage.setItem(STORAGE_KEY, this.token);
          localStorage.setItem(USER_KEY, JSON.stringify(this.user));
        } else {
          localStorage.removeItem(STORAGE_KEY);
          localStorage.removeItem(USER_KEY);
        }
      } catch (err) {
        console.error('Failed to persist auth state', err);
      }
    },

    // Clear auth state
    clearAuth() {
      this.user = null;
      this.token = null;
      this.error = null;
      this.persist();
    },

    // Register new user
    async register(data: RegisterData): Promise<boolean> {
      this.loading = true;
      this.error = null;

      try {
        const response = await fetch(`${API_BASE}/api/auth/register`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(data)
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Registration failed');
        }

        const { user, token } = await response.json();
        this.user = user;
        this.token = token;
        this.persist();
        return true;
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Registration failed';
        return false;
      } finally {
        this.loading = false;
      }
    },

    // Login user
    async login(credentials: LoginCredentials): Promise<boolean> {
      this.loading = true;
      this.error = null;

      try {
        const response = await fetch(`${API_BASE}/api/auth/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(credentials)
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Login failed');
        }

        const { user, token } = await response.json();
        this.user = user;
        this.token = token;
        this.persist();
        return true;
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Login failed';
        return false;
      } finally {
        this.loading = false;
      }
    },

    // Logout user
    async logout(): Promise<void> {
      if (!this.token) return;

      try {
        await fetch(`${API_BASE}/api/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        });
      } catch (err) {
        console.error('Logout request failed', err);
      } finally {
        this.clearAuth();
      }
    },

    // Fetch current user (verify session)
    async fetchCurrentUser(): Promise<boolean> {
      if (!this.token) return false;

      try {
        const response = await fetch(`${API_BASE}/api/auth/me`, {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        });

        if (!response.ok) {
          this.clearAuth();
          return false;
        }

        this.user = await response.json();
        this.persist();
        return true;
      } catch (err) {
        console.error('Failed to fetch current user', err);
        this.clearAuth();
        return false;
      }
    },

    // Update user profile
    async updateProfile(data: { name?: string; address?: Address }): Promise<boolean> {
      if (!this.token) return false;

      this.loading = true;
      this.error = null;

      try {
        const response = await fetch(`${API_BASE}/api/users/profile`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.token}`
          },
          body: JSON.stringify(data)
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Profile update failed');
        }

        this.user = await response.json();
        this.persist();
        return true;
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Profile update failed';
        return false;
      } finally {
        this.loading = false;
      }
    }
  }
});

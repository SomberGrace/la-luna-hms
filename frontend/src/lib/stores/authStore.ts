// frontend/src/lib/stores/authStore.ts
import { writable } from 'svelte/store';
import { isRoleAllowedForApp, type UserRole } from '#lib/mobile';

export type { UserRole } from '#lib/mobile';

export interface UserSession {
  userId: string;
  username: string;
  fullName: string;
  role: UserRole;
  shift: string;
  token: string;
  isLoggedIn: boolean;
}

const DEFAULT_SESSION: UserSession = {
  userId: '',
  username: '',
  fullName: '',
  role: 'frontdesk',
  shift: '',
  token: '',
  isLoggedIn: false
};

const STORAGE_KEY = 'la_luna_session_v2';

function createAuthStore() {
  let initial = DEFAULT_SESSION;
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        initial = JSON.parse(saved);
      } catch (e) {
        // use default
      }
    }
  }

  const { subscribe, set, update } = writable<UserSession>(initial);

  function persist(session: UserSession) {
    if (typeof window !== 'undefined') {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(session));
    }
  }

  return {
    subscribe,
    login: (session: Partial<UserSession>) => {
      update((curr) => {
        const safeRole = session.role ?? curr.role;
        const next = {
          ...curr,
          ...session,
          role: safeRole,
          isLoggedIn: Boolean(session.token || session.username || session.userId) && isRoleAllowedForApp(safeRole)
        };
        persist(next);
        return next;
      });
    },
    logout: () => {
      const reset: UserSession = { ...DEFAULT_SESSION };
      if (typeof window !== 'undefined') {
        localStorage.removeItem(STORAGE_KEY);
        localStorage.removeItem('la_luna_frontdesk_session'); // clean old key
      }
      set(reset);
    },
    updateShift: (newShift: string) => {
      update((curr) => {
        const next = { ...curr, shift: newShift };
        persist(next);
        return next;
      });
    }
  };
}

export const authStore = createAuthStore();

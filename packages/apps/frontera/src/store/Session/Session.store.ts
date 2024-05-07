import type { RootStore } from '@store/root';

import { TransportLayer } from '@store/transport';
import { action, makeAutoObservable } from 'mobx';
import { isHydrated, makePersistable } from 'mobx-persist-store';

// temporary - will be removed once we drop react-query and getGraphQLClient
declare global {
  interface Window {
    __COS_SESSION__?: {
      email: string;
      sessionToken: string | null;
    };
  }
}

type Session = {
  exp: number;
  iat: number;
  tenant: string;
  access_token: string;
  refresh_token: string;
  integrations_token: string;
  profile: {
    id: string;
    name: string;
    email: string;
    locale: string;
    picture: string;
    given_name: string;
    verified_email: boolean;
  };
};

const defaultSession: Session = {
  exp: 0,
  iat: 0,
  tenant: '',
  access_token: '',
  refresh_token: '',
  integrations_token: '',
  profile: {
    id: '',
    name: '',
    email: '',
    locale: '',
    picture: '',
    given_name: '',
    verified_email: false,
  },
};

export class SessionStore {
  value: Session = defaultSession;
  sessionToken: string | null = null;
  error: string | null = null;
  isLoading: 'google' | 'azure-ad' | null = null;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);
    makePersistable(this, {
      name: 'SessionStore',
      properties: ['value', 'sessionToken'],
    }).then(
      action(() => {
        this.loadSession();
      }),
    );
  }

  async loadSession() {
    // Check if the user is already authenticated
    this.isLoading = null;
    if (this.isAuthenticated) {
      // Refresh session data
      await this.fetchSession();

      return;
    }

    // Get the session token from the URL
    const urlParams = new URLSearchParams(window.location.search);
    const sessionToken = urlParams.get('sessionToken');

    if (sessionToken) {
      // Save the session token to the store
      this.sessionToken = sessionToken;
    }
  }

  async fetchSession(options?: {
    onSuccess?: () => void;
    onError?: (error: string) => void;
  }) {
    try {
      const { data } = await this.transportLayer.http.get<{
        session: Session | null;
      }>('/session');
      if (data?.session) {
        this.value = data?.session;
        this.setSessionToWindow();
      }
      options?.onSuccess?.();
    } catch (err) {
      this.error = (err as Error)?.message;
      options?.onError?.(this.error);
      console.error(err);
    }
  }

  async authenticate(provider: 'google' | 'azure-ad') {
    try {
      // initiate the google auth flow
      this.isLoading = provider;
      const { data } = await this.transportLayer.http.get<{ url: string }>(
        '/google-auth',
      );
      window.location.href = data.url;
    } catch (err) {
      this.error = (err as Error)?.message;
    }
  }

  clearSession() {
    this.sessionToken = null;
    this.value = defaultSession;
    this.removeSessionFromWindow();
  }

  /**
   * Temporary: will be removed when we drop react-query & getGraphQLClient
   * Set the session token & session email to the window object
   */
  private setSessionToWindow() {
    window.localStorage.setItem(
      '__COS_SESSION__',
      JSON.stringify({
        email: this.value.profile.email,
        sessionToken: this.sessionToken,
      }),
    );

    window.__COS_SESSION__ = {
      email: this.value.profile.email,
      sessionToken: this.sessionToken,
    };
  }
  private removeSessionFromWindow() {
    window.localStorage.removeItem('__COS_SESSION__');
    delete window.__COS_SESSION__;
  }

  get isHydrated() {
    return isHydrated(this);
  }
  get isAuthenticated() {
    return Boolean(this.sessionToken);
  }
  get isBootstrapped() {
    return this.isHydrated && this.value.profile.email !== '';
  }
}

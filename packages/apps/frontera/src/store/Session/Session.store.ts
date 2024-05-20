import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { isHydrated, makePersistable } from 'mobx-persist-store';
import { action, autorun, runInAction, makeAutoObservable } from 'mobx';

// temporary - will be removed once we drop react-query and getGraphQLClient
declare global {
  interface Window {
    intercomSettings: unknown;
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
  isBootstrapping = true;
  isLoading: 'google' | 'azure-ad' | null = null;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makePersistable(this, {
      name: 'SessionStore',
      properties: ['value', 'sessionToken'],
    }).then(
      action(() => {
        this.loadSession();
      }),
    );

    autorun(() => {
      if (this.sessionToken) {
        this.transport.setHeaders({
          Authorization: `Bearer ${this.sessionToken}`,
          'X-Openline-USERNAME': this.value.profile.email ?? '',
        });
      }
    });

    autorun(() => {
      this.transport.setChannelMeta({
        user_id: this.value.profile.id,
        username: this.value.profile.email,
      });

      window.intercomSettings = {
        api_base: 'https://api-iam.intercom.io',
        app_id: 'pqerb2dx',
        alignment: 'left',
        horizontal_padding: 28,
        vertical_padding: 28,
        name: this.value.profile.name,
        email: this.value.profile.email,
        created_at: `${new Date().valueOf()}`,
      };
    });
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
    const email = urlParams.get('email');
    const id = urlParams.get('id');

    if (sessionToken) {
      // Save the session token & other required data to the store
      this.sessionToken = sessionToken;
      this.value.profile.email = email ?? '';
      this.value.profile.id = id ?? '';

      // refetch session to populate with the rest of the data
      await this.fetchSession();

      return;
    }

    this.isBootstrapping = false;
  }

  async fetchSession(options?: {
    onSuccess?: () => void;
    onError?: (error: string) => void;
  }) {
    try {
      const { data } = await this.transport.http.get<{
        session: Session | null;
      }>('/session');
      runInAction(() => {
        if (data?.session) {
          this.value = data?.session;
          this.setSessionToWindow();
        }
      });
      options?.onSuccess?.();
    } catch (err) {
      this.error = (err as Error)?.message;
      options?.onError?.(this.error);
      console.error(err);
    } finally {
      runInAction(() => {
        this.isBootstrapping = false;
      });
    }
  }

  async authenticate(provider: 'google' | 'azure-ad') {
    try {
      // initiate the google auth flow
      this.isLoading = provider;
      const endpoint =
        provider === 'google' ? '/google-auth' : '/microsoft-auth';
      const { data } = await this.transport.http.get<{ url: string }>(endpoint);
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

  getLocalStorageSession() {
    return window.localStorage.getItem('__COS_SESSION__');
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
    return Boolean(this.sessionToken && this.value.profile.email !== '');
  }
  get isBootstrapped() {
    return this.isHydrated && !this.isBootstrapping;
  }
}

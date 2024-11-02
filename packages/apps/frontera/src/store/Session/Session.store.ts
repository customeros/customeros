import type { RootStore } from '@store/root';

import { AxiosError } from 'axios';
import { Transport } from '@store/transport';
import { Persister } from '@store/persister';
import { toJS, autorun, runInAction, makeAutoObservable } from 'mobx';

import mock from './mock.json';

// temporary - will be removed once we drop react-query and getGraphQLClient
declare global {
  interface Window {
    __COS_SESSION__?: {
      email: string;
      apiKey: string | null;
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
  tenantApiKey: string | null = null;
  error: string | null = null;
  isBootstrapping = true;
  isHydrated = false;
  isLoading: 'google' | 'azure-ad' | null = null;
  private persister = Persister.getSharedInstance('Session');

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);

    this.hydrate();
    autorun(() => {
      if (this.isHydrated) {
        this.loadSession();
      }
    });

    autorun(() => {
      if (this.sessionToken) {
        this.transport.setHeaders({
          Authorization: `Bearer ${this.sessionToken}`,
          'X-Openline-USERNAME': this.value.profile.email ?? '',
        });

        Persister.setTenant(this.value.tenant);
        this.persist();
      }
    });

    autorun(() => {
      this.tenantApiKey = this.root.settings?.tenantApiKey;
    });

    autorun(() => {
      this.transport.setChannelMeta({
        user_id: this.value.profile.id,
        username: this.value.profile.email,
      });
    });
  }

  async loadSession() {
    if (this.root.demoMode) {
      this.value = mock.session as Session;
      this.isBootstrapping = false;
      this.isLoading = null;

      return;
    }

    // Check if the user is already authenticated
    this.isLoading = null;

    if (this.isAuthenticated) {
      // Refresh session data
      await this.fetchSession();

      return;
    }

    const parseJwt = (token: string) => {
      try {
        return JSON.parse(atob(token.split('.')[1]));
      } catch (e) {
        return null;
      }
    };

    // Get the session token from the URL
    const urlParams = new URLSearchParams(window.location.search);
    const sessionToken = urlParams.get('sessionToken') as string;
    const jwtParsed = parseJwt(sessionToken);

    if (sessionToken) {
      // Save the session token & other required data to the store
      runInAction(() => {
        this.sessionToken = sessionToken;
        this.value.tenant = jwtParsed?.tenant ?? '';
        this.value.profile.email = jwtParsed?.profile?.email ?? '';
        this.value.profile.id = jwtParsed?.profile?.id ?? '';
      });

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
      if (err instanceof AxiosError && err.response?.status === 401) {
        this.clearSession();

        window.location.href = '/auth/signin';
      }

      this.error = (err as Error)?.message;
      options?.onError?.(this.error);
      console.error('Error fetching session:', err);
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
        provider === 'google' ? '/google-auth' : '/azure-ad-auth';
      const { data } = await this.transport.http.get<{ url: string }>(endpoint);

      window.location.href = data.url;
    } catch (err) {
      this.error = (err as Error)?.message;
    }
  }

  async clearSession() {
    try {
      this.sessionToken = null;
      this.value = defaultSession;
      this.removeSessionFromWindow();
      this.root.windowManager.clearPersisterKey();

      await this.persister?.clear();
    } catch (e) {
      console.error('Failed clearing persisted data', e);
    }
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
        apiKey: this.root.settings.tenantApiKey,
      }),
    );

    window.__COS_SESSION__ = {
      email: this.value.profile.email,
      sessionToken: this.sessionToken,
      apiKey: this.root.settings.tenantApiKey,
    };
  }

  private removeSessionFromWindow() {
    window.localStorage.removeItem('__COS_SESSION__');
    delete window.__COS_SESSION__;
  }

  get isAuthenticated() {
    if (this.root.demoMode) return true;

    return Boolean(this.sessionToken && this.value.profile.email !== '');
  }

  get isBootstrapped() {
    if (this.root.demoMode) return true;

    return this.isHydrated && !this.isBootstrapping;
  }

  private async persist() {
    try {
      this.persister?.setItem('value', toJS(this.value));
      this.persister?.setItem('sessionToken', this.sessionToken);
      this.persister?.setItem('tenantApiKey', this.tenantApiKey);
    } catch (e) {
      console.error('Failed to persist', e);
    }
  }

  private async hydrate() {
    try {
      const value = await this.persister?.getItem<Session>('value');
      const sessionToken = await this.persister?.getItem<string>(
        'sessionToken',
      );
      const tenantApiKey = await this.persister?.getItem<string>(
        'tenantApiKey',
      );

      runInAction(() => {
        value && (this.value = value);
        sessionToken && (this.sessionToken = sessionToken);
        tenantApiKey && (this.tenantApiKey = tenantApiKey);
        this.isHydrated = true;
      });
    } catch (e) {
      console.error('Failed to hydrate', e);
    }
  }
}

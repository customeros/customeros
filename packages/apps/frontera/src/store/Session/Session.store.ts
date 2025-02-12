import type { RootStore } from '@store/root';

import { set } from 'lodash';
import { match } from 'ts-pattern';
import { AxiosError } from 'axios';
import { Persister } from '@store/persister';
import { Transport } from '@infra/transport.ts';
import { toJS, autorun, runInAction, makeAutoObservable } from 'mobx';

// temporary - will be removed once we drop react-query and getGraphQLClient
declare global {
  interface Window {
    __COS_SESSION__?: {
      email: string;
      tenant: string;
      apiKey: string | null;
      sessionToken: string | null;
    };
  }
}

type Session = {
  exp: number;
  iat: number;
  tenant: string;
  campaign: string;
  access_token: string;
  defaultTenant: string;
  refresh_token: string;
  integrations_token: string;
  profile: {
    id: string;
    name: string;
    email: string;
    locale: string;
    picture: string;
    given_name: string;
    workspaceName?: string;
    verified_email: boolean;
  };
};

const defaultSession: Session = {
  exp: 0,
  iat: 0,
  tenant: '',
  defaultTenant: '',
  campaign: '',
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
    workspaceName: '',
  },
};

export class SessionStore {
  value: Session = defaultSession;
  sessionToken: string | null = null;
  tenantApiKey: string | null = null;
  error: string | null = null;
  isBootstrapping = true;
  isHydrated = false;
  isLoading: 'google' | 'azure-ad' | 'magic-link' | null = null;
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
          'X-Openline-TENANT': this.value.tenant ?? '',
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
    // Check if the user is already authenticated
    this.isLoading = null;

    if (this.isAuthenticated) {
      // Refresh session data
      await this.fetchSession();

      return;
    }

    // Get the session token from the URL
    const urlParams = new URLSearchParams(window.location.search);
    const sessionToken = urlParams.get('sessionToken') as string;
    const magicLinkToken = urlParams.get('mg') as string;
    const jwtParsed = this.parseJwt(sessionToken);

    if (magicLinkToken) {
      await this.validateMagicLinkCode(magicLinkToken);
    }

    if (sessionToken) {
      // Save the session token & other required data to the store
      runInAction(() => {
        this.sessionToken = sessionToken;
        this.value.tenant = jwtParsed?.tenant ?? '';
        this.value.defaultTenant = jwtParsed?.defaultTenant ?? '';
        this.value.profile.email = jwtParsed?.profile?.email ?? '';
        this.value.profile.id = jwtParsed?.profile?.id ?? '';
        this.value.campaign = jwtParsed?.campaign ?? '';
        set(
          this.value.profile,
          'workspaceName',
          jwtParsed?.profile?.workspaceName ?? '',
        );
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
          this.value = data.session;
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

  async authenticate(
    provider: 'google' | 'azure-ad' | 'magic-link',
    payload?: Record<string, unknown>,
    opts?: { onSuccess?: (data?: unknown) => void },
  ) {
    try {
      runInAction(() => {
        this.isLoading = provider;
      });

      const params = new URLSearchParams(window.location.search);
      const from = params.get('from');

      const endpointPath = match(provider)
        .with('google', () => '/google-auth')
        .with('azure-ad', () => '/azure-ad-auth')
        .with('magic-link', () => '/magic-link-auth')
        .otherwise(() => '');

      const endpoint =
        endpointPath + (from ? `?from=${encodeURIComponent(from)}` : '');

      if (provider === 'magic-link') {
        const res = await this.transport.http.post(endpoint, payload);

        opts?.onSuccess?.(res?.data);
      } else {
        const { data } = await this.transport.http.get<{ url: string }>(
          endpoint,
        );

        window.location.href = data.url;
      }
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = null;
      });
    }
  }

  private async validateMagicLinkCode(code: string) {
    try {
      const {
        data: { sessionToken },
      } = await this.transport.http.post<{ sessionToken: string }>(
        '/validate-magic-code',
        { code },
      );

      const jwtParsed = this.parseJwt(sessionToken);

      runInAction(() => {
        this.sessionToken = sessionToken;
        this.value.tenant = jwtParsed?.tenant ?? '';
        this.value.profile.email = jwtParsed?.profile?.email ?? '';
        this.value.profile.id = jwtParsed?.profile?.id ?? '';
        this.value.campaign = jwtParsed?.campaign ?? '';

        this.sessionToken = sessionToken;
      });
    } catch (err) {
      runInAction(() => {
        if (err instanceof AxiosError) {
          this.error = err?.response?.data;
        } else {
          this.error = (err as Error)?.message;
        }
      });
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
        tenant: this.value.tenant,
        email: this.value.profile.email,
        sessionToken: this.sessionToken,
        apiKey: this.root.settings.tenantApiKey,
      }),
    );

    window.__COS_SESSION__ = {
      tenant: this.value.tenant,
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

  private parseJwt(token: string) {
    try {
      return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
      return null;
    }
  }
}

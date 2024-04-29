import { action, makeAutoObservable } from 'mobx';
import { TokenResponse } from '@react-oauth/google';
import { makePersistable } from 'mobx-persist-store';

type SessionValue = {
  id: string;
  name: string;
  email: string;
  avatar: string;
};

export class SessionStore {
  value: SessionValue = {
    id: '',
    name: '',
    email: '',
    avatar: '',
  };
  accessToken: string | null = null;
  error: string | null = null;
  isHydrated = false;

  constructor() {
    makeAutoObservable(this);
    makePersistable(this, {
      name: 'SessionStore',
      properties: ['accessToken', 'value'],
    }).then(action((store) => (this.isHydrated = store.isHydrated)));
  }

  async load(
    tokenResponse: Omit<
      TokenResponse,
      'error' | 'error_description' | 'error_uri'
    >,
  ) {
    this.accessToken = tokenResponse.access_token;

    try {
      const res = await fetch(
        `https://www.googleapis.com/oauth2/v1/userinfo?access_token=${this.accessToken}`,
        {
          headers: {
            Authorization: `Bearer ${this.accessToken}`,
            Accept: 'application/json',
          },
        },
      );

      const data = await res.json();

      this.value.id = data?.id;
      this.value.name = data?.name;
      this.value.email = data?.email;
      this.value.avatar = data?.picture;
    } catch (err) {
      this.error = (err as Error)?.message;
      console.error(err);
    }
  }

  loadError(error: string) {
    this.error = error;
  }

  get isAuthenticated() {
    return Boolean(this.accessToken);
  }
}

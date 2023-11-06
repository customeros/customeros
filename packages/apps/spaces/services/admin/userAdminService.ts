import axios from 'axios';

export interface OAuthToken {
  scope: string;
  expiresAt: Date;
  idToken: string;
  accessToken: string;
  refreshToken: string;
  providerAccountId: string;
}

export interface SignInRequest {
  email: string;
  provider: string;
  oAuthToken: OAuthToken;
}

export function UserSignIn(data: SignInRequest): Promise<unknown> {
  return new Promise((resolve, reject) =>
    axios
      .post((process.env.USER_ADMIN_API_URL as string) + '/signin', data, {
        headers: {
          'X-Openline-API-KEY': process.env.USER_ADMIN_API_KEY as string,
        },
      })
      .then(({ data }) => {
        resolve(data);
      })
      .catch((error) => {
        reject(error);
      }),
  );
}

export function RevokeAccess(
  provider: string,
  data?: unknown,
): Promise<unknown> {
  return new Promise((resolve, reject) =>
    axios
      .post(`/ua/${provider}/revoke`, data)
      .then(({ data }) => {
        resolve(data);
      })
      .catch((error) => {
        reject(error);
      }),
  );
}

import NextAuth from 'next-auth';
import Google from 'next-auth/providers/google';
import {
  OAuthToken,
  SignInRequest,
  UserSignIn,
} from '../../../services/admin/userAdminService';

// This file persists in the pages router, as shown in example in documentation.
// You can find a relevant guide here: https://next-auth.js.org/configuration/nextjs#in-app-router
const providers = [
  Google({
    clientId: process.env.GMAIL_CLIENT_ID as string,
    clientSecret: process.env.GMAIL_CLIENT_SECRET as string,
    authorization: {
      params: {
        prompt: 'consent',
        access_type: 'offline',
        response_type: 'code',
      },
    },
  }),
];

const pages = {
  signIn: '/auth/signin',
};

const callbacks = {
  async jwt({ token, account, profile }: any) {
    // Persist the OAuth access_token and or the user id to the token right after signin
    if (account) {
      token.accessToken = account.access_token;
      token.id = profile.id;
      token.playerIdentityId = account.providerAccountId;

      const oAuthToken: OAuthToken = {
        accessToken: account.access_token,
        refreshToken: account.refresh_token,
        expiresAt: new Date(account.expires_at * 1000),
        scope: account.scope,
        providerAccountId: account.providerAccountId,
        idToken: account.id_token,
      };

      const signInRequest: SignInRequest = {
        email: token.email,
        provider: account.provider,
        oAuthToken: oAuthToken,
      };

      await UserSignIn(signInRequest);
    }
    return token;
  },
  async session({ session, user, token }: any) {
    if (token) {
      session.accessToken = token.accessToken;
      session.user.id = token.id;
      session.user.name = token.name;
      session.user.playerIdentityId = token.playerIdentityId;
    }
    return session;
  },
};

export const authOptions = {
  providers,
  callbacks,
  pages,
};
export default NextAuth(authOptions);

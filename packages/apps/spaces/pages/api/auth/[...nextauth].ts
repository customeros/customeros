import NextAuth, { AuthOptions } from 'next-auth';

import Google, { GoogleProfile } from 'next-auth/providers/google';
import {
  OAuthToken,
  UserSignIn,
  SignInRequest,
} from 'services/admin/userAdminService';

// This file persists in the pages router, as shown in example in documentation.
// You can find a relevant guide here: https://next-auth.js.org/configuration/nextjs#in-app-router

export const authOptions: AuthOptions = {
  providers: [
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
  ],
  callbacks: {
    async jwt({ token, account, profile }) {
      // Persist the OAuth access_token and or the user id to the token right after signin
      if (account) {
        token.accessToken = account.access_token;
        token.id = (profile as GoogleProfile)?.id;
        token.playerIdentityId = account.providerAccountId;

        const oAuthToken: OAuthToken = {
          accessToken: account?.access_token ?? '',
          refreshToken: account?.refresh_token ?? '',
          expiresAt: account?.expires_at
            ? new Date(account.expires_at * 1000)
            : new Date(),
          scope: account.scope ?? '',
          providerAccountId: account.providerAccountId,
          idToken: account.id_token ?? '',
        };

        const signInRequest: SignInRequest = {
          email: token.email ?? '',
          provider: account.provider,
          oAuthToken: oAuthToken,
        };

        await UserSignIn(signInRequest);
      }

      return token;
    },
    async session({ session, token }) {
      if (token) {
        Object.assign(session, {
          accessToken: token.accessToken,
          user: {
            id: token.id,
            name: token.name,
            playerIdentityId: token.playerIdendityId,
          },
        });
      }

      return session;
    },
  },
  pages: {
    signIn: '/auth/signin',
  },
};

export default NextAuth(authOptions);

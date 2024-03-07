import NextAuth, { AuthOptions } from 'next-auth';

import AzureAD from 'next-auth/providers/azure-ad';
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
          access_type: 'offline',
          response_type: 'code',
        },
      },
    }),
    AzureAD({
      name: 'Microsoft',
      clientId: process.env.AZURE_AD_CLIENT_ID as string,
      clientSecret: process.env.AZURE_AD_CLIENT_SECRET as string,
      tenantId: 'common',
    }),
  ],
  callbacks: {
    async jwt({ token, account, profile }) {
      // Persist the OAuth access_token and or the user id to the token right after signin
      if (account) {
        token.accessToken = account.access_token;
        token.id = (profile as GoogleProfile)?.id;
        token.playerIdentityId = account.providerAccountId;

        // Check if the email is available in the profile
        if (profile && profile.email && account.provider === 'google') {
          token.email = profile.email;
        } else if (account.provider === 'azure-ad') {
          // If the email is not available in the profile, fetch it from Microsoft Graph API
          const graphApiResponse = await fetch(
            'https://graph.microsoft.com/v1.0/me',
            {
              headers: {
                Authorization: `Bearer ${token.accessToken}`,
              },
            },
          );
          const graphApiData = await graphApiResponse.json();

          if (graphApiData && graphApiData.userPrincipalName) {
            token.email = graphApiData.userPrincipalName;
          }
        }

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
        session = Object.assign(session, {
          accessToken: token.accessToken,
          user: {
            ...session.user,
            playerIdentityId: token.playerIdentityId,
          },
        });
      }

      return session;
    },
  },
  pages: {
    signIn: '/auth/signin',
  },
  secret: process.env.NEXTAUTH_SECRET,
};

export default NextAuth(authOptions);

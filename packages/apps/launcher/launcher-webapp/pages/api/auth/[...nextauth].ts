import NextAuth, { NextAuthOptions } from "next-auth"
import FusionAuth from "next-auth/providers/fusionauth";

export const authOptions: NextAuthOptions = {
  providers: [
    FusionAuth({
      id: "fusionauth",
      name: "Openline",
      clientId: process.env.NEXTAUTH_OAUTH_CLIENT_ID as string,
      clientSecret: process.env.NEXTAUTH_OAUTH_CLIENT_SECRET as string,
      tenantId: process.env.NEXTAUTH_OAUTH_TENANT_ID as string,
      issuer: process.env.NEXTAUTH_OAUTH_SERVER_URL,
      client: {
        authorization_signed_response_alg: 'HS256',
        id_token_signed_response_alg: 'HS256'
      }
    }),
  ],
  theme: {
    colorScheme: "dark",
  },
  callbacks: {
    async jwt({ token }) {
      token.userRole = "admin"
      return token
    },
  },
}

export default NextAuth(authOptions)

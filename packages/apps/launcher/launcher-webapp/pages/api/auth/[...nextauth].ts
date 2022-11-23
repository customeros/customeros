import NextAuth, { NextAuthOptions } from "next-auth"
import FusionAuth from "next-auth/providers/fusionauth";

export const authOptions: NextAuthOptions = {
  providers: [
    FusionAuth({
      id: "fusionauth",
      name: "Openline",
      clientId: "a31080e4-002c-4ab9-9fb6-5a39ee2e6015",
      clientSecret: "XEBSZzdEke9GZTh1YiAzsirRM6FsB0DwN2R1XaUf_Zg",
      tenantId: "d1cc99c3-9f38-4261-8769-86cddbea71af",
      issuer: "http://localhost:9011",
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

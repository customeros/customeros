import NextAuth from "next-auth"
import Google from "next-auth/providers/google";
export const authOptions = {
    // Configure one or more authentication providers
    providers: [
        Google(
            {
                clientId: process.env.GMAIL_CLIENT_ID as string,
                clientSecret: process.env.GMAIL_CLIENT_SECRET as string,
                authorization: {
                    params: {
                        prompt: "consent",
                        access_type: "offline",
                        response_type: "code",
                    },
                },
            }),
    ],
}
export default NextAuth(authOptions)
import { withAuth } from "next-auth/middleware"
import {NextRequest, NextResponse} from "next/server";

// More on how NextAuth.js middleware works: https://next-auth.js.org/configuration/nextjs#middleware
export default withAuth({
  callbacks: {
    authorized({ req, token }) {
      // `/admin` requires admin role
      if (req.nextUrl.pathname === "/admin") {
        return token?.userRole === "admin"
      }
      // `/me` only requires the user to be logged in
      return !!token
    },
  },
})

export function middleware(request: NextRequest) {
  if (request.nextUrl.pathname.startsWith('/server/registered-apps')) {
    return NextResponse.rewrite(new URL('http://localhost:8070/customer-os/registered-apps', request.url))
  }
}
export const config = {
  matcher: '/server/:path*',
}

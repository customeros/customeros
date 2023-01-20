import {NextRequestWithAuth, withAuth} from "next-auth/middleware"
import {NextResponse} from "next/server";

// More on how NextAuth.js middleware works: https://next-auth.js.org/configuration/nextjs#middleware

export default withAuth(function middleware(request: NextRequestWithAuth) {
        if (!request.nextUrl.pathname.startsWith('/customer-os-api/') && !request.nextUrl.pathname.startsWith('/sa/') && !request.nextUrl.pathname.startsWith('/fs/')) {
            return NextResponse.next()
        } else {
            var newURL = '';
            const requestHeaders = new Headers(request.headers);

            if (request.nextUrl.pathname.startsWith('/customer-os-api/')) {
                newURL = process.env.CUSTOMER_OS_API_PATH + "/" + request.nextUrl.pathname.substring(("/customer-os-api/").length);
                requestHeaders.set('X-Openline-API-KEY', process.env.CUSTOMER_OS_API_KEY as string)
                requestHeaders.set('X-Openline-USERNAME', request.nextauth.token?.email as string)
            } else if (request.nextUrl.pathname.startsWith('/fs/')) {
                newURL = process.env.FILE_STORAGE_API_PATH + "/" + request.nextUrl.pathname.substring(("/fs/").length);
                requestHeaders.set('X-Openline-API-KEY', process.env.FILE_STORAGE_API_KEY as string)
            }else if (request.nextUrl.pathname.startsWith('/sa/')) {
                newURL = process.env.SETTINGS_API_PATH + "/" + request.nextUrl.pathname.substring(("/sa/").length);
                requestHeaders.set('X-Openline-API-KEY', process.env.SETTINGS_API_KEY as string)
                requestHeaders.set('X-Openline-USERNAME', request.nextauth.token?.email as string)
                requestHeaders.set('Content-Type', "application/json")
            }

            if (request.nextUrl.searchParams) {
                newURL = newURL + "?" + request.nextUrl.searchParams.toString()
            }
            console.log("Rewriting url to " + newURL);

            return NextResponse.rewrite(new URL(newURL, request.url),
                {
                    request: {
                        headers: requestHeaders,
                    },
                }
            )
        }
    },
    {
        callbacks: {
            authorized({req, token}) {
                console.log("Got Token: " + JSON.stringify(token));
                if (token) return true
                return false;
            },
        },
    })

export const config = {
    matcher: ['/customer-os-api/(.*)', '/fs/(.*)', '/sa/(.*)'],
}
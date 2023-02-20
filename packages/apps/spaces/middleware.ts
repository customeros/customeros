import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
  console.log('🏷️ ----- : Headers ', request.headers);
  console.log('🏷️ ----- : Cookies ', request.headers.get('cookie'));

  console.log('🏷️ ----- : TEST');
  if (
    !request.nextUrl.pathname.startsWith('/customer-os-api/') &&
    !request.nextUrl.pathname.startsWith('/sa/') &&
    !request.nextUrl.pathname.startsWith('/fs/') &&
    !request.nextUrl.pathname.startsWith('/oasis-api/')
  ) {
    return NextResponse.next();
  }

  return fetch(`${process.env.ORY_SDK_URL}/sessions/whoami`, {
    headers: {
      cookie: request.headers.get('cookie') || '',
    },
  })
    .then((resp) => {
      // there must've been no response (invalid URL or something...)
      if (!resp) {
        return NextResponse.redirect(
          new URL('/api/.ory/ui/login', request.url),
        );
      }

      // the user is not signed in
      if (resp.status === 401) {
        console.log('not signed in');
        return NextResponse.redirect(
          new URL('/api/.ory/ui/login', request.url),
        );
      }

      return resp.json().then((data) => {
        console.log('User is signed in. Proceeding to redirect.');

        return getRedirectUrl(data.identity.traits.email, request);
      });
    })
    .catch((err) => {
      console.log(`Global Session Middleware error: ${JSON.stringify(err)}`);
      if (!err.response) {
        console.log('no response');
        return NextResponse.redirect(
          new URL('/api/.ory/ui/login', request.url),
        );
      }
      switch (err.response?.status) {
        // 422 we need to redirect the user to the location specified in the response
        case 422:
          console.log('422');
          return NextResponse.redirect(
            new URL('/api/.ory/ui/login', request.url),
          );
        //return router.push("/login", { query: { aal: "aal2" } })
        case 401:
          console.log('401');
          // The user is not logged in, so we redirect them to the login page.
          return NextResponse.redirect(
            new URL('/api/.ory/ui/login', request.url),
          );
        case 404:
          console.log('404');
          // the SDK is not configured correctly
          // we set this up so you can debug the issue in the browser
          return NextResponse.redirect(
            new URL('/api/.ory/ui/login', request.url),
          );
        default:
          console.log('default');
          return NextResponse.redirect(
            new URL('/api/.ory/ui/login', request.url),
          );
      }
    });
}

function getRedirectUrl(userName: string, request: NextRequest) {
  var newURL = '';

  const requestHeaders = new Headers(request.headers);

  requestHeaders.set('X-Openline-USERNAME', userName);

  if (request.nextUrl.pathname.startsWith('/customer-os-api/')) {
    newURL =
      process.env.CUSTOMER_OS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/customer-os-api/'.length);

    console.log('🏷️ ----- : API KEYx ', process.env.CUSTOMER_OS_API_KEY);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.CUSTOMER_OS_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/fs/')) {
    newURL =
      process.env.FILE_STORAGE_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/fs/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.FILE_STORAGE_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/sa/')) {
    newURL =
      process.env.SETTINGS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/sa/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.SETTINGS_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/oasis-api/')) {
    newURL =
      process.env.OASIS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/oasis-api/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.OASIS_API_KEY as string,
    );
  }

  if (request.nextUrl.searchParams) {
    newURL = newURL + '?' + request.nextUrl.searchParams.toString();
  }

  console.log('🏷️ ----- requestHeaders: ', requestHeaders);
  return NextResponse.rewrite(new URL(newURL, request.url), {
    request: {
      headers: requestHeaders,
    },
  });
}

export const config = {
  matcher: ['/customer-os-api/(.*)', '/fs/(.*)', '/sa/(.*)', '/oasis-api/(.*)'],
};

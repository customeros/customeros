import { NextRequestWithAuth, withAuth } from 'next-auth/middleware';
import { NextRequest, NextResponse } from 'next/server';

const apiPaths = [
  '/customer-os-api/',
  '/sa/',
  '/fs/',
  '/comms-api/',
  '/oasis-api/',
  '/transcription-api/',
  '/validation-api/',
];

export default withAuth(
  function middleware(request: NextRequestWithAuth) {
    const jwt = request.nextauth.token;
    const nextUrlPath = request.nextUrl.pathname;

    if (!jwt || !jwt.email) return NextResponse.redirect('/auth/signin');

    if (apiPaths.some((path) => nextUrlPath.startsWith(path))) {
      return getRedirectUrl(jwt?.email, jwt?.email, request);
    }

    return NextResponse.next();
  },
  {
    pages: {
      signIn: '/auth/signin',
    },
  },
);

function getRedirectUrl(
  userName: string,
  identityId: string,
  request: NextRequest,
): NextResponse {
  let newURL = '';

  const requestHeaders = new Headers(request.headers);

  if (request.nextUrl.pathname.startsWith('/customer-os-api/')) {
    requestHeaders.set('X-Openline-USERNAME', userName);
    requestHeaders.set('X-Openline-IDENTITY-ID', identityId);
    newURL =
      process.env.CUSTOMER_OS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/customer-os-api/'.length);

    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.CUSTOMER_OS_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/fs/')) {
    requestHeaders.set('X-Openline-USERNAME', userName);
    requestHeaders.set('X-Openline-IDENTITY-ID', identityId);
    newURL =
      process.env.FILE_STORAGE_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/fs/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.FILE_STORAGE_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/sa/')) {
    requestHeaders.set('X-Openline-USERNAME', userName);
    requestHeaders.set('X-Openline-IDENTITY-ID', identityId);
    newURL =
      process.env.SETTINGS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/sa/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.SETTINGS_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/comms-api/')) {
    requestHeaders.set('X-Openline-USERNAME', userName);
    newURL =
      process.env.COMMS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/comms-api/'.length);
    requestHeaders.set(
      'X-Openline-Mail-Api-Key',
      process.env.COMMS_MAIL_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/oasis-api/')) {
    requestHeaders.set('X-Openline-USERNAME', userName);
    requestHeaders.set('X-Openline-IDENTITY-ID', identityId);
    newURL =
      process.env.OASIS_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/oasis-api/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.OASIS_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/transcription-api/')) {
    requestHeaders.set('X-Openline-USERNAME', userName);
    requestHeaders.set('X-Openline-IDENTITY-ID', identityId);
    newURL =
      process.env.TRANSCRIPTION_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/transcription-api/'.length);
    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.TRANSCRIPTION_API_KEY as string,
    );
  } else if (request.nextUrl.pathname.startsWith('/validation-api/')) {
    newURL =
      process.env.VALIDATION_API_PATH +
      '/' +
      request.nextUrl.pathname.substring('/validation-api/'.length);

    requestHeaders.set(
      'X-Openline-API-KEY',
      process.env.VALIDATION_API_KEY as string,
    );
  }

  if (request.nextUrl.searchParams && !request.nextUrl.pathname.startsWith('/comms-api/')) {
    newURL = newURL + '?' + request.nextUrl.searchParams.toString();
  }

  return NextResponse.rewrite(new URL(newURL, request.url), {
    request: {
      headers: requestHeaders,
    },
  });
}

import _, { DefaultSession } from 'next-auth';

declare module 'next-auth' {
  interface Session {
    user: {
      /** The user's postal address. */
      playerIdentityId: string;
    } & DefaultSession['user'];
  }
}

import _, { DefaultSession } from 'next-auth';

declare module 'next-auth' {
  interface Session {
    user: {
      playerIdentityId: string;
    } & DefaultSession['user'];
  }
}

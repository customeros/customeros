import { redirect } from 'next/navigation';

import { getProviders } from 'next-auth/react';
import { getServerSession } from 'next-auth/next';
import { authOptions } from 'pages/api/auth/[...nextauth]';

import SignIn from './signin';

export default async function Page() {
  const session = await getServerSession(authOptions);
  // If the user is already logged in, redirect.
  // Note: Make sure not to redirect to the same page
  // To avoid an infinite loop!
  if (session) {
    return redirect('/organizations');
  }

  const providers = await getProviders();

  return <SignIn providers={providers} />;
}

import { redirect } from 'next/navigation';

export default async function RootPage() {
  redirect('/organizations');
}

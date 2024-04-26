export default async function AuthPageLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className='grid sm:grid-cols-1 md:grid-cols-2 h-[100vh]'>
      {children}
    </div>
  );
}

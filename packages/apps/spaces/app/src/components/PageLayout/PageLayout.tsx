interface PageLayoutProps {
  children: React.ReactNode;
}

export const PageLayout = ({ children }: PageLayoutProps) => {
  return (
    <div
      className='h-screen grid bg-gray-25 grid-cols-[200px_minmax(100px,_1fr)] grid-rows-1fr gap-4 transition-all ease-in-out duration-250'
      style={{ gridTemplateAreas: `"sidebar content"` }}
    >
      {children}
    </div>
  );
};

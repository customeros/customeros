interface PageLayoutProps {
  unstyled?: boolean;
  className?: string;
  children: React.ReactNode;
}

export const PageLayout = ({
  unstyled,
  className,
  children,
}: PageLayoutProps) => {
  if (unstyled) return <div className={className}>{children}</div>;

  return (
    <div
      style={{ gridTemplateAreas: `"sidebar content"` }}
      className='h-screen grid bg-gray-25 grid-cols-[200px_minmax(100px,_1fr)] grid-rows-1fr transition-all ease-in-out duration-250'
    >
      {children}
    </div>
  );
};

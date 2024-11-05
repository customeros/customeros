interface LayoutProps {
  children: React.ReactNode;
}

export const Layout = ({ children }: LayoutProps) => {
  return (
    <div className='py-2 px-6 border-r-[1px] h-full w-[456px] overflow-auto'>
      {children}
    </div>
  );
};

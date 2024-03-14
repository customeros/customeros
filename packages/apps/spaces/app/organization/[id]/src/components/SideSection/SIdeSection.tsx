export const SideSection = ({ children }: { children?: React.ReactNode }) => {
  return (
    <div className='flex h-full min-w-[28rem] drop-shadow-ringPrimary'>
      {children}
    </div>
  );
};

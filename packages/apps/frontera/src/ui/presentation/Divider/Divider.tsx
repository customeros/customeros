import { twMerge } from 'tailwind-merge';

interface DividerProps {
  className?: string;
  children?: React.ReactNode;
}

export const Divider = ({ className, children }: DividerProps) => {
  if (!children) {
    return (
      <div className={twMerge('w-full border-b border-gray-200', className)} />
    );
  }

  return (
    <div className={twMerge('flex gap-2 items-center w-full', className)}>
      <div className='w-full border-b border-gray-200' />
      {children}
      <div className='w-full border-b border-gray-200' />
    </div>
  );
};

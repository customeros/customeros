import { twMerge } from 'tailwind-merge';

interface DividerProps {
  className?: string;
}

export const Divider = ({ className }: DividerProps) => {
  return (
    <div className={twMerge(' w-full border-b border-gray-200', className)} />
  );
};

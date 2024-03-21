import { twMerge } from 'tailwind-merge';

interface SkeletonProps {
  className?: string;
}

export const Skeleton = ({ className }: SkeletonProps) => {
  return (
    <div
      className={twMerge(
        'animate-pulse',
        'bg-gray-200',

        className,
      )}
    />
  );
};

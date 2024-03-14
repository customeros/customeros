import { twMerge } from 'tailwind-merge';

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {}

const defaultClasses = 'animate-pulse bg-gray-300 rounded-sm';

export const Skeleton = (props: SkeletonProps) => (
  <div {...props} className={twMerge(defaultClasses, props.className)} />
);

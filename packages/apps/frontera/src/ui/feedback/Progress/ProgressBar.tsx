import * as RadixProgress from '@radix-ui/react-progress';

import { cn } from '@ui/utils/cn';

export const ProgressBar = ({
  className,
  style,
  value = 0,
  ...props
}: RadixProgress.ProgressProps) => {
  return (
    <RadixProgress.Root
      style={{
        // Fix overflow clipping in Safari
        // https://gist.github.com/domske/b66047671c780a238b51c51ffde8d3a0
        transform: 'translateZ(0)',
        ...style,
      }}
      className={cn(
        'relative overflow-hidden bg-transparent rounded-sm h-0.5 bg-gray-200',
        className,
      )}
      {...props}
    >
      <RadixProgress.Indicator
        className='bg-gray-400 h-full transition-width duration-[660ms] ease-[cubic-bezier(0.65, 0, 0.35, 1)] animate-pulse'
        style={{ width: `${value}%` }}
      />
    </RadixProgress.Root>
  );
};

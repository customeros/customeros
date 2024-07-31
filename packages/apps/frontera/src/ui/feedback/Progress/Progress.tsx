import * as RadixProgress from '@radix-ui/react-progress';

import { cn } from '@ui/utils/cn';

export const Progress = ({
  className,
  style,
  ...props
}: RadixProgress.ProgressProps) => {
  return (
    <RadixProgress.Root
      className={cn(
        'relative overflow-hidden bg-transparent rounded-full h-0.5',
        className,
      )}
      style={{
        // Fix overflow clipping in Safari
        // https://gist.github.com/domske/b66047671c780a238b51c51ffde8d3a0
        transform: 'translateZ(0)',
        ...style,
      }}
      {...props}
    >
      <RadixProgress.Indicator
        style={{ transform: `translateX(-${100 - (props?.value ?? 0)}%)` }}
        className='bg-primary-200 w-full h-full transition-transform duration-[660ms] ease-[cubic-bezier(0.65, 0, 0.35, 1)] animate-pulse'
      />
    </RadixProgress.Root>
  );
};

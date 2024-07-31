import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const LockKeyholeCircle = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M13.7316 13.1947C13.661 12.9831 13.6257 12.8772 13.6276 12.7907C13.6295 12.6996 13.6417 12.6519 13.6836 12.5711C13.7235 12.4942 13.83 12.395 14.0432 12.1967C14.6318 11.649 15 10.8675 15 10C15 8.34315 13.6569 7 12 7C10.3431 7 9 8.34315 9 10C9 10.8675 9.36818 11.649 9.95681 12.1967C10.17 12.395 10.2765 12.4942 10.3164 12.5711C10.3583 12.6519 10.3705 12.6996 10.3724 12.7907C10.3743 12.8772 10.339 12.9831 10.2684 13.1947L9.35099 15.947C9.23249 16.3025 9.17324 16.4803 9.20877 16.6218C9.23987 16.7456 9.31718 16.8529 9.42484 16.9216C9.54783 17 9.7352 17 10.1099 17H13.8901C14.2648 17 14.4522 17 14.5752 16.9216C14.6828 16.8529 14.7601 16.7456 14.7912 16.6218C14.8268 16.4803 14.7675 16.3025 14.649 15.947L13.7316 13.1947Z'
    />
  </svg>
);

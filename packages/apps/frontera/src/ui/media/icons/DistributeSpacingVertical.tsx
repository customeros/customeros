import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const DistributeSpacingVertical = ({
  className,
  ...props
}: IconProps) => (
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
      d='M21 3H3M21 21H3M5 12C5 11.0681 5 10.6022 5.15224 10.2346C5.35523 9.74458 5.74458 9.35523 6.23463 9.15224C6.60218 9 7.06812 9 8 9L16 9C16.9319 9 17.3978 9 17.7654 9.15224C18.2554 9.35523 18.6448 9.74458 18.8478 10.2346C19 10.6022 19 11.0681 19 12C19 12.9319 19 13.3978 18.8478 13.7654C18.6448 14.2554 18.2554 14.6448 17.7654 14.8478C17.3978 15 16.9319 15 16 15L8 15C7.06812 15 6.60218 15 6.23463 14.8478C5.74458 14.6448 5.35523 14.2554 5.15224 13.7654C5 13.3978 5 12.9319 5 12Z'
    />
  </svg>
);

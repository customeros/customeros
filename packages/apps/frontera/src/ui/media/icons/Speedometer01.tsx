import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Speedometer01 = ({ className, ...props }: IconProps) => (
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
      d='M12 2V4.5M12 2C6.47715 2 2 6.47715 2 12M12 2C17.5228 2 22 6.47715 22 12M12 19.5V22M12 22C17.5228 22 22 17.5228 22 12M12 22C6.47715 22 2 17.5228 2 12M4.5 12H2M22 12H19.5M19.0784 19.0784L17.3047 17.3047M4.92163 19.0784L6.69715 17.3029M4.92163 5L6.65808 6.73645M19.0784 5L13.4999 10.5M14 12C14 13.1046 13.1046 14 12 14C10.8954 14 10 13.1046 10 12C10 10.8954 10.8954 10 12 10C13.1046 10 14 10.8954 14 12Z'
    />
  </svg>
);

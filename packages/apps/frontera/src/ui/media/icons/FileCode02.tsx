import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const FileCode02 = ({ className, ...props }: IconProps) => (
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
      d='M5 18.5C5 18.9644 5 19.1966 5.02567 19.3916C5.2029 20.7378 6.26222 21.7971 7.60842 21.9743C7.80337 22 8.03558 22 8.5 22H16.2C17.8802 22 18.7202 22 19.362 21.673C19.9265 21.3854 20.3854 20.9265 20.673 20.362C21 19.7202 21 18.8802 21 17.2V9.98822C21 9.25445 21 8.88757 20.9171 8.5423C20.8436 8.2362 20.7224 7.94356 20.5579 7.67515C20.3724 7.3724 20.113 7.11296 19.5941 6.59411L16.4059 3.40589C15.887 2.88703 15.6276 2.6276 15.3249 2.44208C15.0564 2.27759 14.7638 2.15638 14.4577 2.08289C14.1124 2 13.7455 2 13.0118 2H8.5C8.03558 2 7.80337 2 7.60842 2.02567C6.26222 2.2029 5.2029 3.26222 5.02567 4.60842C5 4.80337 5 5.03558 5 5.5M9 14.5L11.5 12L9 9.5M5 9.5L2.5 12L5 14.5'
    />
  </svg>
);

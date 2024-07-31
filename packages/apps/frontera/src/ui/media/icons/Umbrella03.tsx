import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Umbrella03 = ({ className, ...props }: IconProps) => (
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
      d='M6.25 21.9595L12 12.0002M17 3.33995C12.6868 0.849735 7.28964 1.93783 4.246 5.68314C3.94893 6.0487 3.80039 6.23148 3.75718 6.49336C3.7228 6.70172 3.77373 6.97785 3.88018 7.16024C4.01398 7.38947 4.25111 7.52638 4.72539 7.8002L19.2746 16.2002C19.7489 16.474 19.986 16.6109 20.2514 16.6122C20.4626 16.6132 20.7272 16.5192 20.8905 16.3853C21.0957 16.2169 21.1797 15.9969 21.3477 15.5568C23.0695 11.0483 21.3132 5.83017 17 3.33995ZM17 3.33995C15.0868 2.23538 11.2973 5.21728 8.5359 10.0002M17 3.33995C18.9132 4.44452 18.2255 9.21728 15.4641 14.0002M22 22.0002H2'
    />
  </svg>
);

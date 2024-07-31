import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CloudMoon = ({ className, ...props }: IconProps) => (
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
      d='M16.5 13C19.0768 13 21.2397 11.228 21.8366 8.83604C21.4087 8.94287 20.961 8.99958 20.5 8.99958C17.4624 8.99958 15 6.53715 15 3.49958C15 3.03881 15.0567 2.59128 15.1634 2.16357C12.7717 2.76068 11 4.92345 11 7.50003C11 8.41275 11.2223 9.27355 11.6158 10.0312M5 7V3M3 5H7M6 22C3.79086 22 2 20.2091 2 18C2 16.0221 3.43551 14.3796 5.32148 14.0573C6.12876 11.697 8.3662 10 11 10C13.2882 10 15.2772 11.2809 16.2892 13.1648C16.6744 13.0574 17.0805 13 17.5 13C19.9853 13 22 15.0147 22 17.5C22 19.9853 19.9853 22 17.5 22C13.6667 22 9.83333 22 6 22Z'
    />
  </svg>
);

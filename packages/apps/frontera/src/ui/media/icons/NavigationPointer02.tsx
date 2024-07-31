import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const NavigationPointer02 = ({ className, ...props }: IconProps) => (
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
      d='M5.03685 21.3253C4.45216 21.5826 4.15982 21.7112 3.98042 21.6547C3.8249 21.6057 3.70303 21.484 3.65387 21.3285C3.59717 21.1492 3.72545 20.8567 3.98203 20.2717L11.2634 3.6702C11.495 3.14212 11.6108 2.87808 11.7727 2.79727C11.9133 2.72708 12.0787 2.72708 12.2193 2.79727C12.3812 2.87808 12.497 3.14212 12.7287 3.6702L20.01 20.2717C20.2666 20.8567 20.3949 21.1492 20.3382 21.3285C20.289 21.484 20.1671 21.6057 20.0116 21.6547C19.8322 21.7112 19.5399 21.5826 18.9552 21.3253L12.3182 18.405C12.1995 18.3528 12.1402 18.3267 12.0785 18.3164C12.0239 18.3072 11.9681 18.3072 11.9135 18.3164C11.8519 18.3267 11.7925 18.3528 11.6738 18.405L5.03685 21.3253Z'
    />
  </svg>
);

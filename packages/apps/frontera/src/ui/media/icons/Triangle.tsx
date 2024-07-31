import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Triangle = ({ className, ...props }: IconProps) => (
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
      d='M2.39019 18.0983L10.6151 3.89171C11.0696 3.10655 11.2969 2.71396 11.5935 2.58211C11.8521 2.4671 12.1474 2.4671 12.4061 2.58211C12.7026 2.71396 12.9299 3.10654 13.3844 3.89171L21.6093 18.0983C22.0655 18.8863 22.2936 19.2803 22.2599 19.6037C22.2305 19.8857 22.0827 20.142 21.8534 20.3088C21.5904 20.5 21.1352 20.5 20.2246 20.5H3.77487C2.86435 20.5 2.40908 20.5 2.14613 20.3088C1.91677 20.142 1.769 19.8857 1.73959 19.6037C1.70588 19.2803 1.93398 18.8863 2.39019 18.0983Z'
    />
  </svg>
);

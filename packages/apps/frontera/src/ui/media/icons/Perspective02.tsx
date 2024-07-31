import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Perspective02 = ({ className, ...props }: IconProps) => (
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
      d='M16 5.00007L16 19.0001M10 4.00007L10 20.0001M3 12.0001H21M3 5.98924L3 18.0109C3 19.3749 3 20.0569 3.28134 20.5297C3.52803 20.9442 3.9162 21.2556 4.37434 21.4065C4.89685 21.5785 5.56262 21.4306 6.89418 21.1347L18.4942 18.5569C19.3883 18.3582 19.8354 18.2589 20.1691 18.0185C20.4634 17.8064 20.6945 17.5183 20.8377 17.1849C21 16.807 21 16.3491 21 15.4331V8.56702C21 7.65109 21 7.19312 20.8377 6.8152C20.6945 6.48186 20.4634 6.19373 20.1691 5.98168C19.8354 5.74126 19.3883 5.64191 18.4942 5.44322L6.89418 2.86544C5.56262 2.56954 4.89685 2.42159 4.37434 2.59368C3.9162 2.74457 3.52803 3.05596 3.28134 3.47045C3 3.94318 3 4.6252 3 5.98924Z'
    />
  </svg>
);

import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const NavigationPointerOff02 = ({ className, ...props }: IconProps) => (
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
      d='M10.4712 5.47577L11.2631 3.67044C11.4947 3.14236 11.6105 2.87832 11.7724 2.79751C11.913 2.72732 12.0784 2.72732 12.219 2.79751C12.3809 2.87832 12.4967 3.14236 12.7283 3.67044L15.9006 10.9033M8.94668 8.95153L3.98158 20.272C3.725 20.857 3.59672 21.1495 3.65342 21.3288C3.70258 21.4842 3.82445 21.606 3.97997 21.6549C4.15937 21.7114 4.45171 21.5828 5.0364 21.3256L11.6734 18.4053C11.7921 18.3531 11.8514 18.327 11.9131 18.3166C11.9677 18.3075 12.0235 18.3075 12.0781 18.3166C12.1397 18.327 12.1991 18.3531 12.3178 18.4053L18.9547 21.3256C19.5394 21.5828 19.8318 21.7114 20.0112 21.6549C20.1667 21.606 20.2886 21.4842 20.3377 21.3288C20.3944 21.1495 20.2661 20.857 20.0096 20.272L19.8054 19.8066M22 22L2 2'
    />
  </svg>
);

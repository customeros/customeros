import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Announcement03 = ({ className, ...props }: IconProps) => (
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
      d='M18.5 16C20.433 16 22 13.0899 22 9.5C22 5.91015 20.433 3 18.5 3M18.5 16C16.567 16 15 13.0899 15 9.5C15 5.91015 16.567 3 18.5 3M18.5 16L5.44354 13.6261C4.51605 13.4575 4.05231 13.3731 3.67733 13.189C2.91447 12.8142 2.34636 12.1335 2.11414 11.3159C2 10.914 2 10.4427 2 9.5C2 8.5573 2 8.08595 2.11414 7.68407C2.34636 6.86649 2.91447 6.18577 3.67733 5.81105C4.05231 5.62685 4.51605 5.54254 5.44354 5.3739L18.5 3M5 14L5.39386 19.514C5.43126 20.0376 5.44996 20.2995 5.56387 20.4979C5.66417 20.6726 5.81489 20.8129 5.99629 20.9005C6.20232 21 6.46481 21 6.98979 21H8.7722C9.37234 21 9.67242 21 9.89451 20.8803C10.0897 20.7751 10.2443 20.6081 10.3342 20.4055C10.4365 20.1749 10.4135 19.8757 10.3675 19.2773L10 14.5'
    />
  </svg>
);

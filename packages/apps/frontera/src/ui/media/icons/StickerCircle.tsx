import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const StickerCircle = ({ className, ...props }: IconProps) => (
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
      d='M22.0006 12.1171C22.0006 6.5296 17.471 2 11.8835 2C7.34846 2 3.51036 4.98385 2.22531 9.0954C2.143 9.35878 2.10184 9.49047 2.10572 9.65514C2.10888 9.78904 2.14958 9.95446 2.20891 10.0745C2.28188 10.2222 2.39454 10.3349 2.61986 10.5602L13.4409 21.3807C13.6662 21.606 13.7788 21.7187 13.9265 21.7916C14.0466 21.8509 14.212 21.8916 14.3459 21.8948C14.5106 21.8987 14.6423 21.8575 14.9057 21.7752C19.017 20.49 22.0006 16.652 22.0006 12.1171Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M3.4468 9.7341C3.6873 9.71699 3.93013 9.70829 4.17499 9.70829C9.76253 9.70829 14.2921 14.2379 14.2921 19.8254C14.2921 20.0703 14.2834 20.3131 14.2663 20.5536C14.2364 20.9738 14.2215 21.1839 14.099 21.3137C13.9995 21.4191 13.8298 21.4824 13.6856 21.468C13.508 21.4502 13.3466 21.2887 13.0236 20.9658L3.03464 10.9768C2.71171 10.6539 2.55024 10.4924 2.53246 10.3148C2.51801 10.1706 2.58136 10.0009 2.68675 9.9014C2.81651 9.77892 3.02661 9.76398 3.4468 9.7341Z'
    />
  </svg>
);

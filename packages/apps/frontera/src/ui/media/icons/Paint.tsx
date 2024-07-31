import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Paint = ({ className, ...props }: IconProps) => (
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
      d='M2.99955 13H19.9995M11.9995 3.5L10.4995 2M11.4995 3L20.3682 11.8686C20.7642 12.2646 20.9622 12.4627 21.0364 12.691C21.1016 12.8918 21.1016 13.1082 21.0364 13.309C20.9622 13.5373 20.7642 13.7354 20.3682 14.1314L14.8937 19.6059C13.7056 20.7939 13.1116 21.388 12.4266 21.6105C11.8241 21.8063 11.175 21.8063 10.5725 21.6105C9.88751 21.388 9.29349 20.7939 8.10543 19.6059L4.89366 16.3941C3.70561 15.2061 3.11158 14.612 2.88902 13.9271C2.69324 13.3245 2.69324 12.6755 2.88902 12.0729C3.11158 11.388 3.70561 10.7939 4.89366 9.60589L11.4995 3Z'
    />
  </svg>
);

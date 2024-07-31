import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Film03 = ({ className, ...props }: IconProps) => (
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
      d='M7 2V6M7 16V20M17 4V8M17 18V22M2 6H12M2 16H12M12 8H22M12 18H22M12 20V3.6C12 3.03995 12 2.75992 11.891 2.54601C11.7951 2.35785 11.6422 2.20487 11.454 2.10899C11.2401 2 10.9601 2 10.4 2H6.8C5.11984 2 4.27976 2 3.63803 2.32698C3.07354 2.6146 2.6146 3.07354 2.32698 3.63803C2 4.27976 2 5.11984 2 6.8V15.2C2 16.8802 2 17.7202 2.32698 18.362C2.6146 18.9265 3.07354 19.3854 3.63803 19.673C4.27976 20 5.11984 20 6.8 20H12ZM12 4H17.2C18.8802 4 19.7202 4 20.362 4.32698C20.9265 4.6146 21.3854 5.07354 21.673 5.63803C22 6.27976 22 7.11984 22 8.8V17.2C22 18.8802 22 19.7202 21.673 20.362C21.3854 20.9265 20.9265 21.3854 20.362 21.673C19.7202 22 18.8802 22 17.2 22H13.6C13.0399 22 12.7599 22 12.546 21.891C12.3578 21.7951 12.2049 21.6422 12.109 21.454C12 21.2401 12 20.9601 12 20.4V4Z'
    />
  </svg>
);

import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const ImageUserX = ({ className, ...props }: IconProps) => (
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
      d='M16.5 2.5L21.5 7.5M21.5 2.5L16.5 7.5M22 12V17.2C22 18.8802 22 19.7202 21.673 20.362C21.3854 20.9265 20.9265 21.3854 20.362 21.673C19.7202 22 18.8802 22 17.2 22H6.8C5.11984 22 4.27976 22 3.63803 21.673C3.07354 21.3854 2.6146 20.9265 2.32698 20.362C2 19.7202 2 18.8802 2 17.2V6.8C2 5.11984 2 4.27976 2.32698 3.63803C2.6146 3.07354 3.07354 2.6146 3.63803 2.32698C4.27976 2 5.11984 2 6.8 2H12M2.14551 19.9263C2.61465 18.2386 4.16256 17 5.99977 17H12.9998C13.9291 17 14.3937 17 14.7801 17.0769C16.3669 17.3925 17.6073 18.6329 17.9229 20.2196C17.9998 20.606 17.9998 21.0707 17.9998 22M14 9.5C14 11.7091 12.2091 13.5 10 13.5C7.79086 13.5 6 11.7091 6 9.5C6 7.29086 7.79086 5.5 10 5.5C12.2091 5.5 14 7.29086 14 9.5Z'
    />
  </svg>
);

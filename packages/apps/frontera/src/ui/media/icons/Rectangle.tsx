import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Rectangle = ({ className, ...props }: IconProps) => (
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
      d='M3 10.8C3 9.11984 3 8.27976 3.32698 7.63803C3.6146 7.07354 4.07354 6.6146 4.63803 6.32698C5.27976 6 6.11984 6 7.8 6H16.2C17.8802 6 18.7202 6 19.362 6.32698C19.9265 6.6146 20.3854 7.07354 20.673 7.63803C21 8.27976 21 9.11984 21 10.8V12.2C21 13.8802 21 14.7202 20.673 15.362C20.3854 15.9265 19.9265 16.3854 19.362 16.673C18.7202 17 17.8802 17 16.2 17H7.8C6.11984 17 5.27976 17 4.63803 16.673C4.07354 16.3854 3.6146 15.9265 3.32698 15.362C3 14.7202 3 13.8802 3 12.2V10.8Z'
    />
  </svg>
);

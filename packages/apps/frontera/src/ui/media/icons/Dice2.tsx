import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Dice2 = ({ className, ...props }: IconProps) => (
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
      d='M15.25 8.75H15.255M8.75 15.25H8.755M7.8 21H16.2C17.8802 21 18.7202 21 19.362 20.673C19.9265 20.3854 20.3854 19.9265 20.673 19.362C21 18.7202 21 17.8802 21 16.2V7.8C21 6.11984 21 5.27976 20.673 4.63803C20.3854 4.07354 19.9265 3.6146 19.362 3.32698C18.7202 3 17.8802 3 16.2 3H7.8C6.11984 3 5.27976 3 4.63803 3.32698C4.07354 3.6146 3.6146 4.07354 3.32698 4.63803C3 5.27976 3 6.11984 3 7.8V16.2C3 17.8802 3 18.7202 3.32698 19.362C3.6146 19.9265 4.07354 20.3854 4.63803 20.673C5.27976 21 6.11984 21 7.8 21ZM15.5 8.75C15.5 8.88807 15.3881 9 15.25 9C15.1119 9 15 8.88807 15 8.75C15 8.61193 15.1119 8.5 15.25 8.5C15.3881 8.5 15.5 8.61193 15.5 8.75ZM9 15.25C9 15.3881 8.88807 15.5 8.75 15.5C8.61193 15.5 8.5 15.3881 8.5 15.25C8.5 15.1119 8.61193 15 8.75 15C8.88807 15 9 15.1119 9 15.25Z'
    />
  </svg>
);

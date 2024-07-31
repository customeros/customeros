import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const MoonStar = ({ className, ...props }: IconProps) => (
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
      d='M18 2L18.6178 3.23558C18.8833 3.76656 19.016 4.03205 19.1934 4.26211C19.3507 4.46626 19.5337 4.64927 19.7379 4.80664C19.9679 4.98397 20.2334 5.11672 20.7644 5.38221L22 6L20.7644 6.61779C20.2334 6.88328 19.9679 7.01603 19.7379 7.19336C19.5337 7.35073 19.3507 7.53374 19.1934 7.73789C19.016 7.96795 18.8833 8.23344 18.6178 8.76442L18 10L17.3822 8.76442C17.1167 8.23344 16.984 7.96795 16.8066 7.73789C16.6493 7.53374 16.4663 7.35073 16.2621 7.19336C16.0321 7.01603 15.7666 6.88328 15.2356 6.61779L14 6L15.2356 5.38221C15.7666 5.11672 16.0321 4.98397 16.2621 4.80664C16.4663 4.64927 16.6493 4.46626 16.8066 4.26211C16.984 4.03205 17.1167 3.76656 17.3822 3.23558L18 2Z'
    />
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M21 13.3893C19.689 15.689 17.2145 17.2395 14.3779 17.2395C10.1711 17.2395 6.76075 13.8292 6.76075 9.62233C6.76075 6.78554 8.31149 4.31094 10.6115 3C5.77979 3.45812 2 7.52692 2 12.4785C2 17.7371 6.26292 22 11.5215 22C16.4729 22 20.5415 18.2206 21 13.3893Z'
    />
  </svg>
);

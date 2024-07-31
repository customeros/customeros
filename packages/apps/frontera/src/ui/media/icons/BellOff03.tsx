import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const BellOff03 = ({ className, ...props }: IconProps) => (
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
      d='M15 19C15 20.6569 13.6569 22 12 22C10.3431 22 9 20.6569 9 19M7.37748 7.88479C6.49088 8.81326 6 9.9847 6 11.2C6 13.4818 5.43413 15.1506 4.72806 16.3447C3.92334 17.7056 3.52098 18.3861 3.53686 18.5486C3.55504 18.7346 3.58852 18.7933 3.73934 18.9036C3.87117 19 4.53351 19 5.85819 19H19.88M12 6C11.7071 6 11.4164 6.01856 11.13 6.05493C10.7485 6.10339 10.5577 6.12762 10.3662 6.07557C10.2284 6.0381 10.0206 5.91728 9.91979 5.81604C9.77982 5.67541 9.74922 5.60123 9.68801 5.45287C9.56684 5.15921 9.5 4.83741 9.5 4.5C9.5 3.11929 10.6193 2 12 2C13.3807 2 14.5 3.11929 14.5 4.5C14.5 5.17562 14.232 5.78864 13.7965 6.23856C13.2203 6.08184 12.615 6 12 6ZM12 6C13.5913 6 15.1174 6.54786 16.2426 7.52304C17.3679 8.49823 18 9.82087 18 11.2C18 11.5348 18.0091 11.8563 18.0264 12.1652M21 20L3 4'
    />
  </svg>
);

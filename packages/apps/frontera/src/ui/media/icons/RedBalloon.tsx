import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const RedBalloon = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 23 50'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fill='#F04438'
      fillRule='evenodd'
      clipRule='evenodd'
      d='M13.6957 30.0681H14.3043V49.913H13.6957V30.0681ZM15.8938 28.1358C22.7301 27.2027 28 21.2886 28 14.132C28 6.32711 21.732 0 14 0C6.26801 0 0 6.32711 0 14.132C0 21.2886 5.26993 27.2027 12.1062 28.1358C11.1481 28.8365 9.64987 30.0681 10.9565 30.0681C12.7826 30.0681 12.7826 29.4667 12.7826 29.4667C12.7826 29.4667 13.3913 30.0681 14 30.0681C14.6087 30.0681 15.2174 29.4667 15.2174 29.4667C15.2174 29.4667 15.2174 30.0681 17.0435 30.0681C18.3501 30.0681 16.8519 28.8365 15.8938 28.1358ZM9.13043 2.40545C9.13043 2.40545 4.26087 4.81089 3.04348 8.41907C3.04348 12.6286 3.65217 7.81771 9.13043 4.20953C14.6087 0.601363 9.13043 2.40545 9.13043 2.40545Z'
    />
  </svg>
);

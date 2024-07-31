import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const CustomerOs = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 32 32'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      fillRule='evenodd'
      clipRule='evenodd'
      fill='url(#paint0_angular_6038_109640)'
      d='M16 30C23.732 30 30 23.732 30 16C30 8.26801 23.732 2 16 2C8.26801 2 2 8.26801 2 16C2 23.732 8.26801 30 16 30ZM16 23C19.866 23 23 19.866 23 16C23 12.134 19.866 9 16 9C12.134 9 9 12.134 9 16C9 19.866 12.134 23 16 23Z'
    />
    <defs>
      <radialGradient
        r='1'
        cx='0'
        cy='0'
        gradientUnits='userSpaceOnUse'
        id='paint0_angular_6038_109640'
        gradientTransform='translate(16 21.4099) rotate(90) scale(19.4052 37.0568)'
      >
        <stop offset='0.00188747' stopColor='#4C375A' />
        <stop offset='0.00518777' stopColor='#464068' />
        <stop offset='0.0123669' stopColor='#294868' />
        <stop offset='0.031634' stopColor='#185070' />
        <stop offset='0.0961265' stopColor='#1F688C' />
        <stop offset='0.116357' stopColor='#21759A' />
        <stop offset='0.259931' stopColor='#3DA5BE' />
        <stop offset='0.417431' stopColor='#ACD0D9' />
        <stop offset='0.455328' stopColor='#EFF6F8' />
        <stop offset='0.468935' stopColor='white' />
        <stop offset='0.475924' stopColor='#FFF4D7' />
        <stop offset='0.48356' stopColor='#FEE7A6' />
        <stop offset='0.496565' stopColor='#F7BE33' />
        <stop offset='0.50196' stopColor='#DE7324' />
        <stop offset='0.611997' stopColor='#E9571F' />
        <stop offset='0.741661' stopColor='#DF353C' />
        <stop offset='0.932978' stopColor='#BF1B4E' />
        <stop offset='0.991018' stopColor='#8E2C56' />
      </radialGradient>
    </defs>
  </svg>
);

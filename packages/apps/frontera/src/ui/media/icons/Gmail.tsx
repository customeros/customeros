import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Gmail = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 18 13'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <g id='gmail'>
      <path
        id='Vector'
        fill='#EA4335'
        d='M13.5384 1.39221L9.04807 4.89653L4.45508 1.39221V1.39316L4.46063 1.39789V6.30488L8.99628 9.88497L13.5384 6.44316V1.39221Z'
      />
      <path
        id='Vector_2'
        fill='#FBBC05'
        d='M14.7178 0.539789L13.5386 1.39219V6.44314L17.2492 3.59423V1.87806C17.2492 1.87806 16.7988 -0.573068 14.7178 0.539789Z'
      />
      <path
        id='Vector_3'
        fill='#34A853'
        d='M13.5386 6.44327V12.9944H16.3826C16.3826 12.9944 17.1919 12.9111 17.2501 11.9886V3.59436L13.5386 6.44327Z'
      />
      <path
        id='Vector_4'
        fill='#C5221F'
        d='M4.46084 12.9999V6.30478L4.45508 6.30005L4.46084 12.9999Z'
      />
      <path
        id='Vector_5'
        fill='#C5221F'
        d='M4.4551 1.39315L3.28234 0.54549C1.20135 -0.567367 0.75 1.88281 0.75 1.88281V3.59898L4.4551 6.30014V1.39315Z'
      />
      <path
        id='Vector_6'
        fill='#C5221F'
        d='M4.45508 1.39307V6.30006L4.46084 6.30479V1.3978L4.45508 1.39307Z'
      />
      <path
        id='Vector_7'
        fill='#4285F4'
        d='M0.75 3.59998V11.9942C0.807343 12.9177 1.61754 13.0001 1.61754 13.0001H4.46158L4.4551 6.30019L0.75 3.59998Z'
      />
    </g>
  </svg>
);

import * as React from 'react';
import { SVGProps } from 'react';
const SvgServiceProvider = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17.72 11.05a5.728 5.728 0 0 0-5.725-5.73 5.727 5.727 0 0 0-5.725 5.73v3.81a5.727 5.727 0 0 0 5.725 5.73 5.728 5.728 0 0 0 5.725-5.73v-3.81ZM17.73 9.14h1.91A2.86 2.86 0 0 1 22.5 12v1.91a2.86 2.86 0 0 1-2.86 2.86h-1.91V9.14Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M6.27 16.77H4.36a2.86 2.86 0 0 1-2.86-2.86V12a2.86 2.86 0 0 1 2.86-2.86h1.91v7.63ZM4.36 9.14a7.64 7.64 0 0 1 15.28 0M19.64 16.77v1a4.78 4.78 0 0 1-4.78 4.77'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgServiceProvider;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgHome = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M20 11.75a.74.74 0 0 1-.45-.15L12 5.94 4.45 11.6a.75.75 0 0 1-1.05-.15.75.75 0 0 1 .15-1.05l8-6a.75.75 0 0 1 .9 0l8 6a.75.75 0 0 1 .15 1 .74.74 0 0 1-.6.35z' />
    <path d='M18 19.75H6a.76.76 0 0 1-.75-.75V9.5A.75.75 0 0 1 6 8.75a.75.75 0 0 1 .75.75v8.75h10.5V9.5a.75.75 0 0 1 .75-.75.75.75 0 0 1 .75.75V19a.76.76 0 0 1-.75.75z' />
    <path d='M14 19.75a.76.76 0 0 1-.75-.75v-6.25h-2.5V19a.75.75 0 0 1-.75.75.75.75 0 0 1-.75-.75v-7a.76.76 0 0 1 .75-.75h4a.76.76 0 0 1 .75.75v7a.76.76 0 0 1-.75.75z' />
  </svg>
);
export default SvgHome;

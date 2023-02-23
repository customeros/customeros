import * as React from 'react';
import { SVGProps } from 'react';
const SvgPrint = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M18 16.75h-2a.75.75 0 1 1 0-1.5h2A1.25 1.25 0 0 0 19.25 14v-4A1.25 1.25 0 0 0 18 8.75H6A1.25 1.25 0 0 0 4.75 10v4A1.25 1.25 0 0 0 6 15.25h2a.75.75 0 1 1 0 1.5H6A2.75 2.75 0 0 1 3.25 14v-4A2.75 2.75 0 0 1 6 7.25h12A2.75 2.75 0 0 1 20.75 10v4A2.75 2.75 0 0 1 18 16.75Z' />
      <path d='M16 8.75a.76.76 0 0 1-.75-.75V4.75h-6.5V8a.75.75 0 0 1-1.5 0V4.5A1.25 1.25 0 0 1 8.5 3.25h7a1.25 1.25 0 0 1 1.25 1.25V8a.76.76 0 0 1-.75.75ZM15.5 20.75h-7a1.25 1.25 0 0 1-1.25-1.25v-7a1.25 1.25 0 0 1 1.25-1.25h7a1.25 1.25 0 0 1 1.25 1.25v7a1.25 1.25 0 0 1-1.25 1.25Zm-6.75-1.5h6.5v-6.5h-6.5v6.5Z' />
    </g>
  </svg>
);
export default SvgPrint;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgWindowMinimize = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M18 3.25H6A2.75 2.75 0 0 0 3.25 6v6a.75.75 0 1 0 1.5 0V6A1.25 1.25 0 0 1 6 4.75h12A1.25 1.25 0 0 1 19.25 6v12A1.25 1.25 0 0 1 18 19.25h-6a.75.75 0 1 0 0 1.5h6A2.75 2.75 0 0 0 20.75 18V6A2.75 2.75 0 0 0 18 3.25Z' />
      <path d='M11.21 13.19c.092.04.19.06.29.06h4a.75.75 0 1 0 0-1.5h-2.19l3.22-3.22a.75.75 0 0 0-1.06-1.06l-3.22 3.22V8.5a.75.75 0 1 0-1.5 0v4c0 .1.02.198.06.29a.71.71 0 0 0 .4.4ZM8 14.25H5A1.76 1.76 0 0 0 3.25 16v3A1.76 1.76 0 0 0 5 20.75h3A1.76 1.76 0 0 0 9.75 19v-3A1.76 1.76 0 0 0 8 14.25ZM8.25 19a.25.25 0 0 1-.25.25H5a.25.25 0 0 1-.25-.25v-3a.25.25 0 0 1 .25-.25h3a.25.25 0 0 1 .25.25v3Z' />
    </g>
  </svg>
);
export default SvgWindowMinimize;

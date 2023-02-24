import * as React from 'react';
import { SVGProps } from 'react';
const SvgSlidersH = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M19 8.25h-7a.75.75 0 1 1 0-1.5h7a.75.75 0 1 1 0 1.5ZM10 8.25H5a.75.75 0 0 1 0-1.5h5a.75.75 0 1 1 0 1.5Z' />
      <path d='M10 9.75A.76.76 0 0 1 9.25 9V6a.75.75 0 0 1 1.5 0v3a.76.76 0 0 1-.75.75ZM19 17.25h-7a.75.75 0 1 1 0-1.5h7a.75.75 0 1 1 0 1.5ZM10 17.25H5a.75.75 0 1 1 0-1.5h5a.75.75 0 1 1 0 1.5Z' />
      <path d='M10 18.75a.76.76 0 0 1-.75-.75v-3a.75.75 0 1 1 1.5 0v3a.76.76 0 0 1-.75.75ZM19 12.75h-3a.75.75 0 1 1 0-1.5h3a.75.75 0 1 1 0 1.5ZM14 12.75H5a.75.75 0 1 1 0-1.5h9a.75.75 0 1 1 0 1.5Z' />
      <path d='M14 14.25a.76.76 0 0 1-.75-.75v-3a.75.75 0 1 1 1.5 0v3a.76.76 0 0 1-.75.75Z' />
    </g>
  </svg>
);
export default SvgSlidersH;

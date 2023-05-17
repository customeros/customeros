import * as React from 'react';
import { SVGProps } from 'react';
const SvgArrowRight = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M13 18.75a.74.74 0 0 1-.53-.22.75.75 0 0 1 0-1.06L17.94 12l-5.47-5.47a.75.75 0 0 1 1.06-1.06l6 6a.75.75 0 0 1 0 1.06l-6 6a.74.74 0 0 1-.53.22z' />
      <path d='M19 12.75H5a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5z' />
    </g>
  </svg>
);
export default SvgArrowRight;

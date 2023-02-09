import * as React from 'react';
import { SVGProps } from 'react';
const SvgArrowLeft = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M11 18.75a.74.74 0 0 1-.53-.22l-6-6a.75.75 0 0 1 0-1.06l6-6a.75.75 0 0 1 1.06 1.06L6.06 12l5.47 5.47a.75.75 0 0 1-.53 1.28Z' />
      <path d='M19 12.75H5a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgArrowLeft;

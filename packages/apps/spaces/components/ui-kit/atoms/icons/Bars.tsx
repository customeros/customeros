import * as React from 'react';
import { SVGProps } from 'react';
const SvgBars = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M19 12.75H5a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5ZM19 8.25H5a.75.75 0 0 1 0-1.5h14a.75.75 0 1 1 0 1.5ZM19 17.25H5a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgBars;

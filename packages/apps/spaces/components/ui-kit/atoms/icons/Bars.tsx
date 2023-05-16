import * as React from 'react';
import { SVGProps } from 'react';
const SvgBars = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19 12.75H5a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5zm0-4.5H5a.75.75 0 0 1 0-1.5h14a.75.75 0 1 1 0 1.5zm0 9H5a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5z'
      fill='currentColor'
    />
  </svg>
);
export default SvgBars;

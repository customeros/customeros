import * as React from 'react';
import { SVGProps } from 'react';
const SvgCheckCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M10.5 15.25A.74.74 0 0 1 10 15l-3-3a.75.75 0 0 1 1-1l2.47 2.47L19 5a.75.75 0 0 1 1 1l-9 9a.74.74 0 0 1-.5.25z' />
      <path d='M12 21a9 9 0 0 1-7.87-4.66 8.67 8.67 0 0 1-1.07-3.41 9 9 0 0 1 4.6-8.81 8.67 8.67 0 0 1 3.41-1.07 8.86 8.86 0 0 1 3.55.34.75.75 0 1 1-.43 1.43 7.62 7.62 0 0 0-3-.28 7.5 7.5 0 0 0-5.04 2.73 7.42 7.42 0 0 0-1.64 5.51 7.499 7.499 0 0 0 2.73 5.04 7.42 7.42 0 0 0 5.51 1.64 7.5 7.5 0 0 0 5.04-2.73 7.421 7.421 0 0 0 1.64-5.51.788.788 0 1 1 1.57-.15 9 9 0 0 1-4.61 8.81A8.67 8.67 0 0 1 12.93 21H12z' />
    </g>
  </svg>
);
export default SvgCheckCircle;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgEuro = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M13 20.75h-.15a8.75 8.75 0 0 1-5.93-15 8.54 8.54 0 0 1 6.23-2.46 8.75 8.75 0 0 1 6 2.5.75.75 0 0 1-1.06 1.06 7.26 7.26 0 1 0-.19 10.53l.22-.21a.79.79 0 0 1 1.09 0 .7.7 0 0 1 .05 1l-.05.05-.29.28A8.72 8.72 0 0 1 13 20.75Z' />
      <path d='M17 11.25H3a.75.75 0 1 1 0-1.5h14a.75.75 0 1 1 0 1.5ZM15.5 14.25H3a.75.75 0 1 1 0-1.5h12.5a.75.75 0 1 1 0 1.5Z' />
    </g>
  </svg>
);
export default SvgEuro;

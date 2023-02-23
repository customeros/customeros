import * as React from 'react';
import { SVGProps } from 'react';
const SvgServer = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M20.5 8.5v-3a1 1 0 0 0-1-1h-15a1 1 0 0 0-1 1v3a1 1 0 0 0 1 1 1 1 0 0 0-1 1v3a1 1 0 0 0 1 1 1 1 0 0 0-1 1v3a1 1 0 0 0 1 1h15a1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1 1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1 1 1 0 0 0 1-1Zm-1 10h-15v-3h15v3Zm0-5h-15v-3h15v3Zm0-5h-15v-3h15v3Z' />
      <path d='M6.25 7.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5ZM8.75 7.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5ZM6.25 12.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5ZM8.75 12.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5ZM6.25 17.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5ZM8.75 17.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Z' />
    </g>
  </svg>
);
export default SvgServer;

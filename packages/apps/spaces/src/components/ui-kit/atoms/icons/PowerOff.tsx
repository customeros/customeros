import * as React from 'react';
import { SVGProps } from 'react';
const SvgPowerOff = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 21A9 9 0 0 1 5.64 5.64a.74.74 0 0 1 1.06 0 .75.75 0 0 1 0 1.06 7.5 7.5 0 1 0 10.6 10.6 7.48 7.48 0 0 0 0-10.6.75.75 0 0 1 0-1.06.741.741 0 0 1 1.06 0A9 9 0 0 1 12 21Z' />
      <path d='M12 12.75a.76.76 0 0 1-.75-.75V4a.75.75 0 1 1 1.5 0v8a.76.76 0 0 1-.75.75Z' />
    </g>
  </svg>
);
export default SvgPowerOff;

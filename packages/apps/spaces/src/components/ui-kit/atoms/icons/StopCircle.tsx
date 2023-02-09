import * as React from 'react';
import { SVGProps } from 'react';
const SvgStopCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 21a9 9 0 1 1 0-18 9 9 0 0 1 0 18Zm0-16.5a7.5 7.5 0 1 0 0 15 7.5 7.5 0 0 0 0-15Z' />
      <path d='M14.5 8h-5A1.5 1.5 0 0 0 8 9.5v5A1.5 1.5 0 0 0 9.5 16h5a1.5 1.5 0 0 0 1.5-1.5v-5A1.5 1.5 0 0 0 14.5 8Z' />
    </g>
  </svg>
);
export default SvgStopCircle;

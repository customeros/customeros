import * as React from 'react';
import { SVGProps } from 'react';
const SvgMicrophone = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    stroke='currentColor'
    strokeWidth={2}
    {...props}
  >
    <g strokeLinecap='round' strokeLinejoin='round'>
      <path d='M8 23h8M12 19v4' />
    </g>
    <path
      d='M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 3 3 3 3 0 0 0 3-3V4a3 3 0 0 0-3-3z'
      strokeLinejoin='round'
    />
    <path
      d='M19 10v2a7 7 0 0 1-7 7 7 7 0 0 1-7-7v-2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgMicrophone;

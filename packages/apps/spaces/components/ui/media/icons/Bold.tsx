import * as React from 'react';
import { SVGProps } from 'react';
const SvgBold = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M4 8h5.333a2.667 2.667 0 1 0 0-5.333H4V8Zm0 0h6a2.667 2.667 0 1 1 0 5.333H4V8Z'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgBold;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgXClose = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={16}
    height={16}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='m12 4-8 8m0-8 8 8'
      stroke='#98A2B3'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgXClose;

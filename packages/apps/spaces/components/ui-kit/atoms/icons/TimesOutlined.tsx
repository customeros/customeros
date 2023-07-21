import * as React from 'react';
import { SVGProps } from 'react';
const SvgTimesOutlined = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect
      x={0.75}
      y={0.75}
      width={22.5}
      height={22.5}
      rx={11.25}
      stroke='currentColor'
      strokeWidth={1.5}
    />
    <path
      d='m16 8-8 8m0-8 8 8'
      stroke='currentColor'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgTimesOutlined;

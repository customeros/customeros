import * as React from 'react';
import { SVGProps } from 'react';
const SvgIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={14}
    height={14}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M5 8.334 1.667 5m0 0L5 1.667M1.667 5h4.266c2.24 0 3.36 0 4.216.436a4 4 0 0 1 1.748 1.748c.436.856.436 1.976.436 4.216v.934'
      stroke='#98A2B3'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgIcon;

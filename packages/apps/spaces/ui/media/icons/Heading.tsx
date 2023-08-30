import * as React from 'react';
import { SVGProps } from 'react';
const SvgHeading = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M4 2.667v10.666m8-10.666v10.666M5.333 2.667H2.667M12 8H4m1.333 5.333H2.667m10.666 0h-2.666m2.666-10.666h-2.666'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgHeading;

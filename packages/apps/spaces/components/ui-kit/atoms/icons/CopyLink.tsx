import * as React from 'react';
import { SVGProps } from 'react';
const SvgCopyLink = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 20 20'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M8.333 10.833a4.166 4.166 0 0 0 6.284.45l2.5-2.5a4.167 4.167 0 0 0-5.892-5.891L9.792 4.317m1.875 4.85a4.167 4.167 0 0 0-6.284-.45l-2.5 2.5a4.166 4.166 0 0 0 5.892 5.891l1.425-1.425'
      stroke='currentColor'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgCopyLink;

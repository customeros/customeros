import * as React from 'react';
import { SVGProps } from 'react';
const SvgAngleRight = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M10.25 16.25a.74.74 0 0 1-.53-.25.75.75 0 0 1 0-1.06l3-3-3-3A.75.75 0 0 1 10.78 8l3.5 3.5a.75.75 0 0 1 0 1.06L10.78 16a.74.74 0 0 1-.53.25Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgAngleRight;

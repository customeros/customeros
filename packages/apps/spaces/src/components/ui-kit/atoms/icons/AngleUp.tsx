import * as React from 'react';
import { SVGProps } from 'react';
const SvgAngleUp = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M15.5 14.5a.74.74 0 0 1-.53-.22l-3-3-3 3A.75.75 0 0 1 8 13.22l3.5-3.5a.75.75 0 0 1 1.06 0l3.5 3.5a.75.75 0 0 1-.56 1.28Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgAngleUp;

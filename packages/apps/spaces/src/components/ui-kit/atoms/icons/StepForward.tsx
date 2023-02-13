import * as React from 'react';
import { SVGProps } from 'react';
const SvgStepForward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M8 19.75a.75.75 0 0 1-.29-.06.74.74 0 0 1-.46-.69V5a.74.74 0 0 1 .46-.69.75.75 0 0 1 .82.16l7 7a.75.75 0 0 1 0 1.06l-7 7a.77.77 0 0 1-.53.22Zm.75-12.94v10.38L13.94 12 8.75 6.81Z'
      fill='currentColor'
    />
    <path
      d='M16 19.75a.76.76 0 0 1-.75-.75V5a.75.75 0 1 1 1.5 0v14a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgStepForward;

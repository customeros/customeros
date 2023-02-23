import * as React from 'react';
import { SVGProps } from 'react';
const SvgStepForwardAlt = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M16 19.5a.76.76 0 0 1-.75-.75v-14a.75.75 0 1 1 1.5 0v14a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
    <path
      d='M8 20a.75.75 0 0 1-.29-.06.74.74 0 0 1-.46-.69v-14a.74.74 0 0 1 .46-.69.75.75 0 0 1 .82.16l7 7a.75.75 0 0 1 0 1.06l-7 7A.75.75 0 0 1 8 20Zm.75-12.94v10.38l5.19-5.19-5.19-5.19Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgStepForwardAlt;

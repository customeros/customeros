import * as React from 'react';
import { SVGProps } from 'react';
const SvgPlay = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <path
      d='M8.5 19.75a.76.76 0 0 1-.29-.06.74.74 0 0 1-.46-.69V5a.75.75 0 0 1 .457-.662A.75.75 0 0 1 9 4.47l7 7a.75.75 0 0 1 0 1.06l-7 7a.77.77 0 0 1-.5.22zm.75-12.94v10.38L14.44 12 9.25 6.81z'
      fill='currentColor'
    />
  </svg>
);
export default SvgPlay;

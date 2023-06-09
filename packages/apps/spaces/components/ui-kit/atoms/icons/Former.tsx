import * as React from 'react';
import { SVGProps } from 'react';
const SvgFormer = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 22.5c5.799 0 10.5-4.701 10.5-10.5S17.799 1.5 12 1.5 1.5 6.201 1.5 12 6.201 22.5 12 22.5Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M7.23 7.23h7.63a3.82 3.82 0 0 1 3.82 3.82 3.82 3.82 0 0 1-3.82 3.81H7.23'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m11.04 18.68-3.81-3.82 3.81-3.82'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgFormer;

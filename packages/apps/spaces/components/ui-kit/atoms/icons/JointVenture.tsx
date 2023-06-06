import * as React from 'react';
import { SVGProps } from 'react';
const SvgJointVenture = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M20.59 22.5H3.41L12 11.04l8.59 11.46ZM17.73 1.5H12v4.77h5.73V1.5ZM12 11.05V6.27'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgJointVenture;

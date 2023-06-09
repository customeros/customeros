import * as React from 'react';
import { SVGProps } from 'react';
const SvgProspect = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M7.21 18.69H1.5v3.81h5.71v-3.81ZM12.92 14.88H7.21v7.62h5.71v-7.62Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m15.79 2.5-5.72 7.62h2.86V22.5h5.71V10.12h2.86L15.79 2.5Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgProspect;

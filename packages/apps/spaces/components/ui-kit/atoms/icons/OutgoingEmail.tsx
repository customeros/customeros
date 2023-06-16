import * as React from 'react';
import { SVGProps } from 'react';
const SvgOutgoingEmail = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M.5 8.17h4.78M1.46 12h3.82M2.41 15.83h2.87M5.28 18.7H22.5V5.3H5.28'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m5.28 5.3 8.61 8.42L22.5 5.3M11.98 12l-6.7 6.7M22.5 18.7 15.8 12'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgOutgoingEmail;

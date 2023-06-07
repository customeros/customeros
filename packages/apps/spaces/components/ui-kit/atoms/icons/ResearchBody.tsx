import * as React from 'react';
import { SVGProps } from 'react';
const SvgResearchBody = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M14.86 16.78a7.64 7.64 0 1 0 0-15.28 7.64 7.64 0 0 0 0 15.28ZM1.5 22.5l7.64-7.64'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M14.86 11.04a2.86 2.86 0 1 0 0-5.72 2.86 2.86 0 0 0 0 5.72Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M19.58 15.14a4.76 4.76 0 0 0-9.43 0'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgResearchBody;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgContractManufacturer = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19.63 9.15h-7.65v3.83h1.91l5.74-1.92V9.15Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M13.89 1.5H8.15v18.17h3.83V5.33h7.65V3.41L13.89 1.5ZM3.37 5.33h2.9l1.88-.96M3.37 9.15h2.96l1.82.96M10.07 23.5v-3.83'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgContractManufacturer;

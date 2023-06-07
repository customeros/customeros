import * as React from 'react';
import { SVGProps } from 'react';
const SvgInvestor = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M6.28 22.49h7.56a4.74 4.74 0 0 0 3.37-1.4l4.71-4.71a1.998 1.998 0 0 0 .59-1.45 2 2 0 0 0-3.45-1.41l-5.15 5.15H13a1.9 1.9 0 0 0 1.88-2.22 2 2 0 0 0-2-1.59H9.14l-.93-.46a4.66 4.66 0 0 0-5 .54 4.6 4.6 0 0 0-1.69 3.62v.11M12.95 18.67H8.19M5.33 1.51h1.9A4.77 4.77 0 0 1 12 6.28h-1.91a4.77 4.77 0 0 1-4.76-4.77ZM13.91 7.23H12a4.77 4.77 0 0 1 4.77-4.77h1.91a4.77 4.77 0 0 1-4.77 4.77ZM12 14.86V6.28'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgInvestor;

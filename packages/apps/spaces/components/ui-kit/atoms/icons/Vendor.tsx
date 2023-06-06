import * as React from 'react';
import { SVGProps } from 'react';
const SvgVendor = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17.25 1.5H3.6L1.5 6.75a2.52 2.52 0 0 0 2.62 2.39 2.52 2.52 0 0 0 2.63-2.39 2.52 2.52 0 0 0 2.63 2.39A2.52 2.52 0 0 0 12 6.75a2.52 2.52 0 0 0 2.62 2.39 2.52 2.52 0 0 0 2.63-2.39 2.52 2.52 0 0 0 2.63 2.39 2.52 2.52 0 0 0 2.62-2.39L20.4 1.5h-3.15Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m6.75 7.23.15-.8.9-4.93M17.25 7.23 16.2 1.5M12 1.5v5.73M20.59 9.14V22.5H3.41V9.14M.55 22.5h22.9'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M15.82 13.91H8.18v8.59h7.64v-8.59ZM3.41 17.73h4.77M15.82 17.73h4.77'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgVendor;

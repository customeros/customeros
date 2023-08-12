import * as React from 'react';
import { SVGProps } from 'react';
const SvgCurrencyDollar = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M4 10.667a2.667 2.667 0 0 0 2.667 2.666h2.666a2.667 2.667 0 0 0 0-5.333H6.667a2.667 2.667 0 0 1 0-5.333h2.666A2.667 2.667 0 0 1 12 5.333m-4-4v13.334'
      stroke='currentColor'
      strokeWidth={1.333}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgCurrencyDollar;

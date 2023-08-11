import * as React from 'react';
import { SVGProps } from 'react';
const SvgBillingDetails = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={46}
    height={46}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect x={3} y={3} width={40} height={40} rx={20} fill='#F2F4F7' />
    <g clipPath='url(#billing-details_svg__a)'>
      <path
        d='M26.281 26.281a5.835 5.835 0 1 0-6.563-6.563m6.615 5.782a5.833 5.833 0 1 1-11.667 0 5.833 5.833 0 0 1 11.667 0Z'
        stroke='#667085'
        strokeWidth={1.5}
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
    <rect
      x={3}
      y={3}
      width={40}
      height={40}
      rx={20}
      stroke='#F9FAFB'
      strokeWidth={6}
    />
    <defs>
      <clipPath id='billing-details_svg__a'>
        <path fill='#fff' transform='translate(13 13)' d='M0 0h20v20H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgBillingDetails;

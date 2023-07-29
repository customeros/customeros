import * as React from 'react';
import { SVGProps } from 'react';
const SvgUser = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 28 28'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M23.333 24.5c0-1.628 0-2.442-.2-3.105a4.667 4.667 0 0 0-3.112-3.11c-.662-.202-1.476-.202-3.104-.202h-5.834c-1.628 0-2.442 0-3.104.201a4.667 4.667 0 0 0-3.111 3.111c-.201.663-.201 1.477-.201 3.105M19.25 8.75a5.25 5.25 0 1 1-10.5 0 5.25 5.25 0 0 1 10.5 0Z'
      stroke='currentColor'
      strokeWidth={2.333}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgUser;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgStrikethrough = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M4 10.667a2.667 2.667 0 0 0 2.667 2.666h2.666a2.667 2.667 0 1 0 0-5.333M12 5.333a2.667 2.667 0 0 0-2.667-2.666H6.667A2.667 2.667 0 0 0 4 5.333M2 8h12'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgStrikethrough;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgTag = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M3.41 1.5 1.5 3.41v7.63L12.96 22.5l9.54-9.54L11.04 1.5H3.41Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M6.27 8.18a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82Z'
      stroke='#000'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgTag;

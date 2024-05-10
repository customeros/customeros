import * as React from 'react';
import { SVGProps } from 'react';
const SvgExclamationWaves = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 38 38'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect
      x={6}
      y={6}
      width={26}
      height={26}
      rx={13}
      stroke='#D92D20'
      strokeWidth={2}
      opacity={0.3}
    />
    <rect
      x={1}
      y={1}
      width={36}
      height={36}
      rx={18}
      stroke='#D92D20'
      strokeWidth={2}
      opacity={0.1}
    />
    <g clipPath='url(#exclamationWaves_svg__a)'>
      <path
        d='M19 15.667V19m0 3.333h.008M27.333 19a8.333 8.333 0 1 1-16.666 0 8.333 8.333 0 0 1 16.666 0Z'
        stroke='#D92D20'
        strokeWidth={1.667}
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
    <defs>
      <clipPath id='exclamationWaves_svg__a'>
        <path fill='#fff' transform='translate(9 9)' d='M0 0h20v20H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgExclamationWaves;

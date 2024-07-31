import { SVGProps } from 'react';
const SvgExclamationWaves = (props: SVGProps<SVGSVGElement>) => (
  <svg
    fill='none'
    viewBox='0 0 38 38'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect
      x={6}
      y={6}
      rx={13}
      width={26}
      height={26}
      opacity={0.3}
      strokeWidth={2}
      stroke='#D92D20'
    />
    <rect
      x={1}
      y={1}
      rx={18}
      width={36}
      height={36}
      opacity={0.1}
      strokeWidth={2}
      stroke='#D92D20'
    />
    <g clipPath='url(#exclamationWaves_svg__a)'>
      <path
        stroke='#D92D20'
        strokeWidth={1.667}
        strokeLinecap='round'
        strokeLinejoin='round'
        d='M19 15.667V19m0 3.333h.008M27.333 19a8.333 8.333 0 1 1-16.666 0 8.333 8.333 0 0 1 16.666 0Z'
      />
    </g>
    <defs>
      <clipPath id='exclamationWaves_svg__a'>
        <path fill='#fff' d='M0 0h20v20H0z' transform='translate(9 9)' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgExclamationWaves;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgMicrophoneSlashed = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    fill='none'
    viewBox='0 0 24 24'
    {...props}
  >
    <g
      stroke='#fff'
      strokeLinecap='round'
      strokeLinejoin='round'
      strokeWidth={2}
      clipPath='url(#microphone-slashed_svg__a)'
    >
      <path d='M8 23h8M1 1l22 22M15 9.3V4a3 3 0 0 0-6-.6M9 9v3a3 3 0 0 0 5.1 2.1L9 9Z' />
      <path d='M17 17a7 7 0 0 1-12-5v-2m14 0v2l-.1 1.2M12 19v4' />
    </g>
    <defs>
      <clipPath id='microphone-slashed_svg__a'>
        <path fill='#fff' d='M0 0h24v24H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgMicrophoneSlashed;

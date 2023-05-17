import * as React from 'react';
import { SVGProps } from 'react';
const SvgFastBackward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    xmlnsXlink='http://www.w3.org/1999/xlink'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <g fill='currentColor'>
      <use xlinkHref='#fast-backward_svg__a' />
      <path d='M4 19.5a.76.76 0 0 1-.75-.75v-14A.75.75 0 0 1 4 4a.75.75 0 0 1 .75.75v14a.76.76 0 0 1-.75.75z' />
      <use xlinkHref='#fast-backward_svg__a' x={-8} />
    </g>
    <defs>
      <path
        id='fast-backward_svg__a'
        d='M20 20a.75.75 0 0 1-.53-.22l-7-7a.75.75 0 0 1 0-1.06l7-7a.75.75 0 0 1 .82-.16.74.74 0 0 1 .46.69v14a.74.74 0 0 1-.46.69.75.75 0 0 1-.29.06zm-5.94-7.75 5.19 5.19V7.06l-5.19 5.19z'
      />
    </defs>
  </svg>
);
export default SvgFastBackward;

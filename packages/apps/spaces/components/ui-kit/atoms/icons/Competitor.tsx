import * as React from 'react';
import { SVGProps } from 'react';
const SvgCompetitor = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M7.93 12C8.05 7.37 9.1 1.5 11.9 1.5c2.45 0 2.88 5.15 3 8.58M15.13 15.09c.3 2.14.75 4.22.75 5.5 0 1.28.23 1.91-2 1.91-1.65 0-4.68-3.23-5.65-7.1'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M10.9 13s-3 0-3 4.78-2.09 1.7-3 0c-1-1.91-1-5.73 2-5.73s4 .95 4 .95Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M10.9 13s2.74-1.39 4.66 3 2.85.71 3.07-1.24c.26-2.17-1.28-5.69-4.35-4.46-3.07 1.23-3.38 2.7-3.38 2.7Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path d='M11 13s1 3.82 1.91 4.78' stroke='#000' strokeMiterlimit={10} />
  </svg>
);
export default SvgCompetitor;

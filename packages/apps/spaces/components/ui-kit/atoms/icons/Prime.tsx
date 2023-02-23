import * as React from 'react';
import { SVGProps } from 'react';
const SvgPrime = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='m15.36 11.33-1.1-.25.86 1.22v3.79l2.93-2.44V9.5l-1.35.48-1.34 1.35ZM8.52 11.33 7.17 9.98 5.83 9.5v4.15l2.93 2.44V12.3l.86-1.22-1.1.25Z' />
      <path d='M13.16 11.45h-2.44l-.61-.37-.98 1.47v5.5l.73 1.09.86.86h2.44l.86-.86.73-1.09v-5.5l-.98-1.47-.61.37ZM15.12 18.41l1.58-1.59v-1.58l-1.58 1.34v1.83ZM7.17 16.82l1.59 1.59v-1.83l-1.59-1.34v1.58ZM11.08 10.96h.62V4h-1.35l-.97 2.31-4.41-.36.74 3.06 5.25 1.95h.12ZM14.5 6.31 13.53 4h-1.35v7H13l5.25-2L19 6l-4.5.31Z' />
      <path d='M17.32 5.71 15.6 4h-1.71l.86 1.95 2.57-.24ZM9.98 4H8.27L6.56 5.71l2.57.24L9.98 4Z' />
    </g>
  </svg>
);
export default SvgPrime;

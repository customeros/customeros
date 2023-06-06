import * as React from 'react';
import { SVGProps } from 'react';
const SvgTarget = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M10.09 20.59a8.59 8.59 0 1 0 0-17.18 8.59 8.59 0 0 0 0 17.18Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M10.09 16.77a4.77 4.77 0 1 0 0-9.54 4.77 4.77 0 0 0 0 9.54ZM10.09 12h10.5M22.5 10.09 20.59 12M22.5 13.91 20.59 12'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M10.09 12.95a.95.95 0 1 0 0-1.9.95.95 0 0 0 0 1.9Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgTarget;

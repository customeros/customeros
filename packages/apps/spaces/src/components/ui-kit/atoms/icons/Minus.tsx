import * as React from 'react';
import { SVGProps } from 'react';
const SvgMinus = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path d='M20 13H4a1 1 0 0 1 0-2h16a1 1 0 0 1 0 2Z' fill='currentColor' />
  </svg>
);
export default SvgMinus;

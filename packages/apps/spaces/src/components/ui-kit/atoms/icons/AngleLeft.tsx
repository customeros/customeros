import * as React from 'react';
import { SVGProps } from 'react';
const SvgAngleLeft = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M13.75 16.25a.74.74 0 0 1-.53-.22l-3.5-3.5a.75.75 0 0 1 0-1.06L13.22 8a.75.75 0 0 1 1.06 1l-3 3 3 3a.75.75 0 0 1 0 1.06.74.74 0 0 1-.53.19Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgAngleLeft;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgChevronUp = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17 15.25a.74.74 0 0 1-.53-.22L12 10.56 7.53 15a.75.75 0 0 1-1.06-1l5-5a.75.75 0 0 1 1.06 0l5 5a.75.75 0 0 1 0 1.06.74.74 0 0 1-.53.19Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgChevronUp;

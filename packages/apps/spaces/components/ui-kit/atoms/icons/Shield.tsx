import * as React from 'react';
import { SVGProps } from 'react';
const SvgShield = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 20.75a.87.87 0 0 1-.28-.05A14.27 14.27 0 0 1 3.29 6.43a.74.74 0 0 1 .61-.69 27.12 27.12 0 0 0 7.79-2.42.75.75 0 0 1 .62 0 27.119 27.119 0 0 0 7.79 2.42.74.74 0 0 1 .61.69 14.27 14.27 0 0 1-8.43 14.27.871.871 0 0 1-.28.05ZM4.76 7.11A12.47 12.47 0 0 0 12 19.18a12.47 12.47 0 0 0 7.24-12.07A27.559 27.559 0 0 1 12 4.82a27.56 27.56 0 0 1-7.24 2.29Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgShield;

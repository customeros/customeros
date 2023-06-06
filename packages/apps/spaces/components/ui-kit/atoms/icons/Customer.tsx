import * as React from 'react';
import { SVGProps } from 'react';
const SvgCustomer = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M2.45 20.59h19.1M2.46 7.23l1.9 6.68v6.68h15.28v-6.68l1.91-6.68'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m21.55 7.23-5.73 2.86L12 5.32M12 5.32l-3.82 4.77-5.72-2.86M12 17.72c1.055 0 1.91-1.28 1.91-2.86 0-1.58-.855-2.86-1.91-2.86s-1.91 1.28-1.91 2.86c0 1.58.855 2.86 1.91 2.86ZM15.82 14.86h3.82M8.18 14.86H4.36M12 5.31a.95.95 0 1 0 0-1.9.95.95 0 0 0 0 1.9Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M2.45 7.22a.95.95 0 1 0 0-1.9.95.95 0 0 0 0 1.9ZM21.55 7.22a.95.95 0 1 0 0-1.9.95.95 0 0 0 0 1.9Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgCustomer;

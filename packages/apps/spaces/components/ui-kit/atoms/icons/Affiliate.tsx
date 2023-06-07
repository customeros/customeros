import * as React from 'react';
import { SVGProps } from 'react';
const SvgAffiliate = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 5.34a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M9.14 8.2a2.86 2.86 0 0 1 5.72 0M19.64 20.61a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM16.77 23.48a2.87 2.87 0 0 1 2.87-2.87 2.87 2.87 0 0 1 2.86 2.87M4.36 20.61a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM1.5 23.48a2.87 2.87 0 0 1 2.86-2.87 2.87 2.87 0 0 1 2.87 2.87M12 9.16v4.77M8.18 16.8 12 13.93l3.82 2.87'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgAffiliate;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgLead = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M14.86 20.66a2.86 2.86 0 1 0-5.72 0v2.86M12 15.89a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM19.64 17.8a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM1.5 23.52v-.95a2.87 2.87 0 0 1 2.86-2.87 2.87 2.87 0 0 1 2.87 2.87v.95M16.77 23.52v-.95a2.87 2.87 0 0 1 2.87-2.87 2.87 2.87 0 0 1 2.86 2.87v.95M4.36 17.8a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM14.86 20.66v2.86M12 3.48l.46.94 1.05.15-.76.73.18 1.04-.93-.49-.93.49.18-1.04-.76-.73 1.05-.15.46-.94Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgLead;

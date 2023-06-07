import * as React from 'react';
import { SVGProps } from 'react';
const SvgDistributor = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 20.61a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM9.14 23.48A2.87 2.87 0 0 1 12 20.61a2.87 2.87 0 0 1 2.86 2.87M19.64 20.61a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM16.77 23.48a2.87 2.87 0 0 1 2.87-2.87 2.87 2.87 0 0 1 2.86 2.87M4.36 20.61a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM1.5 23.48a2.87 2.87 0 0 1 2.86-2.87 2.87 2.87 0 0 1 2.87 2.87M4.36 13.93v-3.82h15.28v3.82M12 13.93V7.25M12 7.25a2.86 2.86 0 1 0 0-5.72 2.86 2.86 0 0 0 0 5.72Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgDistributor;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgLost = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M15 6.5h-3.5a1.5 1.5 0 0 0 0 3h1a1.5 1.5 0 0 1 0 3H9M12 14.5v-2M12 6.5v-2'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m3 13.5 9 9 9-9h-3v-12H6v12H3Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgLost;

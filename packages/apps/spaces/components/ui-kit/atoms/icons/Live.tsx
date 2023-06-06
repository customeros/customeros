import * as React from 'react';
import { SVGProps } from 'react';
const SvgLive = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12 11.03a1.91 1.91 0 1 0 0-3.82 1.91 1.91 0 0 0 0 3.82ZM7.28 13.84a6.68 6.68 0 0 1 0-9.45M16.72 4.39a6.68 6.68 0 0 1 0 9.45M4.58 16.54a10.48 10.48 0 0 1 0-14.85M19.42 1.69a10.48 10.48 0 0 1 0 14.85'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M12 11a2.86 2.86 0 0 1 2.86 2.86c.002 2.51-.58 4.984-1.7 7.23l-.21.41h-1.9l-.21-.41a16.18 16.18 0 0 1-1.7-7.23A2.86 2.86 0 0 1 12 11Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgLive;

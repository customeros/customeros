import * as React from 'react';
import { SVGProps } from 'react';
const SvgPartner = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M16.76 6.26a2.39 2.39 0 1 0 0-4.78 2.39 2.39 0 0 0 0 4.78ZM20.59 12a3.84 3.84 0 0 0-3.83-3.83A3.84 3.84 0 0 0 12.93 12M8.15 16.78l3.62.91c.14.015.28.015.42 0a1.7 1.7 0 0 0 1.7-1.7v-.12A1.7 1.7 0 0 0 13 14.4l-3.33-1.64a7.19 7.19 0 0 0-5.5-.39L2.41 13H.5'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M12 17.74h1.91l6.47-1.85a1.68 1.68 0 0 1 2.12 1.61v.11a1.67 1.67 0 0 1-1 1.53L15 22a5.781 5.781 0 0 1-2.32.49 5.788 5.788 0 0 1-1.38-.17l-7-1.74H.5'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgPartner;

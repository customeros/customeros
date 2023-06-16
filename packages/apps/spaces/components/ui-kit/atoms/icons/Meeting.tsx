import * as React from 'react';
import { SVGProps } from 'react';
const SvgMeeting = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M22.5 2.43h-21V7.2h21V2.43ZM22.5 7.2h-21v15.27h21V7.2ZM6.27.52v3.82M17.73.52v3.82M12 .52v3.82'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgMeeting;

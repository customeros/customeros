import * as React from 'react';
import { SVGProps } from 'react';
const SvgCompass = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='m15.94 7.62-4.88 2a2.63 2.63 0 0 0-1.48 1.48l-2 4.88a.34.34 0 0 0 .19.44c.08.03.17.03.25 0l4.88-2a2.632 2.632 0 0 0 1.48-1.48l2-4.88a.34.34 0 0 0-.19-.44.36.36 0 0 0-.25 0ZM12 13a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z' />
      <path d='M12 21a9 9 0 1 1 0-18 9 9 0 0 1 0 18Zm0-16.5a7.5 7.5 0 1 0 0 15 7.5 7.5 0 0 0 0-15Z' />
    </g>
  </svg>
);
export default SvgCompass;

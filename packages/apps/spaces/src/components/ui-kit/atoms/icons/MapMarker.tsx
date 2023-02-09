import * as React from 'react';
import { SVGProps } from 'react';
const SvgMapMarker = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 21.25a.69.69 0 0 1-.41-.13c-.3-.19-7.34-4.92-7.34-10.67a7.75 7.75 0 0 1 15.5 0c0 5.75-7 10.48-7.34 10.67a.69.69 0 0 1-.41.13Zm0-17a6.23 6.23 0 0 0-6.25 6.2c0 4.21 4.79 8.06 6.25 9.13 1.46-1.07 6.25-4.92 6.25-9.13A6.23 6.23 0 0 0 12 4.25Z' />
      <path d='M12 12.75a2.75 2.75 0 1 1 0-5.5 2.75 2.75 0 0 1 0 5.5Zm0-4a1.25 1.25 0 1 0 0 2.5 1.25 1.25 0 0 0 0-2.5Z' />
    </g>
  </svg>
);
export default SvgMapMarker;

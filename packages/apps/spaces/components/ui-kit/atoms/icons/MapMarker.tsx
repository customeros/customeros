import * as React from 'react';
import { SVGProps } from 'react';
const SvgMapMarker = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12 21.25a.69.69 0 0 1-.41-.13c-.3-.19-7.34-4.92-7.34-10.67a7.75 7.75 0 1 1 15.5 0c0 5.75-7 10.48-7.34 10.67a.69.69 0 0 1-.41.13zm0-17a6.23 6.23 0 0 0-6.25 6.2c0 4.21 4.79 8.06 6.25 9.13 1.46-1.07 6.25-4.92 6.25-9.13A6.23 6.23 0 0 0 12 4.25zm0 8.5a2.75 2.75 0 0 1-1.945-4.695A2.75 2.75 0 0 1 14.75 10 2.75 2.75 0 0 1 12 12.75zm0-4a1.25 1.25 0 0 0-.884 2.134A1.25 1.25 0 0 0 13.25 10 1.25 1.25 0 0 0 12 8.75z' />
  </svg>
);
export default SvgMapMarker;

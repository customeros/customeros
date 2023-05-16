import * as React from 'react';
import { SVGProps } from 'react';
const SvgCloud = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M16.5 18.75h-7a6.75 6.75 0 1 1 6.2-9.42c.264-.05.532-.078.8-.08A5.07 5.07 0 0 1 21.25 14a4.75 4.75 0 0 1-4.75 4.75zm-7-12a5.25 5.25 0 0 0 0 10.5h7A3.26 3.26 0 0 0 19.75 14a3.57 3.57 0 0 0-3.25-3.25c-.341.008-.68.072-1 .19a.78.78 0 0 1-.58 0 .731.731 0 0 1-.37-.44A5.24 5.24 0 0 0 9.5 6.75z'
      fill='currentColor'
    />
  </svg>
);
export default SvgCloud;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgArrowDownRight = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17.7 8.46a.75.75 0 1 0-1.5 0v6.68L7.58 6.52a.75.75 0 0 0-1.06 1.06l8.62 8.62H8.46a.75.75 0 1 0 0 1.5H17c.1 0 .198-.02.29-.06a.76.76 0 0 0 .41-.64V8.46Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgArrowDownRight;

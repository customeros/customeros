import * as React from 'react';
import { SVGProps } from 'react';
const SvgMinusCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='M12 21A9 9 0 0 1 5.636 5.636 9 9 0 0 1 21 12a9 9 0 0 1-9 9zm0-16.5a7.5 7.5 0 0 0-5.303 12.803A7.5 7.5 0 0 0 19.5 12 7.5 7.5 0 0 0 12 4.5zm4 8.25H8a.75.75 0 1 1 0-1.5h8a.75.75 0 1 1 0 1.5z' />
  </svg>
);
export default SvgMinusCircle;

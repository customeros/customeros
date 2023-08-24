import * as React from 'react';
import { SVGProps } from 'react';
const SvgUnderline = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M12 2.667v4.666a4 4 0 0 1-8 0V2.667M2.667 14h10.666'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgUnderline;

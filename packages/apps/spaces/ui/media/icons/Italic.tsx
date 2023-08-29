import * as React from 'react';
import { SVGProps } from 'react';
const SvgItalic = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M12.667 2.667h-6m2.666 10.666h-6M10 2.667 6 13.333'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgItalic;

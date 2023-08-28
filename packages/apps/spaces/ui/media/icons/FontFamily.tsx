import * as React from 'react';
import { SVGProps } from 'react';
const SvgFontFamily = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    width='1em'
    height='1em'
    {...props}
  >
    <path
      d='M2.667 4.667c0-.622 0-.932.101-1.177.135-.327.395-.587.722-.722.245-.101.555-.101 1.177-.101h6.666c.622 0 .932 0 1.177.101.327.135.587.395.722.722.101.245.101.555.101 1.177M6 13.333h4M8 2.667v10.666'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgFontFamily;

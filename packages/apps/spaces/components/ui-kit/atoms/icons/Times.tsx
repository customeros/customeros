import * as React from 'react';
import { SVGProps } from 'react';
const SvgTimes = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='m13.06 12 4.42-4.42a.75.75 0 1 0-1.06-1.06L12 10.94 7.58 6.52a.75.75 0 0 0-1.06 1.06L10.94 12l-4.42 4.42a.75.75 0 0 0 1.06 1.06L12 13.06l4.42 4.42a.75.75 0 0 0 1.06-1.06L13.06 12Z'
      fill='#212121'
    />
  </svg>
);
export default SvgTimes;

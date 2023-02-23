import * as React from 'react';
import { SVGProps } from 'react';
const SvgMoon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M12.09 20.66c-.364 0-.729-.024-1.09-.07A8.8 8.8 0 0 1 3.41 13a8.71 8.71 0 0 1 6.83-9.67 1.23 1.23 0 0 1 1.27.48 1.27 1.27 0 0 1 .05 1.4 5.3 5.3 0 0 0-.66 3.47 5.24 5.24 0 0 0 4.38 4.38 5.19 5.19 0 0 0 3.47-.67 1.27 1.27 0 0 1 1.4.07 1.21 1.21 0 0 1 .48 1.26 8.7 8.7 0 0 1-8.54 6.94ZM10 5a7.25 7.25 0 1 0 9 9 6.74 6.74 0 0 1-9.61-5A6.75 6.75 0 0 1 10 5Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgMoon;

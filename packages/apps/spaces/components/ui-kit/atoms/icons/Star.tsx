import * as React from 'react';
import { SVGProps } from 'react';
const SvgStar = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <path
      d='M21.12 9.88a.74.74 0 0 0-.6-.51l-5.42-.79-2.43-4.91a.78.78 0 0 0-.67-.381.78.78 0 0 0-.67.381L8.9 8.58l-5.42.79a.74.74 0 0 0-.6.51.75.75 0 0 0 .18.77L7 14.47l-.93 5.4a.76.76 0 0 0 .3.74.75.75 0 0 0 .79.05L12 18.11l4.85 2.55a.73.73 0 0 0 .35.09.79.79 0 0 0 .44-.14.76.76 0 0 0 .3-.74l-.94-5.4 3.93-3.82a.75.75 0 0 0 .19-.77z'
      fill='currentColor'
    />
  </svg>
);
export default SvgStar;

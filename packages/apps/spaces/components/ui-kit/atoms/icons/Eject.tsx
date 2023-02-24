import * as React from 'react';
import { SVGProps } from 'react';
const SvgEject = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19 12.75H5a.75.75 0 0 1-.53-1.28l7-7a.75.75 0 0 1 1.06 0l7 7a.75.75 0 0 1-.53 1.28Zm-12.19-1.5h10.38L12 6.06l-5.19 5.19ZM18 19.75H6A1.76 1.76 0 0 1 4.25 18v-2A1.76 1.76 0 0 1 6 14.25h12A1.76 1.76 0 0 1 19.75 16v2A1.76 1.76 0 0 1 18 19.75Zm-12-4a.25.25 0 0 0-.25.25v2a.25.25 0 0 0 .25.25h12a.25.25 0 0 0 .25-.25v-2a.25.25 0 0 0-.25-.25H6Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgEject;

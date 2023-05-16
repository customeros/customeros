import * as React from 'react';
import { SVGProps } from 'react';
const SvgVolumeOff = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <path
      d='M17 19.75a.81.81 0 0 1-.47-.16l-4.79-3.84H7a.76.76 0 0 1-.75-.75V9A.76.76 0 0 1 7 8.25h4.74l4.79-3.84a.75.75 0 0 1 .796-.091.75.75 0 0 1 .424.681v14a.77.77 0 0 1-.42.68.78.78 0 0 1-.33.07zm-9.25-5.5H12a.78.78 0 0 1 .47.16l3.78 3V6.56l-3.78 3a.78.78 0 0 1-.47.16H7.75v4.53z'
      fill='currentColor'
    />
  </svg>
);
export default SvgVolumeOff;

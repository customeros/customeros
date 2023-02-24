import * as React from 'react';
import { SVGProps } from 'react';
const SvgVolumeUp = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M13 19.75a.81.81 0 0 1-.47-.16l-4.79-3.84H3a.76.76 0 0 1-.75-.75V9A.76.76 0 0 1 3 8.25h4.74l4.79-3.84a.75.75 0 0 1 1.22.59v14a.76.76 0 0 1-.43.68.71.71 0 0 1-.32.07Zm-9.25-5.5H8a.78.78 0 0 1 .47.16l3.78 3V6.56l-3.78 3a.78.78 0 0 1-.47.19H3.75v4.5ZM18.46 18.07a.76.76 0 0 1-.672-.411.742.742 0 0 1 .112-.829 7.24 7.24 0 0 0 0-9.66.75.75 0 1 1 1.12-1 8.7 8.7 0 0 1 0 11.64.721.721 0 0 1-.56.26Z' />
      <path d='M16.11 15.38a.75.75 0 0 1-.6-1.2 3.6 3.6 0 0 0 0-4.36.75.75 0 0 1 1.2-.9 5.07 5.07 0 0 1 0 6.16.77.77 0 0 1-.6.3Z' />
    </g>
  </svg>
);
export default SvgVolumeUp;

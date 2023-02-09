import * as React from 'react';
import { SVGProps } from 'react';
const SvgDirectionsAlt = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g fill='currentColor'>
      <path d='M12 21.25a1.81 1.81 0 0 1-1.28-.53l-7.44-7.44a1.81 1.81 0 0 1 0-2.56l7.44-7.44a1.86 1.86 0 0 1 2.56 0l7.44 7.44a1.81 1.81 0 0 1 0 2.56l-7.44 7.44a1.81 1.81 0 0 1-1.28.53Zm0-17a.27.27 0 0 0-.21.09l-7.45 7.45a.29.29 0 0 0 0 .42l7.45 7.45a.25.25 0 0 0 .42 0l7.45-7.45a.292.292 0 0 0 .067-.324.291.291 0 0 0-.067-.096l-7.45-7.45a.27.27 0 0 0-.21-.09Z' />
      <path d='M10.73 14.42a.74.74 0 0 1-.53-.22L8 12a.75.75 0 0 1 0-1.06l2.2-2.22a.75.75 0 0 1 1.06 1.06l-1.7 1.68 1.7 1.68a.75.75 0 0 1-.53 1.28Z' />
      <path d='M15.5 15a.76.76 0 0 1-.75-.75v-2H8.5a.75.75 0 1 1 0-1.5h7a.76.76 0 0 1 .75.75v2.79a.76.76 0 0 1-.75.71Z' />
    </g>
  </svg>
);
export default SvgDirectionsAlt;

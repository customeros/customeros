import * as React from 'react';
import { SVGProps } from 'react';
const SvgThumbsUp = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17.25 20.75H5.63a2.38 2.38 0 0 1-2.38-2.36v-5.64a2.38 2.38 0 0 1 2.38-2.37h2l2.73-6.09a1.75 1.75 0 0 1 2.27-.9 3.16 3.16 0 0 1 1.88 2.91v2.47h4a.51.51 0 0 1 .17 0 2.56 2.56 0 0 1 1.58 1 2.3 2.3 0 0 1 .44 1.68l-1.1 7.29a2.38 2.38 0 0 1-2.35 2.01Zm-8.43-1.5h8.43a.87.87 0 0 0 .87-.73l1.12-7.26a.72.72 0 0 0-.16-.56 1.12 1.12 0 0 0-.66-.42h-4.67a.74.74 0 0 1-.75-.76V6.3a1.66 1.66 0 0 0-1-1.53.24.24 0 0 0-.31.13l-2.87 6.39v7.96Zm-3.19-7.37a.87.87 0 0 0-.88.87v5.64a.87.87 0 0 0 .88.86h1.69v-7.37H5.63Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgThumbsUp;

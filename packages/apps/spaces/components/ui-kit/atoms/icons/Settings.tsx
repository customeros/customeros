import * as React from 'react';
import { SVGProps } from 'react';
const SvgSettings = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M20.59 12a8.126 8.126 0 0 0-.15-1.57l2.09-1.2-2.87-5-2.08 1.2a8.65 8.65 0 0 0-2.72-1.56V1.5H9.14v2.41a8.65 8.65 0 0 0-2.72 1.56l-2.08-1.2-2.87 5 2.09 1.2a8.29 8.29 0 0 0 0 3.14l-2.09 1.2 2.87 5 2.08-1.2a8.65 8.65 0 0 0 2.72 1.56v2.33h5.72v-2.41a8.651 8.651 0 0 0 2.72-1.56l2.08 1.2 2.87-5-2.09-1.2a8.126 8.126 0 0 0 .15-1.53Z'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
    <path
      d='M12 15.82a3.82 3.82 0 1 0 0-7.64 3.82 3.82 0 0 0 0 7.64Z'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgSettings;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgOutgoingVoice = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M16.79 6.28a4.78 4.78 0 0 0-4.78-4.78H12a4.78 4.78 0 0 0-4.78 4.78v5.74A4.78 4.78 0 0 0 12 16.8h.01a4.78 4.78 0 0 0 4.78-4.78M15.722 8.861H10'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='m14.146 6.715 2.146 2.146-2.146 2.146M20.61 12a8.608 8.608 0 0 1-14.705 6.101A8.609 8.609 0 0 1 3.39 12M12 23.5v-2.87'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgOutgoingVoice;

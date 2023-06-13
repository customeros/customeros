import * as React from 'react';
import { SVGProps } from 'react';
const SvgContacts = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M9.61 12.03a4.3 4.3 0 1 0 0-8.6 4.3 4.3 0 0 0 0 8.6Z'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
    <path
      d='m1.5 21.57.69-3.46A7.58 7.58 0 0 1 9.61 12 7.56 7.56 0 0 1 17 18.11l.7 3.46M12 11.3a4.3 4.3 0 1 0 0-7.14'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
    <path
      d='m22.5 21.57-.7-3.47a7.551 7.551 0 0 0-9.8-5.69'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgContacts;

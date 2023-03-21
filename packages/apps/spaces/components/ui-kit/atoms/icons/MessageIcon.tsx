import * as React from 'react';
import { SVGProps } from 'react';
const SvgMessageIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={15}
    height={15}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M14.165.837a1.157 1.157 0 0 0-1.174-.28L1.303 4.465a1.15 1.15 0 0 0-.799 1 1.156 1.156 0 0 0 .615 1.122L5.97 8.99l2.396 4.877a1.152 1.152 0 0 0 1.03.633h.08a1.141 1.141 0 0 0 1.007-.8l3.955-11.686a1.131 1.131 0 0 0-.272-1.177ZM1.798 5.57 12 2.16 6.336 7.837 1.798 5.57Zm7.654 7.665-2.27-4.55 5.665-5.677-3.395 10.227Z'
      fill='#fff'
    />
  </svg>
);
export default SvgMessageIcon;

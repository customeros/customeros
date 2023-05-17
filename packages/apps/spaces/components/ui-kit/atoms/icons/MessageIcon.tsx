import * as React from 'react';
import { SVGProps } from 'react';
const SvgMessageIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    width={15}
    height={15}
    fill='none'
    {...props}
  >
    <path
      d='M14.165.837a1.149 1.149 0 0 0-1.174-.281L1.303 4.465a1.15 1.15 0 0 0-.55.378 1.158 1.158 0 0 0-.11 1.276 1.15 1.15 0 0 0 .477.468L5.968 8.99l2.397 4.877a1.15 1.15 0 0 0 1.031.633h.08a1.14 1.14 0 0 0 .626-.245c.179-.143.312-.336.38-.556l3.954-11.685a1.13 1.13 0 0 0-.272-1.177zM1.798 5.57 12 2.159 6.336 7.837 1.798 5.57zm7.653 7.664L7.183 8.686l5.664-5.678-3.395 10.227z'
      fill='#fff'
    />
  </svg>
);
export default SvgMessageIcon;

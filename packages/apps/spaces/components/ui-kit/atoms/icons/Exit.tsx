import * as React from 'react';
import { SVGProps } from 'react';
const SvgExit = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M11.5 22.5h-10v-19l10-2v21ZM15.5 3.5h-4v19h4v-19ZM19.5 9.5l2.98 2.97-2.98 2.98M15.5 12.5h7'
      stroke='#343A40'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgExit;

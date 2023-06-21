import * as React from 'react';
import { SVGProps } from 'react';
const SvgNotes = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='m7.23 20.59-4.78 1 1-4.78L17.89 2.29a2.689 2.689 0 0 1 1.91-.79 2.7 2.7 0 0 1 1.91 4.61L7.23 20.59ZM.55 22.5h22.9M19.64 8.18l-3.82-3.82'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgNotes;

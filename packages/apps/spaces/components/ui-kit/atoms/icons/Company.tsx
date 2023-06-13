import * as React from 'react';
import { SVGProps } from 'react';
const SvgCompany = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <g clipPath='url(#company_svg__a)' stroke='#343A40' strokeMiterlimit={10}>
      <path d='M11.04 14.88h1.92M11.04 11.04h1.92M11.04 7.21h1.92M7.21 14.88h1.92M7.21 11.04h1.92M7.21 7.21h1.92M14.88 14.88h1.91M14.88 11.04h1.91M14.88 7.21h1.91' />
      <path d='M14.88 18.71H9.13v3.83h5.75v-3.83Z' strokeLinecap='square' />
      <path d='M23.5 22.54H.5' />
      <path
        d='M16.79 3.38V1.46H7.21v1.92H4.33v19.16h15.34V3.38h-2.88Z'
        strokeLinecap='square'
      />
    </g>
    <defs>
      <clipPath id='company_svg__a'>
        <path fill='#fff' d='M0 0h24v24H0z' />
      </clipPath>
    </defs>
  </svg>
);
export default SvgCompany;

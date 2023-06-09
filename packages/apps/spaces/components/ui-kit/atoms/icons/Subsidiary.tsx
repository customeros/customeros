import * as React from 'react';
import { SVGProps } from 'react';
const SvgSubsidiary = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M17.72 7.23H6.27V22.5h11.45V7.23ZM22.5 14.86h-4.77v7.64h4.77v-7.64ZM6.27 14.86H1.5v7.64h4.77v-7.64Z'
      stroke='#000'
      strokeMiterlimit={10}
    />
    <path
      d='M22.5 11.05h-4.77v3.82h4.77v-3.82ZM6.27 11.05H1.5v3.82h4.77v-3.82ZM10.09 15.82v6.68M13.91 15.82v6.68M18.68 7.23H5.32V4.36L12 1.5l6.68 2.86v2.87ZM3.41 18.68h2.86M17.73 18.68h2.86M10.09 10.09V12M13.91 10.09V12'
      stroke='#000'
      strokeMiterlimit={10}
    />
  </svg>
);
export default SvgSubsidiary;

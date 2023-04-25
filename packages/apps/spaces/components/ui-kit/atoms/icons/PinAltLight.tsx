import * as React from 'react';
import { SVGProps } from 'react';
const SvgPinAltLight = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={17}
    height={18}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M13.713 12.541c.529.323.807.69.807 1.063s-.278.74-.806 1.063c-.529.322-1.289.59-2.204.777a15.34 15.34 0 0 1-3.01.285c-1.058 0-2.096-.098-3.011-.285-.915-.186-1.675-.455-2.204-.778-.528-.322-.806-.689-.806-1.062 0-.373.278-.74.806-1.063'
      stroke='#878787'
      strokeLinecap='round'
    />
    <path
      d='M13.813 7.583c0 3.435-3.63 5.789-4.889 6.502a.852.852 0 0 1-.848 0c-1.26-.713-4.889-3.067-4.889-6.502 0-3.187 2.575-5.312 5.313-5.312 2.833 0 5.313 2.125 5.313 5.313Z'
      stroke='#878787'
    />
    <circle cx={8.499} cy={7.583} r={2.333} stroke='#878787' />
  </svg>
);
export default SvgPinAltLight;

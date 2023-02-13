import * as React from 'react';
import { SVGProps } from 'react';
const SvgStepBackward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M16 19.75a.77.77 0 0 1-.53-.22l-7-7a.75.75 0 0 1 0-1.06l7-7a.75.75 0 0 1 .82-.16.74.74 0 0 1 .46.69v14a.74.74 0 0 1-.46.69.752.752 0 0 1-.29.06ZM10.06 12l5.19 5.19V6.81L10.06 12Z'
      fill='currentColor'
    />
    <path
      d='M8 19.75a.76.76 0 0 1-.75-.75V5a.75.75 0 0 1 1.5 0v14a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgStepBackward;

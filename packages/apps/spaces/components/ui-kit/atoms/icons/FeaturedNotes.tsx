import * as React from 'react';
import { SVGProps } from 'react';
const SvgFeaturedNotes = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 46 46'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect x={3} y={3} width={40} height={40} rx={20} fill='#F2F4F7' />
    <path
      d='M24.667 22.166h-5m1.666 3.334h-1.666m6.666-6.667h-6.666m10-.166v8.666c0 1.4 0 2.1-.273 2.635a2.5 2.5 0 0 1-1.092 1.093c-.535.272-1.235.272-2.635.272h-5.334c-1.4 0-2.1 0-2.634-.272a2.5 2.5 0 0 1-1.093-1.093c-.273-.535-.273-1.235-.273-2.635v-8.666c0-1.4 0-2.1.273-2.636a2.5 2.5 0 0 1 1.092-1.092c.535-.273 1.235-.273 2.636-.273h5.333c1.4 0 2.1 0 2.635.273a2.5 2.5 0 0 1 1.092 1.093c.273.534.273 1.234.273 2.634Z'
      stroke='#667085'
      strokeWidth={1.5}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
    <rect
      x={3}
      y={3}
      width={40}
      height={40}
      rx={20}
      stroke='#F9FAFB'
      strokeWidth={6}
    />
  </svg>
);
export default SvgFeaturedNotes;

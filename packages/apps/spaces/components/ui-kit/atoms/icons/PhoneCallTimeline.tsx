import * as React from 'react';
import { SVGProps } from 'react';
const SvgPhoneCallTimeline = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 136 50'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M134.742 49.44H.5V13.543h134.999l-.047 35.127v.001c0 .446-.339.767-.71.767Z'
      fill='#fff'
      stroke='#88D894'
    />
    <path
      d='M14.5 8A5.5 5.5 0 0 1 20 2.5h110a5.5 5.5 0 0 1 5.5 5.5v6.543h-121V8Z'
      fill='#E4ECE5'
      stroke='#88D894'
    />
    <path d='M0 8a6 6 0 0 1 6-6h57a6 6 0 0 1 6 6v7.043H0V8Z' fill='#88D894' />
  </svg>
);
export default SvgPhoneCallTimeline;

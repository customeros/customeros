import * as React from 'react';
import { SVGProps } from 'react';
const SvgFeaturedCertificate = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 46 46'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <rect x={3} y={3} width={40} height={40} rx={20} fill='#F2F4F7' />
    <path
      d='M20.5 28.416h5M18.833 25.5h8.334m-10-10.834h11.666c.92 0 1.667.83 1.667 1.852v12.963c0 1.023-.746 1.852-1.667 1.852H17.167c-.92 0-1.667-.829-1.667-1.852V16.518c0-1.022.746-1.852 1.667-1.852Zm5.831 3.51c-.583-.649-1.556-.823-2.286-.229-.73.595-.834 1.589-.26 2.292.574.703 2.546 2.344 2.546 2.344s1.972-1.641 2.546-2.344c.574-.703.483-1.704-.26-2.292-.743-.588-1.703-.42-2.286.23Z'
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
export default SvgFeaturedCertificate;

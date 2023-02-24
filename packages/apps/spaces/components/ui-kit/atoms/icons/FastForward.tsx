import * as React from 'react';
import { SVGProps } from 'react';
const SvgFastForward = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M4 19.75c-.1 0-.198-.02-.29-.06a.74.74 0 0 1-.46-.69V5a.75.75 0 0 1 1.28-.53l7 7a.75.75 0 0 1 0 1.06l-7 7a.75.75 0 0 1-.53.22Zm.75-12.94v10.38L9.94 12 4.75 6.81ZM20 19.25a.76.76 0 0 1-.75-.75v-14a.75.75 0 1 1 1.5 0v14a.76.76 0 0 1-.75.75Z'
      fill='currentColor'
    />
    <path
      d='M12 19.75c-.1 0-.198-.02-.29-.06a.74.74 0 0 1-.46-.69V5a.75.75 0 0 1 1.28-.53l7 7a.75.75 0 0 1 0 1.06l-7 7a.75.75 0 0 1-.53.22Zm.75-12.94v10.38L17.94 12l-5.19-5.19Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgFastForward;

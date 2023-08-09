import * as React from 'react';
import { SVGProps } from 'react';
const SvgDelete = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 12 12'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='m8.5 4.5-3 3m0-3 3 3M1.36 6.48l2.16 2.88c.176.235.264.352.376.437a1 1 0 0 0 .33.165c.134.038.28.038.574.038h3.8c.84 0 1.26 0 1.581-.163a1.5 1.5 0 0 0 .655-.656C11 8.861 11 8.441 11 7.6V4.4c0-.84 0-1.26-.164-1.581a1.5 1.5 0 0 0-.655-.656C9.861 2 9.441 2 8.6 2H4.8c-.293 0-.44 0-.575.038a1 1 0 0 0-.33.165c-.111.085-.199.202-.375.437L1.36 5.52c-.13.172-.194.258-.219.353a.5.5 0 0 0 0 .254c.025.095.09.18.219.353Z'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgDelete;

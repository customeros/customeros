import * as React from 'react';
import { SVGProps } from 'react';
const SvgHashtag = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='none'
    {...props}
  >
    <path
      d='M21 7.25h-3l.77-3.07a.75.75 0 0 0-.55-.91.75.75 0 0 0-.91.55l-.86 3.43H10l.77-3.07a.75.75 0 0 0-.55-.91.75.75 0 0 0-.91.55l-.9 3.43H5a.75.75 0 0 0-.75.75.75.75 0 0 0 .75.75h3l-1.63 6.5H3a.75.75 0 0 0-.75.75.75.75 0 0 0 .75.75h3l-.77 3.07a.75.75 0 0 0 .55.91.75.75 0 0 0 .91-.55l.86-3.43H14l-.77 3.07a.75.75 0 0 0 .55.91.75.75 0 0 0 .91-.55l.86-3.43H19a.75.75 0 0 0 .75-.75.75.75 0 0 0-.75-.75h-3l1.63-6.5H21a.75.75 0 0 0 .75-.75.75.75 0 0 0-.75-.75zm-5 1.5-1.63 6.5H8l1.63-6.5H16z'
      fill='currentColor'
    />
  </svg>
);
export default SvgHashtag;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgBriefcase = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19 6.25h-3.75V5a1.89 1.89 0 0 0-2-1.75h-2.5a1.89 1.89 0 0 0-2 1.75v1.25H5A1.76 1.76 0 0 0 3.25 8v10A1.76 1.76 0 0 0 5 19.75h14A1.76 1.76 0 0 0 20.75 18V8A1.76 1.76 0 0 0 19 6.25ZM10.25 5c0-.08.19-.25.5-.25h2.5c.31 0 .5.17.5.25v1.25h-3.5V5ZM5 7.75h14a.25.25 0 0 1 .25.25v3.25H4.75V8A.25.25 0 0 1 5 7.75Zm3.75 5h6.5v1.5h-6.5v-1.5ZM19 18.25H5a.25.25 0 0 1-.25-.25v-5.25h2.5V15a.76.76 0 0 0 .75.75h8a.76.76 0 0 0 .75-.75v-2.25h2.5V18a.25.25 0 0 1-.25.25Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgBriefcase;

import * as React from 'react';
import { SVGProps } from 'react';
const SvgMap = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width={24}
    height={24}
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M19.9 4.09a1.75 1.75 0 0 0-1.66-.15l-3.57 1.53-4.75-2a1.51 1.51 0 0 0-1.18 0L4.38 5.31a1.88 1.88 0 0 0-1.13 1.75v11.25a1.91 1.91 0 0 0 .85 1.6 1.75 1.75 0 0 0 1.66.15l3.57-1.53 4.75 2a1.518 1.518 0 0 0 1.18 0l4.36-1.87a1.88 1.88 0 0 0 1.13-1.75V5.69a1.911 1.911 0 0 0-.85-1.6Zm-9.82 1 3.84 1.64v12.13l-3.84-1.64V5.09ZM5.17 18.68a.25.25 0 0 1-.25 0 .4.4 0 0 1-.17-.35V7.06A.39.39 0 0 1 5 6.69l3.58-1.55v12.08l-3.41 1.46Zm14.08-1.74a.39.39 0 0 1-.22.37l-3.61 1.55V6.78l3.41-1.46a.25.25 0 0 1 .25 0 .4.4 0 0 1 .17.35v11.27Z'
      fill='currentColor'
    />
  </svg>
);
export default SvgMap;

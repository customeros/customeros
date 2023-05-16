import * as React from 'react';
import { SVGProps } from 'react';
const SvgPrime = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns='http://www.w3.org/2000/svg'
    viewBox='0 0 24 24'
    fill='currentColor'
    {...props}
  >
    <path d='m15.36 11.33-1.1-.25.86 1.22v3.79l2.93-2.44V9.5l-1.35.48-1.34 1.35zm-6.84 0L7.17 9.98 5.83 9.5v4.15l2.93 2.44V12.3l.86-1.22-1.1.25zm4.64.12h-2.44l-.61-.37-.98 1.47v5.5l.73 1.09.86.86h2.44l.86-.86.73-1.09v-5.5l-.98-1.47-.61.37zm1.96 6.96 1.58-1.59v-1.58l-1.58 1.34v1.83zm-7.95-1.59 1.59 1.59v-1.83l-1.59-1.34v1.58zm3.91-5.86h.62V4h-1.35l-.97 2.31-4.41-.36.74 3.06 5.25 1.95h.12zm3.42-4.65L13.53 4h-1.35v7H13l5.25-2L19 6l-4.5.31zm2.82-.6L15.6 4h-1.71l.86 1.95 2.57-.24zM9.98 4H8.27L6.56 5.71l2.57.24L9.98 4z' />
  </svg>
);
export default SvgPrime;

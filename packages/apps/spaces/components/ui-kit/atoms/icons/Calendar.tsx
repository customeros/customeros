import * as React from 'react';
import { SVGProps } from 'react';
const SvgCalendar = (props: SVGProps<SVGSVGElement>) => (
  <svg
    viewBox='0 0 16 16'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
    {...props}
  >
    <path
      d='M14 6.667H2m8.667-5.333V4M5.333 1.333V4M5.2 14.667h5.6c1.12 0 1.68 0 2.108-.218a2 2 0 0 0 .874-.874c.218-.428.218-.988.218-2.108v-5.6c0-1.12 0-1.68-.218-2.108a2 2 0 0 0-.874-.874c-.428-.218-.988-.218-2.108-.218H5.2c-1.12 0-1.68 0-2.108.218a2 2 0 0 0-.874.874C2 4.187 2 4.747 2 5.867v5.6c0 1.12 0 1.68.218 2.108a2 2 0 0 0 .874.874c.428.218.988.218 2.108.218Z'
      stroke='currentColor'
      strokeWidth={1.333}
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </svg>
);
export default SvgCalendar;

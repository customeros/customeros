import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Linkedin = ({ className, ...props }: IconProps) => (
  <svg
    width='32'
    height='32'
    fill='none'
    viewBox='0 0 32 32'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <rect x='2' y='2' rx='14' width='28' height='28' fill='#1275B1' />
    <path
      fill='white'
      d='M12.6186 9.69215C12.6186 10.6267 11.8085 11.3843 10.8093 11.3843C9.81004 11.3843 9 10.6267 9 9.69215C9 8.7576 9.81004 8 10.8093 8C11.8085 8 12.6186 8.7576 12.6186 9.69215Z'
    />
    <path fill='white' d='M9.24742 12.6281H12.3402V22H9.24742V12.6281Z' />
    <path
      fill='white'
      d='M17.3196 12.6281H14.2268V22H17.3196C17.3196 22 17.3196 19.0496 17.3196 17.2049C17.3196 16.0976 17.6977 14.9855 19.2062 14.9855C20.911 14.9855 20.9008 16.4345 20.8928 17.5571C20.8824 19.0244 20.9072 20.5219 20.9072 22H24V17.0537C23.9738 13.8954 23.1508 12.4401 20.4433 12.4401C18.8354 12.4401 17.8387 13.1701 17.3196 13.8305V12.6281Z'
    />
  </svg>
);

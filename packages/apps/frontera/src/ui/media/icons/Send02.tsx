import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}

export const Send02 = ({ className, ...props }: IconProps) => (
  <svg
    fill='none'
    viewBox='0 0 24 24'
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <path
      strokeWidth='2'
      stroke='currentColor'
      strokeLinecap='round'
      strokeLinejoin='round'
      d='M12.0005 18.9998V11.9998M12.292 19.0845L19.2704 21.4207C19.8173 21.6038 20.0908 21.6953 20.2594 21.6296C20.4059 21.5726 20.517 21.45 20.5594 21.2987C20.6082 21.1244 20.4903 20.8613 20.2545 20.3349L12.766 3.6222C12.5354 3.1075 12.4201 2.85015 12.2594 2.77041C12.1199 2.70113 11.956 2.70087 11.8162 2.7697C11.6553 2.84892 11.5392 3.1059 11.3069 3.61986L3.75244 20.3359C3.51474 20.8619 3.39589 21.1248 3.44422 21.2993C3.48619 21.4508 3.59697 21.5737 3.74329 21.6312C3.91178 21.6974 4.18567 21.6064 4.73346 21.4246L11.786 19.0838C11.8799 19.0527 11.9268 19.0371 11.9749 19.0309C12.0175 19.0255 12.0606 19.0255 12.1032 19.0311C12.1512 19.0374 12.1981 19.0531 12.292 19.0845Z'
    />
  </svg>
);

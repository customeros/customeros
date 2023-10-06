import { Icon, IconProps } from '@ui/media/Icon';

export const Plus = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='plus'>
      <path
        id='Icon'
        d='M12 5V19M5 12H19'
        stroke='currentColor'
        strokeWidth='2'
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
  </Icon>
);

import { Icon, IconProps } from '@ui/media/Icon';

export const ArrowDown = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='arrow-down'>
      <path
        id='Icon'
        d='M12 5V19M12 19L19 12M12 19L5 12'
        stroke='currentColor'
        strokeWidth='2'
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
  </Icon>
);

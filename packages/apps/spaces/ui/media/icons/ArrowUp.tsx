import { Icon, IconProps } from '@ui/media/Icon';

export const ArrowUp = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='arrow-up'>
      <path
        id='Icon'
        d='M12 19V5M12 5L5 12M12 5L19 12'
        stroke='currentColor'
        strokeWidth='2'
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
  </Icon>
);

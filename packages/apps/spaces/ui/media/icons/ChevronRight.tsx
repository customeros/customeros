import { Icon, IconProps } from '@ui/media/Icon';

export const ChevronRight = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='chevron-right'>
      <path
        id='Icon'
        d='M9 18L15 12L9 6'
        stroke='currentColor'
        strokeWidth='2'
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
  </Icon>
);

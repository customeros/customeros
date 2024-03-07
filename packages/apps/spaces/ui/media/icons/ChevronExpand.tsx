import { Icon, IconProps } from '@ui/media/Icon';

export const ChevronExpand = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='chevron-expand'>
      <path
        id='Icon'
        d='M7 15L12 20L17 15M7 9L12 4L17 9'
        stroke='currentColor'
        strokeWidth='2'
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
  </Icon>
);

import { Icon, IconProps } from '@ui/media/Icon';

export const ChevronCollapse = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <g id='chevron-collapse'>
      <path
        id='Icon'
        d='M7 19L12 14L17 19M7 5L12 10L17 5'
        stroke='currentColor'
        strokeWidth='2'
        strokeLinecap='round'
        strokeLinejoin='round'
      />
    </g>
  </Icon>
);

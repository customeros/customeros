import { Icon, IconProps } from '@ui/media/Icon';

export const ChevronUp = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <path
      d='M18 15L12 9L6 15'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </Icon>
);

import { Icon, IconProps } from '@ui/media/Icon';

export const XCircle = (props: IconProps) => (
  <Icon viewBox='0 0 24 24' fill='none' boxSize='4' {...props}>
    <path
      d='M15 9L9 15M9 9L15 15M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12Z'
      stroke='currentColor'
      strokeWidth='2'
      strokeLinecap='round'
      strokeLinejoin='round'
    />
  </Icon>
);

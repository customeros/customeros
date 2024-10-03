import { cn } from '@ui/utils/cn.ts';
import { IconButton, IconButtonProps } from '@ui/form/IconButton';

interface FloatingToolbarButtonProps extends IconButtonProps {
  active?: boolean;
}

export const FloatingToolbarButton = ({
  onClick,
  active,
  icon,
  ...rest
}: FloatingToolbarButtonProps) => {
  return (
    <IconButton
      {...rest}
      size='xs'
      icon={icon}
      variant='ghost'
      onClick={onClick}
      className={cn(
        'rounded-sm text-gray-100 hover:text-inherit focus:text-inherit hover:bg-gray-600 focus:bg-gray-600 focus:text-gray-100 hover:text-gray-100',
        {
          'bg-gray-600 text-gray-100': active,
        },
      )}
    />
  );
};

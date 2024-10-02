import { ReactElement } from 'react';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';

interface FloatingToolbarButtonProps {
  label: string;
  active?: boolean;
  icon: ReactElement;
  onClick: () => void;
}

export const FloatingToolbarButton = ({
  label,
  onClick,
  active,
  icon,
}: FloatingToolbarButtonProps) => {
  return (
    <IconButton
      size='xs'
      icon={icon}
      variant='ghost'
      onClick={onClick}
      aria-label={label}
      className={cn(
        'rounded-sm text-gray-25 hover:text-inherit focus:text-inherit hover:bg-gray-600 focus:bg-gray-600 focus:text-white hover:text-white',
        {
          'bg-gray-600 text-gray-25': active,
        },
      )}
    />
  );
};

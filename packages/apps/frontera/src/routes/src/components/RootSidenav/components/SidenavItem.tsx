import { ReactElement, MouseEventHandler } from 'react';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';

interface SidenavItemProps {
  href?: string;
  label: string;
  isActive?: boolean;
  onClick?: () => void;
  icon: ((isActive: boolean) => ReactElement) | ReactElement;
}

export const SidenavItem = ({
  label,
  icon,
  onClick,
  isActive,
}: SidenavItemProps) => {
  const handleClick: MouseEventHandler = (e) => {
    e.preventDefault();
    onClick?.();
  };

  const dynamicClasses = cn(
    isActive
      ? ['font-semibold', 'bg-gray-100']
      : ['font-normal', 'bg-transparent'],
  );

  return (
    <Button
      size='md'
      variant='ghost'
      colorScheme='gray'
      onClick={handleClick}
      leftIcon={typeof icon === 'function' ? icon(!!isActive) : icon}
      className={`w-full justify-start px-3 text-gray-700 focus:shadow-sidenavItemFocus ${dynamicClasses}`}
    >
      {label}
    </Button>
  );
};

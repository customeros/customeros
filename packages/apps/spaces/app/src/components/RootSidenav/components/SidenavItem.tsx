'use client';
import { ReactElement, MouseEventHandler } from 'react';

import { Button } from '@ui/form/Button';

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

  return (
    <Button
      px='3'
      w='full'
      size='md'
      variant='ghost'
      fontSize='sm'
      textDecoration='none'
      fontWeight={isActive ? 'semibold' : 'regular'}
      justifyContent='flex-start'
      borderRadius='md'
      bg={isActive ? 'gray.100' : 'transparent'}
      color={isActive ? 'gray.700' : 'gray.500'}
      onClick={handleClick}
      leftIcon={typeof icon === 'function' ? icon(!!isActive) : icon}
    >
      {label}
    </Button>
  );
};

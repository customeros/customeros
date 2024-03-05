'use client';
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
      ? ['font-semibold', 'bg-gray-100', 'text-gray-700']
      : ['font-normal', 'bg-transparent', 'text-gray-500'],
  );

  return (
    <Button
      size='md'
      variant='tertiary'
      // px='3'
      // w='full'
      // size='md'
      // variant='ghost'
      // fontSize='sm'
      // textDecoration='none'
      // fontWeight={isActive ? 'semibold' : 'regular'}
      // justifyContent='flex-start'
      // borderRadius='md'
      // bg={isActive ? 'gray.100' : 'transparent'}
      // color={isActive ? 'gray.700' : 'gray.500'}
      onClick={handleClick}
      leftIcon={typeof icon === 'function' ? icon(!!isActive) : icon}
      // _focus={{
      //   boxShadow: 'sidenavItemFocus',
      // }}
      className={`w-full justify-start px-3 ${dynamicClasses} focus:shadow-sidenavItemFocus`}
    >
      {label}
    </Button>
  );
};

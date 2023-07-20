'use client';
import { ReactElement, useEffect, MouseEventHandler } from 'react';
import { usePathname, useRouter } from 'next/navigation';

import { Button } from '@ui/form/Button';

interface SidenavItemProps {
  href?: string;
  label: string;
  isActive?: boolean;
  icon: ((isActive: boolean) => ReactElement) | ReactElement;
  onClick?: () => void;
}

export const SidenavItem = ({
  label,
  icon,
  href,
  onClick,
  isActive,
}: SidenavItemProps) => {
  const router = useRouter();
  const pathname = usePathname();
  const _isActive = isActive
    ? isActive
    : href
    ? pathname?.startsWith(href)
    : false;

  const handleClick: MouseEventHandler = (e) => {
    e.preventDefault();
    onClick?.();
    if (href) {
      router.push(href);
    }
  };

  useEffect(() => {
    if (href) {
      router.prefetch(href);
    }
  }, [href]);

  return (
    <Button
      px='3'
      w='full'
      as={href ? 'a' : 'button'}
      // href={href}
      size='lg'
      variant='ghost'
      fontSize='md'
      textDecoration='none'
      fontWeight={_isActive ? 'bold' : 'normal'}
      justifyContent='flex-start'
      borderRadius='xl'
      border='3px solid transparent'
      borderColor={_isActive ? 'gray.200' : 'transparent'}
      color='gray.700'
      onClick={handleClick}
      leftIcon={typeof icon === 'function' ? icon(!!_isActive) : icon}
    >
      {label}
    </Button>
  );
};

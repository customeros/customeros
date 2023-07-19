'use client';
import { ReactElement, useEffect, MouseEventHandler } from 'react';
import { usePathname, useRouter } from 'next/navigation';

import { Button } from '@ui/form/Button';

interface SidebarItemProps {
  href?: string;
  label: string;
  icon: (isActive: boolean) => ReactElement;
  onClick?: () => void;
}

export const SidebarItem = ({
  label,
  icon,
  href,
  onClick,
}: SidebarItemProps) => {
  const router = useRouter();
  const pathname = usePathname();
  const isActive = href ? pathname?.startsWith(href) : false;

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
      as='a'
      href={href}
      size='lg'
      variant='ghost'
      fontSize='md'
      textDecoration='none'
      fontWeight={isActive ? 'bold' : 'normal'}
      justifyContent='flex-start'
      borderRadius='xl'
      border='3px solid transparent'
      borderColor={isActive ? 'gray.200' : 'transparent'}
      color='gray.700'
      leftIcon={icon(!!isActive)}
      onClick={handleClick}
    >
      {label}
    </Button>
  );
};

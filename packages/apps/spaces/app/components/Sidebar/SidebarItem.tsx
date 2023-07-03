'use client';
import { ReactNode } from 'react';
import { usePathname } from 'next/navigation';

import { Tooltip } from '@spaces/atoms/tooltip';
import { Flex } from '@ui/layout/Flex';
import { Link } from '@ui/navigation/Link';

interface SidebarItemProps {
  href?: string;
  label: string;
  icon?: ReactNode;
  onClick?: () => void;
}

export const SidebarItem = ({
  label,
  icon,
  href,
  onClick,
}: SidebarItemProps) => {
  const pathname = usePathname();
  const isActive = href ? pathname?.startsWith(href) : false;

  return (
    <Flex
      tabIndex={0}
      role='button'
      onClick={onClick}
      flexDir='column'
      justify='center'
      py='3'
      px='0'
      cursor='pointer'
      transition='all 0.5s ease'
      w='full'
      h='50px'
      bg={isActive ? '#e4e4e4' : 'transparent'}
      _hover={{
        bg: '#e4e4e4',
      }}
    >
      <Link
        href={href ?? ''}
        display='flex'
        justifyContent='center'
        id={`icon-${label}`}
      >
        {icon}
      </Link>
      <Tooltip
        content={label}
        showDelay={300}
        autoHide={false}
        position='right'
        target={`#icon-${label}`}
      />
    </Flex>
  );
};

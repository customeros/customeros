import { ReactElement, MouseEventHandler } from 'react';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface SidenavItemProps {
  href?: string;
  label: string;
  countTag?: number;
  dataTest?: string;
  isActive?: boolean;
  onClick?: () => void;
  icon: ((isActive: boolean) => ReactElement) | ReactElement;
}

export const SidenavItem = ({
  label,
  icon,
  onClick,
  isActive,
  dataTest: dataTest,
  countTag,
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
      data-test={dataTest}
    >
      <div className='w-full flex justify-between '>
        <div>{label}</div>
        {countTag !== undefined && countTag > 0 && (
          <Tag colorScheme='gray' size='sm' variant='solid'>
            <TagLabel>{countTag}</TagLabel>
          </Tag>
        )}
      </div>
    </Button>
  );
};

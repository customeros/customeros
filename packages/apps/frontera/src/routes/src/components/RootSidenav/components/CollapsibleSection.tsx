import { ReactNode } from 'react';

import { cn } from '@ui/utils/cn.ts';
import { Icon } from '@ui/media/Icon';

interface CollapsibleSectionProps {
  title: string;
  isOpen: boolean;
  children: ReactNode;
  onToggle: () => void;
}

export const CollapsibleSection = ({
  title,
  isOpen,
  onToggle,
  children,
}: CollapsibleSectionProps) => {
  return (
    <div>
      <button
        onClick={onToggle}
        className='w-full gap-1 flex justify-flex-start items-center pl-3.5 cursor-pointer text-gray-500 hover:text-gray-700 transition-colors'
      >
        <span className='text-sm  text-gray-500'>{title}</span>

        <Icon
          name='chevron-down'
          className={cn('w-3 h-3', {
            'transform -rotate-90': !isOpen,
          })}
        />
      </button>
      {isOpen && <div className='mt-1'>{children}</div>}
    </div>
  );
};

import React from 'react';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { CommandKbd } from '@ui/overlay/CommandMenu';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ActionItemProps {
  tooltip?: string;
  dataTest?: string;
  onClick: () => void;
  shortcutKey?: string;
  icon: React.ReactElement;
  children: React.ReactNode;
}

export const ActionItem = ({
  icon,
  onClick,
  dataTest,
  tooltip,
  shortcutKey,
  children,
}: ActionItemProps) => {
  return (
    <Tooltip
      className=''
      label={
        tooltip ? (
          <div className='flex items-center text-sm'>
            {tooltip}
            <span
              className={cn(
                'bg-gray-600 text-xs rounded-sm leading-[1.125rem]',
                shortcutKey && 'pl-2',
              )}
            >
              {shortcutKey}
            </span>
          </div>
        ) : (
          <>
            <div className='flex items-center text-sm'>
              Open command menu
              <CommandKbd className='bg-gray-600 text-gray-25 mx-1' />
              <div className='bg-gray-600 text-xs h-5 w-5 rounded-sm flex justify-center items-center'>
                K
              </div>
            </div>
          </>
        )
      }
    >
      <Button
        leftIcon={icon}
        onClick={onClick}
        colorScheme='gray'
        data-test={dataTest}
        className='bg-gray-700 text-gray-25 hover:bg-gray-800 hover:text-gray-25 focus:bg-gray-800'
      >
        {children}
      </Button>
    </Tooltip>
  );
};

import React from 'react';

import { Button } from '@ui/form/Button/Button';
import { CommandKbd } from '@ui/overlay/CommandMenu';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ActionItemProps {
  dataTest?: string;
  onClick: () => void;
  icon: React.ReactElement;
  tooltip?: React.ReactNode;
  children: React.ReactNode;
}

export const ActionItem = ({
  icon,
  onClick,
  dataTest,
  tooltip,
  children,
}: ActionItemProps) => {
  return (
    <Tooltip
      className='p-1.5'
      label={
        tooltip ? (
          tooltip
        ) : (
          <>
            <div className='flex items-center text-sm'>
              Open command menu
              <CommandKbd className='bg-gray-600 text-gray-25 mx-1' />
              <div className='bg-gray-600 text-xs min-h-5 min-w-5 rounded flex justify-center items-center'>
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
        dataTest={dataTest}
        className='bg-gray-700 text-gray-25 hover:bg-gray-800 hover:text-gray-25 focus:bg-gray-800'
      >
        {children}
      </Button>
    </Tooltip>
  );
};

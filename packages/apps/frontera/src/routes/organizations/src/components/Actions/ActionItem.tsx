import { Button } from '@ui/form/Button/Button';
import { Command } from '@ui/media/icons/Command.tsx';
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
      className='p-1 pl-2'
      label={
        tooltip ? (
          <div className='flex items-center text-sm'>
            {tooltip}{' '}
            <span className='bg-gray-600 text-xs px-1.5 rounded-sm leading-[1.125rem] ml-3'>
              {shortcutKey}
            </span>
          </div>
        ) : (
          <>
            <div className='flex items-center text-sm'>
              Open command menu
              <div className='bg-gray-600 h-5 w-5 rounded-sm ml-3 mr-1 flex flex justify-center items-center'>
                <Command className='size-3' />
              </div>
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

import { Button } from '@ui/form/Button/Button';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface ActionItemProps {
  tooltip?: string;
  onClick: () => void;
  shortcutKey?: string;
  icon: React.ReactElement;
  children: React.ReactNode;
}

export const ActionItem = ({
  icon,
  onClick,
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
        ) : undefined
      }
    >
      <Button
        leftIcon={icon}
        onClick={onClick}
        colorScheme='gray'
        className='bg-gray-700 text-gray-25 hover:bg-gray-800 hover:text-gray-25 focus:bg-gray-800'
      >
        {children}
      </Button>
    </Tooltip>
  );
};
import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Switch } from '@ui/form/Switch';
import { Icon, IconName } from '@ui/media/Icon';

interface HeaderProps {
  icon?: IconName;
  agentName: string;
  isActive: boolean;
  onToggleActive: () => void;
}

export const Header = observer(
  ({ onToggleActive, agentName, isActive, icon }: HeaderProps) => {
    return (
      <div className='w-full border-b border-b-gray-200 px-3 py-[11px] flex justify-between items-center'>
        <div className='flex items-center gap-1'>
          <Link
            to='/agents'
            className='text-md font-medium text-grayModern-500 hover:text-grayModern-700 transition-colors'
          >
            Agents
          </Link>
          <Icon name='chevron-right' className='w-4 h-4 text-gray-500' />
          <div className='flex items-center gap-1'>
            {icon && (
              <div className='flex items-center rounded-sm p-0.5 bg-grayModern-100'>
                <Icon name={icon} className='w-4 h-4 text-grayModern-500' />
              </div>
            )}

            <div className='flex items-center gap-1'>
              <p className='text-md font-medium'>{agentName}</p>
            </div>

            <div className='ml-4 flex items-center'>
              <Switch checked={isActive} onChange={onToggleActive} />
            </div>
          </div>
        </div>
      </div>
    );
  },
);

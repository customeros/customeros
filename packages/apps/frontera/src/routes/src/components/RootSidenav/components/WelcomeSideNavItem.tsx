import { useNavigate } from 'react-router-dom';
import { useState, ReactElement, MouseEventHandler } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { buttonSize } from '@ui/form/Button/Button';
import { ghostButton } from '@ui/form/Button/Button.variants';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

interface SidenavItemProps {
  href?: string;
  label: string;
  dataTest?: string;
  isActive?: boolean;
  onClick?: () => void;
  rightElement?: ReactElement | null;
  icon: ((isActive: boolean) => ReactElement) | ReactElement;
}

export const WelcomeSidenavItem = observer(
  ({ label, icon, onClick, isActive = false, dataTest }: SidenavItemProps) => {
    const [isEditing, setIsEditing] = useState(false);
    const store = useStore();
    const navigate = useNavigate();

    const handleClick: MouseEventHandler = (e) => {
      e.preventDefault();
      onClick?.();
    };

    const dynamicClasses = cn(
      isActive
        ? ['font-medium', 'bg-grayModern-100']
        : ['font-normal', 'bg-transparent'],
    );

    const handleArchive = () => {
      store.globalCache.updateOnboardingDetails(
        {
          showOnboardingPage: false,
        },
        {
          onSuccess: () => {
            if (isActive) {
              navigate(
                `/finder?preset=${store.tableViewDefs.organizationsPreset}`,
              );
            }
          },
        },
      );
    };

    return (
      <div
        data-test={dataTest}
        onClick={handleClick}
        className={cn(
          buttonSize({ size: 'sm' }),
          (ghostButton({ colorScheme: 'grayModern' }),
          `flex w-full justify-start gap-2 px-3 text-gray-700 cursor-pointer hover:bg-grayModern-100 *:hover:text-gray-700  group focus:shadow-EditableSideNavItemFocus mb-[2px] rounded-md ${dynamicClasses}`),
        )}
      >
        <div className='mt-[-1px]'>
          {typeof icon === 'function' ? icon(isActive!) : icon}
        </div>
        <div
          aria-selected={isActive}
          className={cn(
            'w-full text-justify whitespace-nowrap overflow-hidden overflow-ellipsis',
          )}
        >
          {label}
        </div>

        <div
          className={cn(
            'justify-end opacity-0 w-0 group-hover:opacity-100 group-focus:opacity-100 group-hover:w-6 group-focus:w-6',
            {
              'opacity-100 w-6': isEditing,
            },
          )}
        >
          <Menu open={isEditing} onOpenChange={setIsEditing}>
            <MenuButton className='min-w-6 h-5 rounded-md outline-none focus:outline-none text-gray-400 hover:text-gray-500'>
              <Icon name='dots-vertical' className='text-inherit' />
            </MenuButton>

            <MenuList align='end' side='bottom'>
              <MenuItem onClick={handleArchive}>
                <Icon name='archive' className='text-gray-500' />
                Archive
              </MenuItem>
            </MenuList>
          </Menu>
        </div>
      </div>
    );
  },
);

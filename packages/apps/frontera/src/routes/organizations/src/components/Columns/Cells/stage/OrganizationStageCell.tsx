import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { OrganizationStage } from '@graphql/types';
import { Seeding } from '@ui/media/icons/Seeding.tsx';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart.tsx';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart.tsx';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle.tsx';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  stageOptions,
  getStageOptions,
} from '@organization/components/Tabs/panels/AboutPanel/util.ts';

const iconMap = {
  Customer: <ActivityHeart className='text-gray-500' />,
  Prospect: <Seeding className='text-gray-500' />,
  'Not a fit': <MessageXCircle className='text-gray-500' />,
  'Former Customer': <BrokenHeart className='text-gray-500' />,
};

interface RenewalLikelihoodCellProps {
  id: string;
}

export const OrganizationStageCell = observer(
  ({ id }: RenewalLikelihoodCellProps) => {
    const [isEdit, setIsEdit] = useState(false);
    const [metaKey, setMetaKey] = useState(false);

    const store = useStore();
    const organization = store.organizations.value.get(id);
    useEffect(() => {
      store.ui.setIsEditingTableCell(isEdit);
    }, [isEdit, store.ui]);

    const selectedStageOption = stageOptions.find(
      (option) => option.value === organization?.value.stage,
    );
    const applicableStageOptions = getStageOptions(
      organization?.value?.relationship,
    );

    const menuHandleChange = (value: OrganizationStage) => {
      organization?.update((org) => {
        org.stage = value;

        return org;
      });
    };

    return (
      <div
        className='flex-1'
        onDoubleClick={() => setIsEdit(true)}
        onKeyDown={(e) => e.metaKey && setMetaKey(true)}
        onKeyUp={() => metaKey && setMetaKey(false)}
        onClick={() => metaKey && setIsEdit(true)}
        onBlur={() => setIsEdit(false)}
      >
        <Menu>
          <MenuButton className='outline-none focus:outline-none'>
            <span>
              {selectedStageOption?.label ? (
                selectedStageOption?.label
              ) : (
                <span className='text-gray-400'>Unknown</span>
              )}
            </span>
          </MenuButton>
          <MenuList side='bottom' align='center' className='min-w-[280px]'>
            {applicableStageOptions.map((option) => (
              <MenuItem
                className='ml-0'
                key={option.value}
                onClick={() => {
                  menuHandleChange(option.value);
                }}
              >
                {iconMap[option.label as keyof typeof iconMap]}
                {option.label}
              </MenuItem>
            ))}
          </MenuList>
        </Menu>
      </div>
    );
  },
);

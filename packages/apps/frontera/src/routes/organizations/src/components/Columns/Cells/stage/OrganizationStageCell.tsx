import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { OrganizationStage } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03.tsx';
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
        className='flex gap-1 group/stage'
        onKeyDown={(e) => e.metaKey && setMetaKey(true)}
        onKeyUp={() => metaKey && setMetaKey(false)}
        onClick={() => metaKey && setIsEdit(true)}
      >
        <p
          className={cn(
            'cursor-default text-gray-700',
            !selectedStageOption?.value && 'text-gray-400',
          )}
          data-test='organization-stage-in-all-orgs-table'
          onDoubleClick={() => setIsEdit(true)}
        >
          {selectedStageOption?.label ?? 'Not applicable'}
        </p>
        <Menu open={isEdit} onOpenChange={setIsEdit}>
          <MenuButton>
            {!!applicableStageOptions.length && (
              <IconButton
                className={cn(
                  'rounded-md opacity-0 group-hover/stage:opacity-100',
                  isEdit && 'opacity-100',
                )}
                aria-label='edit stage'
                size='xxs'
                variant='ghost'
                id='edit-button'
                onClick={() => setIsEdit(true)}
                icon={<Edit03 className='text-gray-500' />}
              />
            )}
          </MenuButton>

          <MenuList
            side='bottom'
            align='center'
            className={cn('min-w-[280px]', {
              hidden: !applicableStageOptions.length,
            })}
          >
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

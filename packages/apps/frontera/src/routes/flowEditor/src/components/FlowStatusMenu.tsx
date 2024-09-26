import React from 'react';

import { observer } from 'mobx-react-lite';

import { FlowStatus } from '@graphql/types';
import { Play } from '@ui/media/icons/Play';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { DotLive } from '@ui/media/icons/DotLive';
import { StopCircle } from '@ui/media/icons/StopCircle';
import { Tag, TagLabel, TagLeftIcon } from '@ui/presentation/Tag';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

interface FlowStatusMenuSelectProps {
  id: string;
}

export const FlowStatusMenu = observer(({ id }: FlowStatusMenuSelectProps) => {
  const store = useStore();

  const flow = store.flows.value.get(id);
  const status = flow?.value.status;

  if (status !== FlowStatus.Active) {
    return (
      <Button
        size='xs'
        variant='outline'
        leftIcon={<Play />}
        colorScheme='primary'
        onClick={() => {
          store.ui.commandMenu.toggle('StartFlow');
        }}
      >
        Start flow
      </Button>
    );
  }

  return (
    <>
      <Menu>
        <MenuButton
          className='text-success-500'
          data-test='flow-editor-status-change-button'
        >
          <Tag
            variant='outline'
            colorScheme='success'
            className='h-full rounded-md px-2'
          >
            <TagLeftIcon>
              <div>
                <DotLive className='text-success-500 mr-1 [&>*:nth-child(1)]:fill-success-200 [&>*:nth-child(1)]:stroke-success-300 [&>*:nth-child(2)]:fill-success-600 ' />
              </div>
            </TagLeftIcon>
            <TagLabel className='text-success-500'>Live</TagLabel>
          </Tag>
        </MenuButton>
        <MenuList align='end' side='bottom' className='p-0'>
          <MenuItem
            className='flex items-center text-base'
            data-test='contract-menu-delete-contract'
            onClick={() => store.ui.commandMenu.toggle('StopFlow')}
          >
            <StopCircle className='mr-1 text-gray-500' />
            Stop flow...
          </MenuItem>
        </MenuList>
      </Menu>
    </>
  );
});

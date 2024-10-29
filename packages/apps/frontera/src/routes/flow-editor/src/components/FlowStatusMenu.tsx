import React from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { FlowStatus } from '@graphql/types';
import { Play } from '@ui/media/icons/Play';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { DotLive } from '@ui/media/icons/DotLive';
import { StopCircle } from '@ui/media/icons/StopCircle';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
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
      <Tooltip
        label={
          status === FlowStatus.Scheduling
            ? 'We’re scheduling this flow’s contacts'
            : ''
        }
      >
        <div>
          <Button
            size='xs'
            variant='outline'
            leftIcon={<Play />}
            dataTest='start-flow'
            loadingText='Scheduling...'
            isLoading={status === FlowStatus.Scheduling}
            colorScheme={status === FlowStatus.Scheduling ? 'gray' : 'primary'}
            onClick={() => {
              store.ui.commandMenu.toggle('StartFlow');
            }}
            className={cn({
              'text-gray-500 pointer-events-none':
                status === FlowStatus.Scheduling,
            })}
            leftSpinner={
              <Spinner
                size='sm'
                label='Scheduling'
                className='text-gray-500 fill-gray-200'
              />
            }
          >
            Start flow
          </Button>
        </div>
      </Tooltip>
    );
  }

  return (
    <>
      <Menu>
        <MenuButton
          className='text-success-500 h-full'
          data-test='flow-editor-status-change-button'
        >
          <Tag
            variant='outline'
            colorScheme='success'
            className='h-full rounded-md px-2 py-1'
          >
            <TagLeftIcon>
              <div>
                <DotLive className='text-success-500 mr-1 [&>*:nth-child(1)]:fill-success-200 [&>*:nth-child(1)]:stroke-success-300 [&>*:nth-child(2)]:fill-success-600 ' />
              </div>
            </TagLeftIcon>
            <TagLabel className='text-success-500'>Live</TagLabel>
          </Tag>
        </MenuButton>
        <MenuList align='end' side='bottom' className='p-0 z-[11]'>
          <MenuItem
            className='flex items-center '
            data-test='stop-flow-menu-button'
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

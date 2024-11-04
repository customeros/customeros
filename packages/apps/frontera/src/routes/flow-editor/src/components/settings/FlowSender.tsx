import { ReactNode } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowSenderStore } from '@store/FlowSenders/FlowSender.store';

import { Avatar } from '@ui/media/Avatar';
import { User01 } from '@ui/media/icons/User01';
import { Delete } from '@ui/media/icons/Delete';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover/Popover.tsx';

export const FlowSender = observer(
  ({ id, flowId }: { id: string; flowId: string }) => {
    const store = useStore();
    const flowSender = store.flowSenders.value.get(id) as FlowSenderStore;
    const userMailboxes = flowSender?.user?.value?.mailboxes;

    return (
      <div className='flex justify-between'>
        <div className='flex'>
          <Avatar
            size='xs'
            textSize='xxs'
            name={flowSender?.user?.name ?? 'Unnamed'}
            icon={<User01 className='text-gray-500 size-3' />}
            src={flowSender?.user?.value?.profilePhotoUrl ?? ''}
            className={'w-5 h-5 min-w-5 mr-2 border border-gray-200'}
          />
          <span className='flex-1 text-sm'>{flowSender?.user?.name}</span>
        </div>
        <div className='flex'>
          <Popover>
            <PopoverTrigger asChild>
              <Button
                size='xxs'
                variant='ghost'
                leftIcon={<Mail01 className='text-inherit ' />}
              >
                {userMailboxes?.length ?? 0}{' '}
                {userMailboxes?.length === 1 ? 'mailbox' : 'mailboxes'}
              </Button>
            </PopoverTrigger>
            <PopoverContent align='end' side='bottom' className={'px-4 py-3 '}>
              <p className='text-sm font-medium mb-1'>Linked mailboxes</p>
              <ul className='list-disc px-4 text-sm'>
                {userMailboxes?.map((mailbox) => (
                  <li className='' key={`mailbox-item-${mailbox}`}>
                    {mailbox}
                  </li>
                ))}
              </ul>
            </PopoverContent>
          </Popover>

          <FlowSenderMenu senderId={id} flowId={flowId}>
            <IconButton
              size='xxs'
              aria-label={''}
              variant='ghost'
              icon={<DotsVertical className='text-inherit' />}
            />
          </FlowSenderMenu>
        </div>
      </div>
    );
  },
);

const FlowSenderMenu = observer(
  ({
    flowId,
    senderId,
    children,
  }: {
    flowId: string;
    senderId: string;
    children: ReactNode;
  }) => {
    const store = useStore();

    return (
      <Menu>
        <MenuButton
          data-test='flow-sender-mailbox-menu'
          className='outline-none focus:outline-none'
        >
          {children}
        </MenuButton>
        <MenuList align='end' side='bottom' className='min-w-[280px]'>
          <MenuItem
            onClick={() => {
              store.flowSenders.deleteFlowSender(senderId, flowId);
            }}
          >
            <Delete />
            Remove sender
          </MenuItem>
        </MenuList>
      </Menu>
    );
  },
);

import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { DisconnectSlackChannelUsecase } from '@domain/usecases/agents/capabilities/disconnect-slack-channel.usecase.ts';

import { Icon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button.tsx';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import {
  Menu,
  MenuList,
  MenuItem,
  MenuButton,
} from '@ui/overlay/Menu/Menu.tsx';

export const DisconnectSlackMenu = observer(() => {
  const usecase = useMemo(() => new DisconnectSlackChannelUsecase(), []);

  return (
    <>
      <Menu>
        <MenuButton asChild>
          <Button size='xxs' variant='ghost'>
            <Icon stroke={'none'} name={'dot-live-success'} />
            {'Connected'}
          </Button>
        </MenuButton>
        <MenuList>
          <MenuItem
            onClick={() => {
              usecase.toggleConfirmationDialog(true);
            }}
          >
            <Icon name='link-broken-02' className='text-grayModern-500' />
            Disconnect Slack
          </MenuItem>
        </MenuList>
      </Menu>
      <ConfirmDeleteDialog
        hideCloseButton
        colorScheme='primary'
        isOpen={usecase.isOpen}
        label={`Disconnect Slack?`}
        confirmButtonLabel='Disconnect'
        onConfirm={() => usecase.execute()}
        onClose={() => usecase.toggleConfirmationDialog(false)}
        description='Disconnecting Slack will revoke its access to your workspace. This means that any part of the app, including agent capabilities, that relies on it will stop working.'
      />
    </>
  );
});

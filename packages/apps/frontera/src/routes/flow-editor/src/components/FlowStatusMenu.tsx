import React from 'react';

import { observer } from 'mobx-react-lite';
import { useReactFlow } from '@xyflow/react';
import { FlowActionType } from '@store/Flows/types.ts';

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
  handleOpenSettingsPanel: () => void;
}

export const FlowStatusMenu = observer(
  ({ id, handleOpenSettingsPanel }: FlowStatusMenuSelectProps) => {
    const store = useStore();

    const flow = store.flows.value.get(id);
    const status = flow?.value.status;
    const { getNodes } = useReactFlow();

    const showValidationMessage = (
      meta: (typeof ValidationMessage)[keyof typeof ValidationMessage],
      openSettings = false,
    ) => {
      if (openSettings) {
        store.ui.commandMenu.setCallback(handleOpenSettingsPanel);
      }
      store.ui.commandMenu.toggle('FlowValidationMessage', {
        ...store.ui.commandMenu.context,
        meta,
      });
    };

    const handleFlowValidation = () => {
      const steps = getNodes();
      const hasOnlyControlSteps = steps.every(
        (step) => step.type === 'trigger' || step.type === 'control',
      );

      const hasEmailSteps = steps.some(
        (e) => e.data.action === FlowActionType.EMAIL_NEW,
      );
      const hasLinkedInSteps = steps.some(
        (e) => e.data.action === FlowActionType.LINKEDIN_CONNECTION_REQUEST,
      );

      if (hasOnlyControlSteps) {
        showValidationMessage(ValidationMessage.NO_STEPS);

        return true;
      }

      if (!flow?.value.senders.length) {
        showValidationMessage(ValidationMessage.NO_SENDERS, true);

        return true;
      }
      const senders = flow.value.senders
        .map(
          (sender) =>
            sender.user?.id && store.users.value.get(sender.user.id)?.value,
        )
        .filter(Boolean);

      if (
        hasEmailSteps &&
        hasLinkedInSteps &&
        senders.every(
          (sender) =>
            !sender || !sender?.mailboxes.length || !sender?.hasLinkedInToken,
        )
      ) {
        showValidationMessage(
          ValidationMessage.NO_MAILBOXES_AND_LINKEDIN,
          true,
        );

        return true;
      }

      if (
        hasEmailSteps &&
        senders.every((sender) => !sender || !sender?.mailboxes.length)
      ) {
        showValidationMessage(ValidationMessage.NO_MAILBOXES, true);

        return true;
      }

      if (
        hasLinkedInSteps &&
        senders.every((sender) => !sender || !sender?.hasLinkedInToken)
      ) {
        showValidationMessage(ValidationMessage.NO_LINKEDIN, true);

        return true;
      }

      return false;
    };

    const handleStartFlow = () => {
      const hasErrors = handleFlowValidation();

      if (hasErrors) return;

      store.ui.commandMenu.toggle('StartFlow');
    };

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
              onClick={handleStartFlow}
              loadingText='Scheduling...'
              isLoading={status === FlowStatus.Scheduling}
              colorScheme={
                status === FlowStatus.Scheduling ? 'gray' : 'primary'
              }
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
  },
);
const ValidationMessage = {
  NO_STEPS: {
    title: 'Add at least one step',
    description: 'To start this flow, add at least one step to it',
    buttonText: 'Got it',
  },
  NO_SENDERS: {
    title: 'Add a sender first',
    description: 'To start this flow, add at least one sender to it',
    buttonText: 'Go to flow settings',
  },
  NO_MAILBOXES_AND_LINKEDIN: {
    title: 'Add mailboxes and the LinkedIn extension',
    description:
      'To start this flow, add some mailboxes and the CustomerOS LinkedIn browser extension to a sender.',
    buttonText: 'Go to flow settings',
  },
  NO_MAILBOXES: {
    title: 'Add some mailboxes first',
    description: 'To start this flow, add some mailboxes to a sender',
    buttonText: 'Go to flow settings',
  },
  NO_LINKEDIN: {
    title: 'Add the LinkedIn browser extension first',
    description:
      'To start this flow, add the CustomerOS LinkedIn browser extension to a sender',
    buttonText: 'Go to flow settings',
  },
};

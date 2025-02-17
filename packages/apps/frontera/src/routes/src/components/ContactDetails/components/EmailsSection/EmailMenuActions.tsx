import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { EditEmailCase } from '@domain/usecases/command-menu/edit-email.usecase';
import { EmailVerificationStatus } from '@finder/components/Columns/contacts/filterTypes';

import { Copy01 } from '@ui/media/icons/Copy01';
import { Star01 } from '@ui/media/icons/Star01';
import { Send03 } from '@ui/media/icons/Send03';
import { IconButton } from '@ui/form/IconButton';
import { useEvent } from '@shared/hooks/useEvent';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { Stars02 } from '@ui/media/icons/Stars02';
import { TextInput } from '@ui/media/icons/TextInput';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  EmailDeliverable,
  EmailValidationDetails,
} from '@shared/types/__generated__/graphql.types';

import { emailStatuses } from './EmailValidationMessage';

const editEmailUseCase = EditEmailCase.getInstance();

interface EmailMenuActionsProps {
  id: string;
  idx: number;
  email: string;
  contactId: string;
}

export const EmailMenuActions = observer(
  ({ contactId, idx, email }: EmailMenuActionsProps) => {
    const store = useStore();
    const [_, copyToClipboard] = useCopyToClipboard();
    const navigate = useNavigate();

    const { dispatchEvent } = useEvent<{ email: string; openEditor: string }>(
      'openEmailEditor',
    );

    const contactStore = store.contacts.value.get(contactId);

    const isPrimaryEmail = contactStore?.value.emails[idx]?.primary;
    const emailValidationDetails =
      contactStore?.value.emails[idx]?.emailValidationDetails;
    const isNotDeliverable = checkEmailStatus(emailValidationDetails, email);

    if (!contactStore) return;

    return (
      <Menu>
        <MenuButton asChild>
          <IconButton
            size='xxs'
            variant='ghost'
            icon={<DotsVertical />}
            aria-label='add new email'
          />
        </MenuButton>
        <MenuList>
          <MenuItem
            className='group/edit-email'
            onClick={() => {
              editEmailUseCase.setEmail(email);
              store.ui.commandMenu.setType('EditEmail');
              store.ui.commandMenu.setContext({
                ids: [contactStore?.value.id ?? ''],
                entity: 'Contact',
                property: 'email',
              });
              store.ui.commandMenu.setOpen(true);
            }}
          >
            <div className='flex items-center gap-2'>
              <TextInput className='group-hover/edit-email:text-gray-700 text-gray-500' />
              <span>Edit email</span>
            </div>
          </MenuItem>
          {!isPrimaryEmail && (
            <MenuItem
              className='group/primary-email'
              onClick={() => {
                contactStore.value.emails.forEach((e) => (e.primary = false));
                contactStore?.draft();
                contactStore.value.emails[idx].primary = true;
                contactStore?.commit();
              }}
            >
              <div className='flex items-center gap-2'>
                <Star01 className='group-hover/primary-email:text-gray-700 text-gray-500' />
                <span>Make primary</span>
              </div>
            </MenuItem>
          )}
          {isNotDeliverable?.value !== EmailVerificationStatus.InvalidMailbox &&
            isNotDeliverable?.value !== EmailVerificationStatus.MailboxFull &&
            isNotDeliverable?.value !==
              EmailVerificationStatus.IncorrectFormat &&
            contactStore.value.primaryOrganizationId && (
              <MenuItem
                className='group/send-email'
                onClick={() => {
                  navigate(
                    `/organization/${contactStore.value.primaryOrganizationId}?tab=people`,
                  );
                  setTimeout(() => {
                    dispatchEvent({ email: email, openEditor: 'email' });
                  }, 0);
                }}
              >
                <Send03 className='text-gray-500 group-hover/send-email:text-gray-700' />
                Send email to contact
              </MenuItem>
            )}
          <MenuItem
            className='group/add-email'
            onClick={() => {
              store.ui.commandMenu.setContext({
                ids: [contactStore?.id || ''],
                entity: 'Contact',
                property: 'email',
              });
              store.ui.commandMenu.setType('AddEmail');
              store.ui.commandMenu.setOpen(true);
            }}
          >
            <PlusCircle className='text-gray-500 group-hover/add-email:text-gray-700' />
            Add email
          </MenuItem>
          <MenuItem
            className='group/work-email'
            onClick={() => {
              contactStore.findEmail();
            }}
          >
            <Stars02 className='text-gray-500 group-hover/work-email:text-gray-700' />
            Find work email
          </MenuItem>
          <MenuItem
            className='group/archive-email'
            onClick={() => {
              contactStore.draft();
              contactStore.value.emails[idx].email = '';

              contactStore.value.emails.splice(idx, 1);
              contactStore?.commit();
            }}
          >
            <Archive className='text-gray-500 group-hover/archive-email:text-gray-700' />
            Archive email
          </MenuItem>

          <MenuItem
            className='group/copy-email'
            onClick={() => copyToClipboard(email || '', 'Email copied')}
          >
            <Copy01 className='group-hover/copy-email:text-gray-700 text-gray-500' />
            Copy email
          </MenuItem>
        </MenuList>
      </Menu>
    );
  },
);

function isValidEmail(email: string) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  return emailRegex.test(email);
}

function checkEmailStatus(emailData?: EmailValidationDetails, email?: string) {
  if (!email) {
    return null;
  }

  if (email && !emailData) {
    const isValidSyntax = isValidEmail(email);

    if (!isValidSyntax) return emailStatuses.INCORRECT_FORMAT;
  }

  if (emailData?.deliverable === EmailDeliverable.Deliverable) {
    if (emailData.isFirewalled) return emailStatuses.DELIVERABLE_FIREWALL;
    if (emailData.isFreeAccount) return emailStatuses.DELIVERABLE_FREE_ACCOUNT;

    return emailStatuses.DELIVERABLE_NO_RISK;
  }

  if (emailData?.deliverable === EmailDeliverable.Unknown) {
    return emailData.isCatchAll
      ? emailStatuses.CATCH_ALL
      : emailStatuses.UNABLE_TO_VALIDATE;
  }

  if (emailData?.deliverable === EmailDeliverable.Undeliverable) {
    if (emailData.isMailboxFull) return emailStatuses.MAILBOX_FULL;

    return emailStatuses.INVALID_MAILBOX;
  }

  if (emailData?.verified === false) return emailStatuses.NOT_VERIFIED;
  if (emailData?.verifyingCheckAll)
    return emailStatuses.VERIFICATION_IN_PROGRESS;
}

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const UsersCard = observer(() => {
  const store = useStore();
  const { open, onOpen, onClose } = useDisclosure();
  const [error1, error2] = store.mailboxes.invalidUsernames;
  const hasDomains =
    store.mailboxes.value.size > 0
      ? store.mailboxes.extendedBundle.size > 0
      : store.mailboxes.baseBundle.size > 0;

  const dirty = store.mailboxes.dirty;
  const isUsername1Dirty = dirty.get('username1') || false;
  const isUsername2Dirty = dirty.get('username2') || false;

  if (!hasDomains) return null;

  return (
    <>
      <Card className='py-2 px-3 bg-white'>
        <CardHeader className='flex flex-col'>
          <div className='flex items-end gap-1 pb-1'>
            <span className='font-medium text-sm '>Add 2 usernames</span>
            <IconButton
              size='xxs'
              variant='ghost'
              onClick={onOpen}
              aria-label='info'
              icon={<InfoCircle />}
            />
          </div>
          <span className='text-sm'>
            Your usernames will apply to all selected domains
          </span>
        </CardHeader>

        <CardContent className='flex flex-col p-0'>
          <Input
            size='sm'
            variant='outline'
            placeholder='E.g. john'
            value={store.mailboxes.usernames[0]}
            dataTest='settings-mailboxes-first-username'
            invalid={isUsername1Dirty && error1.length > 0}
            onKeyDown={(e) => {
              e.stopPropagation();
            }}
            className={cn(
              'w-full',
              (!isUsername1Dirty ||
                (isUsername1Dirty && error1.length === 0)) &&
                'mb-[18px]',
            )}
            onChange={(e) => {
              store.mailboxes.setUsername(
                0,
                e.target.value.trim().toLowerCase(),
              );

              if (isUsername1Dirty) {
                store.mailboxes.validateUsernames();
              }
            }}
            onBlur={() => {
              if (
                !isUsername1Dirty &&
                store.mailboxes.usernames[0].length > 0
              ) {
                store.mailboxes.setDirty('username1');
                store.mailboxes.validateUsernames();
              }
            }}
          />
          {isUsername1Dirty && error1.length > 0 && (
            <span className='text-[12px] ml-[9px] text-error-400'>
              {error1}
            </span>
          )}
          <Input
            size='sm'
            variant='outline'
            placeholder='E.g. melinda'
            value={store.mailboxes.usernames[1]}
            dataTest='settings-mailboxes-second-username'
            invalid={isUsername2Dirty && error2.length > 0}
            onKeyDown={(e) => {
              e.stopPropagation();
            }}
            className={cn(
              'w-full mt-0.5',
              (!isUsername2Dirty ||
                (isUsername2Dirty && error2.length === 0)) &&
                'mb-[18px]',
            )}
            onChange={(e) => {
              store.mailboxes.setUsername(
                1,
                e.target.value.trim().toLowerCase(),
              );

              if (isUsername2Dirty) {
                store.mailboxes.validateUsernames();
              }
            }}
            onBlur={() => {
              if (
                !isUsername2Dirty &&
                store.mailboxes.usernames[1].length > 0
              ) {
                store.mailboxes.setDirty('username2');
                store.mailboxes.validateUsernames();
              }
            }}
          />
          {isUsername2Dirty && error2.length > 0 && (
            <span className='text-[12px] ml-[9px] text-error-400'>
              {error2}
            </span>
          )}
        </CardContent>
      </Card>

      <InfoDialog
        isOpen={open}
        onClose={onClose}
        onConfirm={onClose}
        confirmButtonLabel='Got it'
        label='Mailbox best practices'
        body={
          <p className='text-sm'>
            Two mailboxes per domain has shown to be effective for maintaining
            deliverability, avoiding spam filters, and supporting rotation and
            inbox warming.
          </p>
        }
      />
    </>
  );
});

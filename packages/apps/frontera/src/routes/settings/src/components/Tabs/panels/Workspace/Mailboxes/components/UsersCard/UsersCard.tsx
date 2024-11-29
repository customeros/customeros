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

  if (store.mailboxes.baseBundle.size === 0) return null;

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
            invalid={error1.length > 0}
            value={store.mailboxes.usernames[0]}
            onBlur={store.mailboxes.validateUsernames}
            className={cn('w-full', error1.length === 0 && 'mb-[18px]')}
            onChange={(e) => {
              store.mailboxes.setUsername(0, e.target.value.trim());
            }}
          />
          {error1.length > 0 && (
            <span className='text-[12px] ml-[9px] text-error-400'>
              {error1}
            </span>
          )}
          <Input
            size='sm'
            variant='outline'
            placeholder='E.g. melinda'
            invalid={error2.length > 0}
            value={store.mailboxes.usernames[1]}
            onBlur={store.mailboxes.validateUsernames}
            className={cn('w-full mt-0.5', error1.length === 0 && 'mb-[18px]')}
            onChange={(e) => {
              store.mailboxes.setUsername(1, e.target.value.trim());
            }}
          />
          {error2.length > 0 && (
            <span className='text-[12px] ml-[9px] text-error-400'>
              {error1}
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

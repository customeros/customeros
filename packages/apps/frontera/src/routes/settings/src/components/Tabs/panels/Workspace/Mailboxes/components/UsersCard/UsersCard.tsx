import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const UsersCard = observer(() => {
  const checkout = true;
  const store = useStore();

  return (
    store.mailboxes.baseBundle.size > 0 && (
      <Card className='py-2 px-3 bg-white'>
        <CardHeader className='flex flex-col'>
          <div className='flex items-end gap-1 pb-1'>
            <span className='font-medium text-sm '>Add 2 usernames</span>
            <IconButton
              size='xxs'
              variant='ghost'
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
            className='w-full mb-[17px]'
            value={store.mailboxes.usernames[0]}
            required={checkout && store.mailboxes.usernames[0].length === 0}
            onChange={(e) => {
              store.mailboxes.setUsername(0, e.target.value.trim());
            }}
          />
          {checkout && store.mailboxes.usernames[0].length === 0 && (
            <span className='text-[12px] ml-[9px] text-error-400'>
              Houston, we have a blank...
            </span>
          )}
          <Input
            size='sm'
            variant='outline'
            placeholder='E.g. melinda'
            className='w-full mt-[2px]'
            value={store.mailboxes.usernames[1]}
            required={checkout && store.mailboxes.usernames[1].length === 0}
            onChange={(e) => {
              store.mailboxes.setUsername(1, e.target.value.trim());
            }}
          />
        </CardContent>
      </Card>
    )
  );
});

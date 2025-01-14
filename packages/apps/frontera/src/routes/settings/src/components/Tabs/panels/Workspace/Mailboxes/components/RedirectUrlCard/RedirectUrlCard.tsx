import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const RedirectUrlCard = observer(() => {
  const store = useStore();
  const dirty = store.mailboxes.dirty.get('redirectUrl') || false;
  const hasDomains =
    store.mailboxes.value.size > 0
      ? store.mailboxes.extendedBundle.size > 0
      : store.mailboxes.baseBundle.size > 0;

  if (!hasDomains) return null;

  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='flex flex-col'>
        <div className='flex items-end gap-1 pb-1'>
          <span className='font-medium text-sm'>Domains redirect to...</span>
        </div>
        <span className='text-sm'>
          Your new domains will redirect to this website
        </span>
      </CardHeader>

      <CardContent className='flex flex-col p-0'>
        <Input
          size='sm'
          variant='outline'
          placeholder='Website URL'
          value={store.mailboxes.redirectUrl}
          dataTest='settings-mailboxes-redirect-url'
          invalid={dirty && store.mailboxes.invalidRedirectUrl.length > 0}
          onKeyDown={(e) => {
            e.stopPropagation();
          }}
          className={cn(
            'w-full',
            (!dirty ||
              (dirty && store.mailboxes.invalidRedirectUrl.length === 0)) &&
              'mb-[18px]',
          )}
          onChange={(e) => {
            store.mailboxes.setRedirectUrl(e.target.value.trim());

            if (dirty) {
              store.mailboxes.validateRedirectUrl();
            }
          }}
          onBlur={() => {
            if (!dirty && store.mailboxes.redirectUrl.length > 0) {
              store.mailboxes.setDirty('redirectUrl');
              store.mailboxes.validateRedirectUrl();
            }
          }}
        />
        {dirty && store.mailboxes.invalidRedirectUrl.length > 0 && (
          <span className='text-[12px] ml-[9px] text-error-400'>
            {store.mailboxes.invalidRedirectUrl}
          </span>
        )}
      </CardContent>
    </Card>
  );
});

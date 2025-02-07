import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';

export const ForwardEmailReplies = observer(() => {
  return (
    <article className='flex flex-col gap-4'>
      <div>
        <h1 className='text-sm font-medium pr-4'>Forward email replies</h1>
        <p className='text-sm mt-1'>
          When someone replies to an email, forward it to this email
        </p>
      </div>

      <Input
        size='xs'
        variant='outline'
        placeholder='Email'
        aria-label={'Email'}
        className='max-w-[320px]'
      />
    </article>
  );
});

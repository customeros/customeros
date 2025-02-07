import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { Input } from '@ui/form/Input';

export const ManageCampaignExecution = observer(() => {
  return (
    <article className='flex flex-col gap-4'>
      <h1 className='text-sm font-medium pr-4'>Schedule email delivery</h1>
      <div className=''>
        <h2 className='text-sm font-medium mb-1'>Send schedule</h2>
        <p className='text-sm'>
          Messages will be scheduled for weekdays, using the UTC timezone
        </p>
        <div className='flex items-center mt-3'>
          <div className='relative'>
            <Icon
              name={'clock'}
              className='absolute left-2 top-1/2 transform -translate-y-1/2 text-grayModern-500'
            />
            <Input
              size='xs'
              variant='outline'
              placeholder='8:00'
              className='max-w-[120px] pl-8'
            />
          </div>
          <div className='mx-2 w-3 h-[1px] bg-grayModern-700'></div>

          <div className='relative'>
            <Icon
              name={'clock'}
              className='absolute left-2 top-1/2 transform -translate-y-1/2 text-grayModern-500'
            />
            <Input
              size='xs'
              variant='outline'
              placeholder='18:00'
              className='max-w-[120px] pl-8'
            />
          </div>
        </div>
      </div>
    </article>
  );
});

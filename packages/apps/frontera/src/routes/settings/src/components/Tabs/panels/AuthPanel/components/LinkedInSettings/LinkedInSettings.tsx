import { useState } from 'react';
// import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Switch } from '@ui/form/Switch';
// import { useStore } from '@shared/hooks/useStore';

export const LinkedInSettings = observer(() => {
  // const store = useStore();
  // const [queryParams] = useSearchParams();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <article className='flex-col flex relative max-w-[550px] px-6 '>
        <div className='flex items-center w-full'>
          <div className='flex items-center gap-1'>
            <h2 className='text-gray-700 text-sm font-medium'>LinkedIn</h2>
          </div>
          <div className='w-full border-b border-gray-100 mx-2' />

          <div className='flex items-center'>
            <Switch
              isChecked={isOpen}
              colorScheme='primary'
              size='sm'
              onChange={(isChecked) => setIsOpen(isChecked)}
            />
          </div>
        </div>

        <p className='line-clamp-2 mt-2 mb-3 text-sm'>
          Import your LinkedIn connections by providing your email and password
        </p>

        {isOpen && (
          <>
            <label className='font-semibold text-sm'>
              Email or Phone
              <Input
                name='emailOrPhone'
                placeholder='olivia@untitledui.com'
                autoComplete='off'
                className='overflow-hidden overflow-ellipsis font-normal'
                value={''}
                onChange={() => {}}
              />
            </label>
            <label className='font-semibold text-sm'>
              Password
              <Input
                name='linkedInPassword'
                placeholder='*********'
                autoComplete='off'
                className='overflow-hidden overflow-ellipsis font-normal'
                value={''}
                onChange={() => {}}
              />
            </label>
          </>
        )}
      </article>
    </>
  );
});

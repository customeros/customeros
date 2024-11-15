import { useState } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const UsersCard = () => {
  const [user1, setUser1] = useState('');
  const [user2, setUser2] = useState('');
  const [userName, setUserName] = useLocalStorage<string[]>('userName', []);

  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='flex flex-col'>
        <div className='flex items-end gap-1'>
          <span className='font-medium'>Add 2 usernames</span>
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='info'
            icon={<InfoCircle />}
          />
        </div>
        <span>Your usernames will apply to all selected domains</span>
      </CardHeader>

      <CardContent className='flex flex-col gap-2 p-0 '>
        <Input
          size='sm'
          variant='outline'
          className='w-full'
          value={userName[0]}
          placeholder='E.g. john'
          onChange={(e) => setUser1(e.target.value)}
          onBlur={() => {
            const foundIndex = userName.findIndex((item) => item === user1);

            if (foundIndex === -1) {
              setUserName((prev) => [...prev, user1]);
            }
          }}
        />
        <Input
          size='sm'
          variant='outline'
          className='w-full'
          value={userName[1]}
          placeholder='E.g. melinda'
          onChange={(e) => setUser2(e.target.value)}
          onBlur={() => {
            const foundIndex = userName.findIndex((item) => item === user2);

            if (foundIndex === -1) {
              setUserName((prev) => [...prev, user2]);
            }
          }}
        />
      </CardContent>
    </Card>
  );
};

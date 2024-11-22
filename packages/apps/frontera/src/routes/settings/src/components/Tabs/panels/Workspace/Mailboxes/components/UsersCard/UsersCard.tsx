import { useState } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

interface UsersCardProps {
  checkout: boolean;
}

export const UsersCard = ({ checkout }: UsersCardProps) => {
  const [userName, setUserName] = useLocalStorage<string[]>('userName', []);
  const [storedBrandName, _setStoredBrandName] = useLocalStorage<string[]>(
    'brandName',
    [],
  );
  const [user1, setUser1] = useState<string>('');
  const [user2, setUser2] = useState<string>('');

  return (
    storedBrandName.length > 0 && (
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

        <CardContent className='flex flex-col p-0 '>
          <Input
            size='sm'
            variant='outline'
            className='w-full'
            value={userName[0]}
            placeholder='E.g. john'
            required={checkout && user1.length === 0}
            onChange={(e) => {
              setUser1(e.target.value.trim());
            }}
            onBlur={() => {
              if (user1?.length === 0) return;
              const foundIndex = userName.findIndex((item) => item === user1);

              if (foundIndex === -1) {
                setUserName((prev) => [...prev, user1]);
              }
            }}
          />
          {checkout && user1?.length === 0 && (
            <span className='text-[12px] ml-[9px] text-error-400'>
              Houston, we have a blank...
            </span>
          )}
          <Input
            size='sm'
            variant='outline'
            value={userName[1]}
            placeholder='E.g. melinda'
            className='w-full mt-[2px]'
            required={checkout && user2.length === 0 && user1.length === 0}
            onChange={(e) => {
              setUser2(e.target.value.trim());
            }}
            onBlur={() => {
              if (user2?.length === 0) return;
              const foundIndex = userName.findIndex((item) => item === user2);

              if (foundIndex === -1) {
                setUserName((prev) => [...prev, user2]);
              }
            }}
          />
        </CardContent>
      </Card>
    )
  );
};

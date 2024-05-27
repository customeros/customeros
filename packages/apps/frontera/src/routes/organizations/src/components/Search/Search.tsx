import { useEffect, startTransition } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { UserPresence } from '@shared/components/UserPresence';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

export const Search = observer(() => {
  const store = useStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    startTransition(() => {
      const value = event.target.value;

      setSearchParams(
        (prev) => {
          if (!value) {
            prev.delete('search');
          } else {
            prev.set('search', value);
          }

          return prev;
        },
        { replace: true },
      );
    });
  };

  useEffect(() => {
    setSearchParams((prev) => {
      prev.delete('search');

      return prev;
    });
  }, [preset]);

  return (
    <div className='flex items-center justify-between pr-4 w-full'>
      <InputGroup className='w-full bg-gray-25 hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-2'>
        <LeftElement className='ml-2'>
          <SearchSm className='size-5' />
        </LeftElement>
        <Input
          size='lg'
          autoCorrect='off'
          spellCheck={false}
          variant='unstyled'
          onChange={handleChange}
          placeholder={`organization name...`}
          defaultValue={searchParams.get('search') ?? ''}
        />
      </InputGroup>
      <UserPresence channelName={`finder:${store.session.value.tenant}`} />
    </div>
  );
});

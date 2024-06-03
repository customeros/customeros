import { useSearchParams } from 'react-router-dom';
import { useRef, useEffect, startTransition } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { ViewSettings } from '@shared/components/ViewSettings';
import { UserPresence } from '@shared/components/UserPresence';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

export const Search = observer(() => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
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

  useKeyBindings({
    '/': () => {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 0);
    },
  });

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-3'
    >
      <InputGroup className='w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-2'>
        <LeftElement className='ml-2'>
          <SearchSm className='size-5' />
        </LeftElement>
        <Input
          size='lg'
          ref={inputRef}
          autoCorrect='off'
          spellCheck={false}
          variant='unstyled'
          onChange={handleChange}
          placeholder={
            store.ui.isSearching !== 'organizations'
              ? `Search organizations (/ to search)`
              : 'e.g. CustomerOS...'
          }
          defaultValue={searchParams.get('search') ?? ''}
          onKeyUp={(e) => {
            if (e.code === 'Escape') {
              inputRef.current?.blur();
              store.ui.setIsSearching(null);
            }
          }}
          onFocus={() => {
            store.ui.setIsSearching('organizations');
            wrapperRef.current?.setAttribute('data-focused', '');
          }}
          onBlur={() => {
            store.ui.setIsSearching(null);
            wrapperRef.current?.removeAttribute('data-focused');
          }}
        />
      </InputGroup>
      <UserPresence channelName={`finder:${store.session.value.tenant}`} />
      <ViewSettings type='organizations' />
    </div>
  );
});

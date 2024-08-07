import { useRef, startTransition } from 'react';
import { useSearchParams } from 'react-router-dom';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

export const Search = observer(() => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    startTransition(() => {
      const value = event.target.value
        .toLowerCase()
        .normalize('NFD')
        .replace(/[\u0300-\u036f]/g, '');

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

  useKeyBindings(
    {
      '/': () => {
        setTimeout(() => {
          inputRef.current?.focus();
        }, 0);
      },
    },
    {
      when: !store.ui.isEditingTableCell && !store.ui.isFilteringTable,
    },
  );

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-2'
    >
      <InputGroup className='relative w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2 flex items-center gap-1'>
          <SearchSm className='size-5' />
          <span className='text-nowrap font-medium'>
            {store.ui.searchCount} opps:
          </span>
        </LeftElement>
        <Input
          size='md'
          ref={inputRef}
          autoCorrect='off'
          spellCheck={false}
          variant='unstyled'
          onChange={handleChange}
          defaultValue={searchParams.get('search') ?? ''}
          onBlur={() => {
            store.ui.setIsSearching(null);
            wrapperRef.current?.removeAttribute('data-focused');
          }}
          onFocus={() => {
            store.ui.setIsSearching('opportunities');
            wrapperRef.current?.setAttribute('data-focused', '');
          }}
          placeholder={
            store.ui.isSearching
              ? 'e.g. Opportunity, organization or owner...'
              : '/ to search'
          }
          onKeyUp={(e) => {
            if (
              e.code === 'Escape' ||
              e.code === 'ArrowUp' ||
              e.code === 'ArrowDown'
            ) {
              inputRef.current?.blur();
              store.ui.setIsSearching(null);
            }
          }}
        />
      </InputGroup>
    </div>
  );
});

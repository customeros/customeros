import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect, startTransition } from 'react';

import { useKeyBindings } from 'rooks';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import { SortingState } from '@tanstack/table-core';

import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { ViewSettings } from '@shared/components/ViewSettings';
import { UserPresence } from '@shared/components/UserPresence';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';

import {
  getColumnSortFn,
  getPredefinedFilterFn,
} from '../Columns/columnsDictionary';

export const Search = observer(() => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [sorting, _setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);

  const searchTerm = searchParams?.get('search');
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const organizations = store.organizations.toComputedArray((arr) => {
    const predefinedFilter = getPredefinedFilterFn(tableViewDef?.getFilters());
    if (predefinedFilter) {
      arr = arr.filter(predefinedFilter);
    }
    if (searchTerm) {
      arr = arr.filter((org) =>
        org.value.name?.toLowerCase().includes(searchTerm?.toLowerCase()),
      );
    }
    const columnId = sorting[0]?.id;
    const isDesc = sorting[0]?.desc;
    const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
      getColumnSortFn(columnId),
    );

    return computed;
  });

  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;
  const multiResultPlaceholder = (() => {
    switch (tableViewName) {
      case 'Nurture':
        return 'prospects';
      case 'Customers':
        return 'customers';
      case 'Leads':
        return 'leads';
      case 'Churn':
        return 'churned';
      case 'All orgs':
        return 'organizations';
      default:
        return 'organizations';
    }
  })();

  const singleResultPlaceholder = (() => {
    switch (tableViewName) {
      case 'Nurture':
        return 'prospect';
      case 'Customers':
        return 'customer';
      case 'Leads':
        return 'lead';
      case 'Churn':
        return 'churned';
      case 'All orgs':
        return 'organization';
      default:
        return 'organization';
    }
  })();

  const toatalOrganizations = organizations.length;

  const tableName =
    toatalOrganizations === 1
      ? singleResultPlaceholder
      : multiResultPlaceholder;

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

  useKeyBindings(
    {
      '/': () => {
        setTimeout(() => {
          inputRef.current?.focus();
        }, 0);
      },
    },
    {
      when: !store.ui.isEditingTableCell,
    },
  );

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-3'
    >
      <InputGroup className='w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2'>
          <div className='flex flex-row items-center gap-1'>
            <SearchSm className='size-5' />
            <span className='font-medium'>{toatalOrganizations}</span>
            <span className='font-medium'>{tableName}:</span>
          </div>
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
              ? `/ to search`
              : 'e.g. CustomerOS...'
          }
          defaultValue={searchParams.get('search') ?? ''}
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

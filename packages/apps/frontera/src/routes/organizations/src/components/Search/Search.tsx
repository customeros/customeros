import { useSearchParams } from 'react-router-dom';
import { useRef, useMemo, useState, useEffect, startTransition } from 'react';

import { useKeyBindings } from 'rooks';
import { Store } from '@store/store.ts';
import { inPlaceSort } from 'fast-sort';
import { observer } from 'mobx-react-lite';
import { SortingState } from '@tanstack/react-table';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { Button } from '@ui/form/Button/Button.tsx';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { ViewSettings } from '@shared/components/ViewSettings';
import { UserPresence } from '@shared/components/UserPresence';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';
import {
  Contact,
  TableIdType,
  Organization,
  TableViewType,
} from '@graphql/types';
import {
  getAllFilterFns,
  getColumnSortFn,
} from '@organizations/components/Columns/Dictionaries/columnsDictionary.tsx';
import {
  getContactFilterFn,
  getOrganizationFilterFn,
} from '@organizations/components/Columns/Dictionaries/SortAndFilterDictionary';

export const Search = observer(() => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [sorting, _setSorting] = useState<SortingState>([
    { id: 'ORGANIZATIONS_LAST_TOUCHPOINT', desc: true },
  ]);
  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;
  const tableViewType = store.tableViewDefs.getById(preset || '')?.value
    .tableType;

  const multiResultPlaceholder = (() => {
    switch (tableViewName) {
      case 'Targets':
        return 'targets';
      case 'Customers':
        return 'customers';
      case 'Contacts':
        return 'contacts';
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
      case 'Targets':
        return 'target';
      case 'Customers':
        return 'customer';
      case 'Contacts':
        return 'contact';
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

  const searchTerm = searchParams?.get('search');

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');

  const tableType = tableViewDef?.value?.tableType;
  const dataSet = useMemo(() => {
    if (tableType === TableViewType.Organizations) {
      return store.organizations;
    }
    if (tableType === TableViewType.Contacts) {
      return store.contacts;
    }

    return store.organizations;
  }, [tableType]);

  const filterFunction = useMemo(() => {
    if (tableType === TableViewType.Organizations) {
      return getOrganizationFilterFn;
    }
    if (tableType === TableViewType.Contacts) {
      return getContactFilterFn;
    }

    return getOrganizationFilterFn;
  }, [tableType]);

  // @ts-expect-error fixme
  const data = dataSet?.toComputedArray((arr) => {
    const filters = getAllFilterFns(tableViewDef?.getFilters(), filterFunction);
    if (filters) {
      // @ts-expect-error fixme

      arr = arr.filter((v) => filters.every((fn) => fn(v)));
    }

    if (searchTerm) {
      arr = arr.filter((entity) =>
        entity.value?.name
          ?.toLowerCase()
          .includes(searchTerm?.toLowerCase() as string),
      ) as Store<Contact>[] | Store<Organization>[];
    }
    if (tableType) {
      const columnId = sorting[0]?.id;
      const isDesc = sorting[0]?.desc;
      // @ts-expect-error fixme
      const computed = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
        getColumnSortFn(columnId, tableType),
      );

      return computed;
    }

    return arr;
  });

  const toatalResults = data?.length;

  const tableName =
    toatalResults === 1 ? singleResultPlaceholder : multiResultPlaceholder;

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
      when: !store.ui.isEditingTableCell || !store.ui.isFilteringTable,
    },
  );

  const placeholder =
    tableType === TableViewType.Contacts
      ? 'e.g. Isabella Evans'
      : 'e.g. CustomerOS...';

  const contactTableDef = store.tableViewDefs
    .toArray()
    .find((e) => e.value.tableType === TableViewType.Contacts)?.value.id;
  const targetTableDef = store.tableViewDefs
    .toArray()
    .find((e) => e.value.tableId === TableIdType.Nurture)?.value?.id;

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-3'
    >
      <InputGroup className='w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2'>
          <div className='flex flex-row items-center gap-1'>
            <SearchSm className='size-5' />
            <span className={'font-medium break-keep w-max mb-[2px]'}>
              {`${toatalResults} ${tableName}:`}
            </span>
          </div>
        </LeftElement>
        <Input
          size='md'
          ref={inputRef}
          autoCorrect='off'
          spellCheck={false}
          variant='unstyled'
          onChange={handleChange}
          placeholder={
            store.ui.isSearching !== 'organizations'
              ? `/ to search`
              : placeholder
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

      {(tableViewType === TableViewType.Contacts ||
        tableViewName === 'Targets') && (
        <ButtonGroup className='flex items-center '>
          <Button
            size='xs'
            className={cn('bg-white !border-r px-4', {
              'bg-gray-50 text-gray-500 font-normal': preset !== targetTableDef,
            })}
            onClick={() => {
              setSearchParams((prev) => {
                prev.set('preset', targetTableDef as string);

                return prev;
              });
            }}
          >
            Targets
          </Button>
          <Button
            size='xs'
            className={cn('bg-white px-4', {
              'bg-gray-50 text-gray-500 font-normal':
                preset !== contactTableDef,
            })}
            onClick={() => {
              setSearchParams((prev) => {
                prev.set('preset', contactTableDef as string);

                return prev;
              });
            }}
          >
            Contacts
          </Button>
        </ButtonGroup>
      )}

      {tableViewType && <ViewSettings type={tableViewType} />}
    </div>
  );
});

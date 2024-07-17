import { useSearchParams } from 'react-router-dom';
import { useRef, useEffect, startTransition } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Input } from '@ui/form/Input/Input';
import { Star06 } from '@ui/media/icons/Star06';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { TableIdType, TableViewType } from '@graphql/types';
import { ViewSettings } from '@shared/components/ViewSettings';
import { UserPresence } from '@shared/components/UserPresence';
import { InputGroup, LeftElement } from '@ui/form/InputGroup/InputGroup';
import { TargetNavigation } from '@organizations/components/TargetNavigation';
import { DownloadCsvButton } from '@organizations/components/DownloadCsvButton/DownloadCsvButton.tsx';

interface SearchProps {
  open: boolean;
  onOpen: () => void;
  onClose: () => void;
}

export const Search = observer(({ onClose, onOpen, open }: SearchProps) => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');

  const displayIcp = useFeatureIsOn('icp');

  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;
  const tableViewType = store.tableViewDefs.getById(preset || '')?.value
    .tableType;

  const tableId = store.tableViewDefs.getById(preset || '')?.value.tableId;
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
      case 'Past':
        return 'invoices';
      case 'Upcoming':
        return 'invoices';
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
      case 'Past':
        return 'invoice';
      case 'Upcoming':
        return 'invoice';
      case 'All orgs':
        return 'organization';
      default:
        return 'organization';
    }
  })();
  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;
  const toatalResults = store.ui.searchCount;

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
    onClose();
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
      when: !store.ui.isEditingTableCell && !store.ui.isFilteringTable,
    },
  );

  const placeholder =
    tableType === TableViewType.Contacts
      ? 'e.g. Isabella Evans'
      : 'e.g. CustomerOS...';

  const handleToogleFlow = () => {
    if (open) {
      onClose();
    } else {
      onOpen();
    }
  };

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-2'
    >
      <InputGroup className='w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2'>
          <div className='flex flex-row items-center gap-1'>
            <SearchSm className='size-5' />
            <span
              className={'font-medium break-keep w-max mb-[2px]'}
              data-test={`search-${tableName}`}
            >
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

      <TargetNavigation />

      {tableViewType && <ViewSettings type={tableViewType} />}

      {TableIdType.Leads === tableId && displayIcp && (
        <IconButton
          icon={<Star06 />}
          aria-label='toogle-flow'
          size='xs'
          onClick={handleToogleFlow}
        />
      )}
      {tableType !== TableViewType.Invoices && <DownloadCsvButton />}
    </div>
  );
});

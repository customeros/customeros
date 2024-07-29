import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect, startTransition } from 'react';

import { match } from 'ts-pattern';
import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Input } from '@ui/form/Input/Input';
import { Star06 } from '@ui/media/icons/Star06';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { TableIdType, TableViewType } from '@graphql/types';
import { ViewSettings } from '@shared/components/ViewSettings';
import { UserPresence } from '@shared/components/UserPresence';
import { useTablePlaceholder } from '@organizations/hooks/useTablePlaceholder.tsx';
import { ContactOrgViewToggle } from '@organizations/components/ContactOrgViewToggle';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '@ui/form/InputGroup/InputGroup';
import { DownloadCsvButton } from '@organizations/components/DownloadCsvButton/DownloadCsvButton.tsx';
import { SearchBarFilterData } from '@organizations/components/SearchBarFilterData/SearchBarFilterData.tsx';
import { CreateNewOrganizationModal } from '@organizations/components/shared/CreateNewOrganizationModal.tsx';

interface SearchProps {
  open: boolean;
  onOpen: () => void;
  onClose: () => void;
}

export const Search = observer(({ onClose, onOpen, open }: SearchProps) => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const [lastSearchForPreset] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-last-search-for-preset`, { root: 'root' });
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;

  const { multi: multiResultPlaceholder, single: singleResultPlaceholder } =
    useTablePlaceholder(tableViewName);
  const displayIcp = useFeatureIsOn('icp');

  useEffect(() => {
    onClose();

    setSearchParams(
      (prev) => {
        if (preset && lastSearchForPreset?.[preset]) {
          prev.set('search', lastSearchForPreset[preset]);
        } else {
          prev.delete('search');
        }

        return prev;
      },
      { replace: true },
    );
    if (preset && inputRef?.current) {
      inputRef.current.value = lastSearchForPreset[preset] ?? '';
    }
  }, [preset]);
  const tableViewType = store.tableViewDefs.getById(preset || '')?.value
    .tableType;
  const tableId = store.tableViewDefs.getById(preset || '')?.value.tableId;

  const tableViewDef = store.tableViewDefs.getById(preset ?? '1');
  const tableType = tableViewDef?.value?.tableType;
  const totalResults = store.ui.searchCount;

  const tableName =
    totalResults === 1 ? singleResultPlaceholder : multiResultPlaceholder;

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
  const placeholder = match(tableType)
    .with(TableViewType.Contacts, () => 'e.g. Isabella Evans')
    .with(TableViewType.Organizations, () => 'e.g. CustomerOS...')
    .with(TableViewType.Invoices, () => 'e.g. My contract')
    .otherwise(() => 'e.g. Organization name...');

  const handleToogleFlow = () => {
    if (open) {
      onClose();
    } else {
      onOpen();
    }
  };

  const allowCreation = totalResults === 0 && !!searchParams.get('search');

  useKeyBindings(
    {
      Enter: () => {
        store.ui.setIsEditingTableCell(true);
        setIsCreateModalOpen(true);
      },
    },
    { when: allowCreation },
  );

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full data-[focused]:animate-focus gap-2'
    >
      <InputGroup className='relative w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2'>
          <SearchBarFilterData />
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
        <RightElement>
          {allowCreation && (
            <div
              role='button'
              onClick={() => inputRef.current?.focus()}
              className='flex flex-row items-center gap-1 absolute top-[8px] cursor-text'
              style={{
                left: `calc(${measureRef?.current?.offsetWidth}px + 58px)`,
              }}
            >
              <Tag variant='subtle' colorScheme='grayBlue' className='mb-[2px]'>
                <TagLabel className='capitalize'>Enter</TagLabel>
              </Tag>
              <span className='font-normal text-gray-400 break-keep w-max mb-[2px]'>
                to create
              </span>
            </div>
          )}
        </RightElement>
      </InputGroup>
      <UserPresence channelName={`finder:${store.session.value.tenant}`} />

      <ContactOrgViewToggle />

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
      <span
        ref={measureRef}
        className={`z-[-1] absolute h-0 inline-block invisible`}
      >
        {searchParams.get('search')} {`${totalResults} ${tableName}:`}
      </span>

      <CreateNewOrganizationModal
        isOpen={isCreateModalOpen}
        setIsOpen={setIsCreateModalOpen}
      />
    </div>
  );
});

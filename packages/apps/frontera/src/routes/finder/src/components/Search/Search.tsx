import { useSearchParams } from 'react-router-dom';
import { useRef, useEffect, startTransition } from 'react';

import { match } from 'ts-pattern';
import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { TableViewMenu } from '@finder/components/TableViewMenu/TableViewMenu';
import { useSearchPersistence } from '@finder/components/Search/useSearchPersistance.ts';
import { CreateSequenceButton } from '@finder/components/Search/CreateSequenceButton.tsx';
import { TableViewsToggleNavigation } from '@finder/components/TableViewsToggleNavigation';
import { SearchBarFilterData } from '@finder/components/SearchBarFilterData/SearchBarFilterData';

import { Input } from '@ui/form/Input/Input';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { UserPlus01 } from '@ui/media/icons/UserPlus01.tsx';
import { TableIdType, TableViewType } from '@graphql/types';
import { UserPresence } from '@shared/components/UserPresence';
import {
  InputGroup,
  LeftElement,
  RightElement,
} from '@ui/form/InputGroup/InputGroup';

export const Search = observer(() => {
  const store = useStore();
  const wrapperRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLDivElement>(null);
  const floatingActionPropmterRef = useRef<HTMLDivElement>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const preset = searchParams.get('preset');
  const timeoutRef = useRef<NodeJS.Timeout>();

  const { lastSearchForPreset } = useSearchPersistence();

  useEffect(() => {
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
    preset &&
      store.organizations.setSearchTerm(lastSearchForPreset?.[preset], preset);

    if (preset && inputRef?.current) {
      inputRef.current.value = lastSearchForPreset[preset] ?? '';
    }
  }, [preset]);

  const tableViewDef = store.tableViewDefs.getById(preset || '');

  const tableType = tableViewDef?.value?.tableType;
  const totalResults = store.ui.searchCount;

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    startTransition(() => {
      const value = event.target.value;

      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        preset && store.organizations.setSearchTerm(value, preset);
      }, 200);

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
      when:
        !store.ui.isEditingTableCell &&
        !store.ui.isFilteringTable &&
        !store.ui.commandMenu.isOpen,
    },
  );

  const placeholder = match(tableType)
    .with(TableViewType.Flow, () => 'by flow name...')
    .with(TableViewType.Contacts, () => 'by name, organization or email...')
    .with(TableViewType.Contracts, () => 'by contract name...')
    .with(TableViewType.Organizations, () => 'by organization name...')
    .with(TableViewType.Invoices, () => 'by contract name...')
    .with(
      TableViewType.Opportunities,
      () => 'by name, organization or owner...',
    )
    .otherwise(() => 'by organization name...');

  const createNewEntityModalType:
    | null
    | 'CreateNewFlow'
    | 'AddContactsBulk'
    | 'AddNewOrganization' = match(tableType)
    .with(TableViewType.Flow, (): 'CreateNewFlow' => 'CreateNewFlow')
    .with(TableViewType.Contacts, (): 'AddContactsBulk' => 'AddContactsBulk')
    .with(
      TableViewType.Organizations,
      (): 'AddNewOrganization' => 'AddNewOrganization',
    )
    .otherwise(() => null);

  const allowCreation =
    ![TableIdType.Contracts, TableIdType.FlowActions].includes(
      tableViewDef?.value?.tableId as TableIdType,
    ) &&
    totalResults === 0 &&
    !!searchParams.get('search');

  return (
    <div
      ref={wrapperRef}
      className='flex items-center justify-between pr-1 w-full gap-2 bg-white border-b'
    >
      <InputGroup className='relative w-full bg-transparent hover:border-transparent focus-within:border-transparent focus-within:hover:border-transparent gap-1'>
        <LeftElement className='ml-2'>
          <SearchBarFilterData dataTest={'search-orgs'} />
        </LeftElement>
        <Input
          readOnly
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
            if (tableType === TableViewType.Organizations) return;
            store.ui.setIsSearching('organizations');
            wrapperRef.current?.setAttribute('data-focused', '');
          }}
          placeholder={(() => {
            if (tableType === TableViewType.Organizations) return '';

            return store.ui.isSearching !== 'organizations'
              ? `/ to search`
              : placeholder;
          })()}
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
          onKeyDown={(e) => {
            if (e.key === 'a' && (e.metaKey || e.ctrlKey)) {
              e.preventDefault();
              e.currentTarget.select();
            }

            if (e.key === 'Enter' && allowCreation) {
              e.stopPropagation();
              e.preventDefault();

              if (createNewEntityModalType) {
                store.ui.commandMenu.setType(createNewEntityModalType);
                store.ui.commandMenu.setOpen(true);
              }
            }
            e.stopPropagation();
          }}
        />
        <RightElement>
          {allowCreation && (
            <div
              role='button'
              ref={floatingActionPropmterRef}
              onClick={() => inputRef.current?.focus()}
              className='flex flex-row items-center gap-1 absolute top-[11px] cursor-text'
              style={{
                left: `calc(${measureRef?.current?.offsetWidth ?? 0}px + 24px)`,
              }}
            >
              <Tag variant='subtle' className='mb-[2px]' colorScheme='grayBlue'>
                <TagLabel className='capitalize'>Enter</TagLabel>
              </Tag>
              <span className='font-normal text-gray-400 break-keep w-max text-sm'>
                to create
              </span>
            </div>
          )}
        </RightElement>
      </InputGroup>
      <UserPresence channelName={`finder:${store.session.value.tenant}`} />

      <TableViewsToggleNavigation />

      {tableViewDef?.value.tableId === TableIdType.FlowActions && (
        <CreateSequenceButton />
      )}

      {tableViewDef?.value.tableId === TableIdType.Contacts && (
        <Button
          size='xs'
          variant='outline'
          colorScheme='primary'
          leftIcon={<UserPlus01 />}
          onClick={() => {
            store.ui.commandMenu.setOpen(true);
            store.ui.commandMenu.setType('AddContactsBulk');
          }}
        >
          Add contacts
        </Button>
      )}
      {tableViewDef?.value.tableId !== TableIdType.FlowActions && (
        <TableViewMenu />
      )}
      <span ref={measureRef} className={`z-[-1] absolute h-0 invisible flex`}>
        <div className='ml-2'>
          <SearchBarFilterData />
        </div>
        {inputRef?.current?.value ?? ''}
      </span>
    </div>
  );
});

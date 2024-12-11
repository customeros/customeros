import { useSearchParams } from 'react-router-dom';
import React, { useState, useCallback } from 'react';

import { debounce } from 'lodash';
import { observer } from 'mobx-react-lite';
import { useDidMount, useKeyBindings } from 'rooks';

import { Input } from '@ui/form/Input';
import { Avatar } from '@ui/media/Avatar';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { User03 } from '@ui/media/icons/User03.tsx';
import { PlusCircle } from '@ui/media/icons/PlusCircle.tsx';
import { Command, CommandCancelIconButton } from '@ui/overlay/CommandMenu';
import {
  ColumnViewType,
  OrganizationStage,
  OrganizationRelationship,
} from '@graphql/types';

export const AddNewOrganization = observer(() => {
  const store = useStore();
  const [allowSubmit, _setAllowSubmit] = useState(false);
  const { organizations, tableViewDefs } = useStore();
  const [searchParams] = useSearchParams();

  const [search, setSearch] = useState('');
  const [options, setOptions] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const preset = searchParams?.get('preset');

  const tableViewName = tableViewDefs.getById(`${preset}`)?.value.name;

  useDidMount(async () => {
    const result = await store.organizations.getGlobalOrganizationOptions(
      '',
      30,
    );

    const safeOptions = Array.isArray(result) ? result : [];

    setOptions(safeOptions);
  });

  const debouncedSearch = useCallback(
    debounce(async (searchTerm: string) => {
      try {
        setIsLoading(true);

        const result = await store.organizations.getGlobalOrganizationOptions(
          searchTerm,
          30,
        );

        const safeOptions = Array.isArray(result) ? result : [];

        setOptions(safeOptions);
      } catch (error) {
        console.error('Error fetching organization options:', error);
        setOptions([]);
      } finally {
        setIsLoading(false);
      }
    }, 300),
    [],
  );

  const handleSearchChange = (value: string) => {
    setSearch(value);
    debouncedSearch.cancel();
    debouncedSearch(value);
  };

  const handleClose = () => {
    store.ui.commandMenu.clearContext();
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.clearCallback();
  };

  const handleAddNewOrganization = () => {
    if (!allowSubmit) return;

    const payload = defaultValuesNewOrganization(tableViewName ?? '');

    organizations.create({
      ...payload,
      name: search || 'Unnamed',
    });

    store.ui.commandMenu.toggle('AddNewOrganization');
  };

  useKeyBindings(
    {
      // Enter: handleConfirm,
      Escape: () => (store.ui.commandMenu.isOpen = false),
    },
    { when: allowSubmit },
  );

  return (
    <Command label={`Rename `} shouldFilter={false}>
      <div className='p-6 pt-4 pb-4 flex justify-between items-center gap-1 '>
        <p className='font-semibold'>Search 300,000+ organizations</p>
        <CommandCancelIconButton onClose={handleClose} />
      </div>

      <div className='p-6 pt-0 pb-4 flex flex-col gap-3'>
        <div className='flex items-center'>
          <Input
            value={search}
            className='text-base p-0 !border-none '
            placeholder='Search by name, website or LinkedIn...'
            onChange={(e) => handleSearchChange(e.target.value)}
          />

          <Spinner
            size='sm'
            label='Loading...'
            className='text-gray-300 fill-gray-500'
          />
        </div>

        <Command.List className='!px-0 !pb-0 border-t border-gray-200'>
          {isLoading && (
            <div className='text-gray-400 flex gap-2 mb-2 ml-4'>
              <Spinner
                size='sm'
                label='Loading...'
                className='text-gray-300 fill-gray-500 size-3'
              />
              Searching...
            </div>
          )}

          {options.length === 0 && !isLoading && search.length && (
            <Command.Item
              onSelect={() => {
                const tableViewDef = store.tableViewDefs.getById(preset ?? '');

                handleAddNewOrganization();
                tableViewDef?.setSorting(
                  ColumnViewType.OrganizationsUpdatedDate,
                  true,
                );
              }}
            >
              <div className='flex items-center gap-2 overflow-hidden'>
                <PlusCircle className='text-primary-700  ' />

                <span className='truncate '>Add {search}</span>
              </div>
            </Command.Item>
          )}

          {options.length > 0 &&
            options.map((option, index) => (
              <Command.Item
                key={option?.id || `option-${index}`}
                onSelect={() => {
                  const tableViewDef = store.tableViewDefs.getById(
                    preset ?? '',
                  );

                  organizations.addOrganizationByGlobalOrgId(option.id);
                  tableViewDef?.setSorting(
                    ColumnViewType.OrganizationsUpdatedDate,
                    true,
                  );
                }}
              >
                <div className='flex items-center gap-2 overflow-hidden'>
                  <Avatar
                    size='xxs'
                    textSize='xxs'
                    name={option.name}
                    className='ml-[1px]'
                    variant='outlineSquare'
                    icon={<User03 className='text-primary-700  ' />}
                    src={
                      option?.logoUrl
                        ? option.logoUrl
                        : option?.iconUrl || undefined
                    }
                  />

                  <span className='truncate '>{option?.name || 'Unnamed'}</span>
                  <span className='whitespace-nowrap truncate'>•</span>
                  <span className='whitespace-nowrap truncate'>
                    {getFormattedLink(option.website)}
                  </span>
                </div>
              </Command.Item>
            ))}
        </Command.List>
      </div>
    </Command>
  );
});

const defaultValuesNewOrganization = (organizationName: string) => {
  switch (organizationName) {
    case 'Customers':
      return {
        relationship: OrganizationRelationship.Customer,
        stage: OrganizationStage.Onboarding,
      };
    case 'Leads':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Lead,
      };
    case 'Nurture':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Target,
      };
    case 'All orgs':
      return {
        relationship: OrganizationRelationship.Prospect,
        stage: OrganizationStage.Target,
      };

    case 'Churn':
      return {
        relationship: OrganizationRelationship.FormerCustomer,
        stage: OrganizationStage.PendingChurn,
      };
    default:
      return {};
  }
};

const getFormattedLink = (url: string): string => {
  return url.replace(/^(https?:\/\/)?(www\.)?([^/?#]+).*/i, '$3');
};

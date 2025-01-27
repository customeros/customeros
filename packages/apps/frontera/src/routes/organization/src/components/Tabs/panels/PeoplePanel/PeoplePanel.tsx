import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { differenceInCalendarMonths } from 'date-fns';
import { SearchSortContact } from '@domain/usecases/contact-details/search-sort-contacts.usecase';

import { Input } from '@ui/form/Input';
import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Users02 } from '@ui/media/icons/Users02';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { UsersPlus } from '@ui/media/icons/UsersPlus';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { ContactDetails } from '@shared/components/ContactDetails';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

import { SortOptionsMenu } from './components/SortOptionsMenu';
const searchSortContactUseCase = new SearchSortContact();

export const PeoplePanel = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const [expandAll, setExpandAll] = useState(false);
  const organization = store.organizations.getById(id);
  const contacts = organization?.contacts;

  const searchSortContact = searchSortContactUseCase;

  const search = searchSortContact.getSearch();
  const sortDir = searchSortContact.getSortDirection();
  const sortBy = searchSortContact.getSort();

  const filteredContactsState =
    contacts?.filter((v) => {
      if (!search) return true;

      return (
        v.name.toLowerCase().includes(search.toLowerCase()) ||
        v.firstName.toLowerCase().includes(search.toLowerCase()) ||
        (v.primaryOrganizationJobRoleTitle || '')
          .toLowerCase()
          .includes(search.toLowerCase())
      );
    }).length === 0 && search;

  useEffect(() => {
    if (organization?.value.contacts.length) {
      store.contacts.preload(organization?.value.contacts);
    }
  }, []);

  if (!contacts) return null;

  return (
    <OrganizationPanel
      withFade
      title='People'
      isLoading={store.contacts.isLoading}
      bgImage={
        !contacts?.length
          ? '/backgrounds/organization/half-circle-pattern.svg'
          : ''
      }
      actionItem={
        !!contacts.length && (
          <IconButton
            size='xs'
            variant='outline'
            aria-label='Add contact'
            className='text-gray-500'
            dataTest={'org-people-add-contact'}
            icon={<UsersPlus className='text-gray-500' />}
            spinner={
              <Spinner
                size='sm'
                label='adding'
                className='text-gray-300 fill-gray-400'
              />
            }
            onClick={() => {
              store.ui.commandMenu.setType('AddSingleContact');
              store.ui.commandMenu.setContext({
                ids: [id],
                entity: 'Organization',
                property: 'contacts',
              });
              store.ui.commandMenu.setOpen(true);
            }}
          >
            Add
          </IconButton>
        )
      }
    >
      {!contacts.length && (
        <div className='flex flex-col items-center mt-4'>
          <div className='border-1 border-gray-200 p-3 rounded-md mb-6 mt-5'>
            <FeaturedIcon colorScheme='gray'>
              <Users02 className='text-gray-700 size-6' />
            </FeaturedIcon>
          </div>
          <span className='text-gray-700 font-medium'>Assemble the team</span>
          <span className='text-gray-700 mt-1 mb-6 text-center text-sm'>
            Start by adding people that work at {organization?.value.name}, and
            keep track of everyone from decision-makers to day-to-day
            collaborators.
          </span>
          <div>
            <Button
              variant='outline'
              loadingText='Adding'
              colorScheme={'primary'}
              dataTest='org-people-add-someone'
              isDisabled={store.contacts.isLoading}
              onClick={() => {
                store.ui.commandMenu.setType('AddSingleContact');
                store.ui.commandMenu.setContext({
                  ids: [id],
                  entity: 'Organization',
                  property: 'contacts',
                });
                store.ui.commandMenu.setOpen(true);
              }}
            >
              Add someone
            </Button>
          </div>
        </div>
      )}

      {contacts.length > 0 && (
        <div className='flex items-center justify-between'>
          <div className='flex items-center gap-2'>
            <SearchSm className='text-gray-500 size-4' />
            <Input
              size='xs'
              type='text'
              value={search}
              variant='unstyled'
              placeholder='Search name or title...'
              onChange={(e) => searchSortContact.setSearch(e.target.value)}
            />
          </div>
          <div className='flex items-center'>
            <SortOptionsMenu searchSortContact={searchSortContact} />
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='collapse/decollapse'
              onClick={() => setExpandAll(!expandAll)}
              icon={!expandAll ? <ChevronExpand /> : <ChevronCollapse />}
            />
          </div>
        </div>
      )}

      {/* Filtered contacts */}
      {contacts
        .sort((a, b) => {
          if (sortBy === 'First name') {
            if (sortDir === 'asc') {
              return a.name.localeCompare(b.name);
            } else {
              return b.name.localeCompare(a.name);
            }
          }

          if (sortBy === 'Created') {
            const dateA = new Date(a.createdAt);
            const dateB = new Date(b.createdAt);

            if (sortDir === 'asc') {
              return dateA.getTime() - dateB.getTime();
            } else {
              return dateB.getTime() - dateA.getTime();
            }
          }

          if (sortBy === 'Updated') {
            const dateA = new Date(a.updatedAt);
            const dateB = new Date(b.updatedAt);

            if (sortDir === 'asc') {
              return dateA.getTime() - dateB.getTime();
            } else {
              return dateB.getTime() - dateA.getTime();
            }
          }

          if (sortBy === 'Tenure') {
            const aTenure = Math.abs(
              differenceInCalendarMonths(
                new Date(a.primaryOrganizationJobRoleStartDate),
                new Date(),
              ),
            );

            const bTenure = Math.abs(
              differenceInCalendarMonths(
                new Date(b.primaryOrganizationJobRoleStartDate),
                new Date(),
              ),
            );

            if (sortDir === 'asc') {
              return aTenure - bTenure;
            } else {
              return bTenure - aTenure;
            }
          }

          return 0;
        })
        .filter((v) => {
          if (!search) return true;

          return (
            v.name.toLowerCase().includes(search.toLowerCase()) ||
            v.firstName.toLowerCase().includes(search.toLowerCase()) ||
            (v.primaryOrganizationJobRoleTitle || '')
              .toLowerCase()
              .includes(search.toLowerCase())
          );
        })
        .map((contact) => (
          <div
            key={contact?.id}
            className='group/card'
            style={{ width: '100%' }}
          >
            <ContactDetails
              isExpandble={true}
              expandAll={expandAll}
              id={contact?.id ?? ''}
            />
          </div>
        ))}

      {filteredContactsState && (
        <div className='text-center text-gray-500 mt-4 text-sm'>
          No matches found—looks like a ghost town in here
        </div>
      )}
    </OrganizationPanel>
  );
});

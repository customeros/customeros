import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
<<<<<<< HEAD
import { SearchSortContact } from '@domain/usecases/people-contact-card/search-sort-contacts.usecase';
=======
<<<<<<< Updated upstream
>>>>>>> eb554d103 (contact card issues)

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { Users03 } from '@ui/media/icons/Users03';
import { useStore } from '@shared/hooks/useStore';
<<<<<<< HEAD
import { SearchSm } from '@ui/media/icons/SearchSm';
=======
=======
import { differenceInCalendarMonths } from 'date-fns';
import { SearchSortContact } from '@domain/usecases/people-contact-card/search-sort-contacts.usecase';

import { Input } from '@ui/form/Input';
import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Users02 } from '@ui/media/icons/Users02';
import { SearchSm } from '@ui/media/icons/SearchSm';
>>>>>>> Stashed changes
>>>>>>> eb554d103 (contact card issues)
import { UsersPlus } from '@ui/media/icons/UsersPlus';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

import { SortOptionsMenu } from './components/SortOptionsMenu';
import { ContactCard } from './components/ContactCard/ContactCard';
import { CreateNewContactModal } from './components/CreateNewContactModal';
const searchSortContactUseCase = new SearchSortContact();

export const PeoplePanel = observer(() => {
  const store = useStore();
  const { open, onOpen, onClose } = useDisclosure();
  const id = useParams()?.id as string;
  const organization = store.organizations.getById(id);
  const contacts = store.organizations.getById(id)?.contacts;

  if (!contacts) return null;

  const searchSortContact = searchSortContactUseCase;

  const search = searchSortContact.getSearch();
  const sortDir = searchSortContact.getSortDirection();
  const sortBy = searchSortContact.getSort();

  const filteredContactsState =
    contacts.filter((v) => {
      if (!search) return true;

      return (
        v.name.toLowerCase().includes(search.toLowerCase()) ||
        v.firstName.toLowerCase().includes(search.toLowerCase()) ||
        (v.primaryOrganizationJobRoleTitle || '')
          .toLowerCase()
          .includes(search.toLowerCase())
      );
    }).length === 0 && search;

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
          <Button
            size='xs'
            variant='outline'
            onClick={() => onOpen()}
            aria-label='Add contact'
            leftIcon={<UsersPlus />}
            dataTest={'org-people-add-contact'}
          >
            Add
          </Button>
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
              onClick={() => onOpen()}
              dataTest='org-people-add-someone'
              isDisabled={store.contacts.isLoading}
            >
              Add someone
            </Button>
          </div>
        </div>
      )}
<<<<<<< HEAD

      {contacts.length > 0 && (
        <div className='flex items-center justify-between'>
          <div className='flex items-center gap-2'>
            <SearchSm className='text-gray-500 size-4' />
            <Input
              type='text'
              value={search}
              variant='unstyled'
              placeholder='Search name or title...'
              onChange={(e) => searchSortContact.setSearch(e.target.value)}
            />
          </div>
          <SortOptionsMenu searchSortContact={searchSortContact} />
        </div>
=======
<<<<<<< Updated upstream
      {contacts.map((contact) => (
        <div key={contact?.id} className='group/card' style={{ width: '100%' }}>
          <ContactCard id={contact?.id} />
        </div>
      ))}
=======

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
          <SortOptionsMenu searchSortContact={searchSortContact} />
        </div>
>>>>>>> eb554d103 (contact card issues)
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
<<<<<<< HEAD
            const aTenure =
              a.primaryOrganizationJobRoleStartDate -
              a.primaryOrganizationJobRoleEndDate;
            const bTenure =
              b.primaryOrganizationJobRoleStartDate -
              b.primaryOrganizationJobRoleEndDate;
=======
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
>>>>>>> eb554d103 (contact card issues)

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
            <ContactCard id={contact?.id} />
          </div>
        ))}

      {filteredContactsState && (
        <div className='text-center text-gray-500 mt-4 text-sm'>
          No matches found—looks like a ghost town in here
        </div>
      )}
<<<<<<< HEAD
=======
>>>>>>> Stashed changes
>>>>>>> eb554d103 (contact card issues)

      <CreateNewContactModal orgId={id} open={open} onClose={onClose} />
    </OrganizationPanel>
  );
});

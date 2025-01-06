import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { SearchSortContact } from '@domain/usecases/people-contact-card/search-sort-contacts';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { Users03 } from '@ui/media/icons/Users03';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { UsersPlus } from '@ui/media/icons/UsersPlus';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

import { SortOptionsMenu } from './components/SortOptionsMenu';
import { ContactCard } from './components/ContactCard/ContactCard';
import { CreateNewContactModal } from './components/CreateNewContactModal';
const searchSortContactUseCase = SearchSortContact.getInstance();

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
          <div className='border-1 border-gray-200 p-3 rounded-md mb-6'>
            <Users03 className='text-gray-700 size-6' />
          </div>
          <span className='text-gray-700 font-semibold'>
            Let’s add some people
          </span>
          <span className='text-gray-500 mt-1 mb-6 text-center'>
            With the right people, you&apos;ll create meaningful interactions
            and results. Start by adding yourself, your colleagues or anyone
            from {organization?.value.name}.
          </span>
          <div>
            <Button
              variant='outline'
              loadingText='Adding'
              onClick={() => onOpen()}
              dataTest='org-people-add-someone'
              isDisabled={store.contacts.isLoading}
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
              type='text'
              value={search}
              variant='unstyled'
              placeholder='Search name or title...'
              onChange={(e) => searchSortContact.setSearch(e.target.value)}
            />
          </div>
          <SortOptionsMenu searchSortContact={searchSortContact} />
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
            const aTenure =
              a.primaryOrganizationJobRoleStartDate -
              a.primaryOrganizationJobRoleEndDate;
            const bTenure =
              b.primaryOrganizationJobRoleStartDate -
              b.primaryOrganizationJobRoleEndDate;

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

      <CreateNewContactModal orgId={id} open={open} onClose={onClose} />
    </OrganizationPanel>
  );
});

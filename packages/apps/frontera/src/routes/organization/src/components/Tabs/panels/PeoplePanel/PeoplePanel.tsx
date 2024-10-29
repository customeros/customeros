import { MouseEvent } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Contact } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Users03 } from '@ui/media/icons/Users03';
import { useStore } from '@shared/hooks/useStore';
import { UsersPlus } from '@ui/media/icons/UsersPlus';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { ContactCard } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/ContactCard';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

export const PeoplePanel = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const organization = store.organizations.value.get(id);

  const contacts =
    organization?.contacts.slice().sort((a, b) => {
      return a?.createdAt > b?.createdAt ? -1 : 1;
    }) ?? [];

  const handleAddContact = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    e.stopPropagation();

    const organizationId = organization?.id;

    if (!organizationId) return;

    store.contacts.create(organizationId);
  };

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
            onClick={handleAddContact}
            dataTest={'org-people-add-contact'}
            icon={<UsersPlus className='text-gray-500' />}
            spinner={
              <Spinner
                size='sm'
                label='adding'
                className='text-gray-300 fill-gray-400'
              />
            }
          >
            Add
          </IconButton>
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
              onClick={handleAddContact}
              dataTest='org-people-add-someone'
              isDisabled={store.contacts.isLoading}
              rightSpinner={
                <Spinner
                  size='sm'
                  label='adding'
                  className='text-gray-300 fill-gray-400'
                />
              }
            >
              Add someone
            </Button>
          </div>
        </div>
      )}
      {!!contacts.length &&
        contacts.map((contact) => (
          <div key={contact.metadata.id} style={{ width: '100%' }}>
            <ContactCard
              id={contact.metadata.id}
              contact={contact as Contact}
              organizationName={organization?.value.name}
            />
          </div>
        ))}
    </OrganizationPanel>
  );
});

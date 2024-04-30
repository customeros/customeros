'use client';
import { MouseEvent } from 'react';
import { useParams } from 'react-router-dom';

import { useQueryClient } from '@tanstack/react-query';

import { Contact } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { Users03 } from '@ui/media/icons/Users03';
import { UsersPlus } from '@ui/media/icons/UsersPlus';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { invalidateQuery } from '@organization/components/Tabs/panels/PeoplePanel/util';
import { useCreateContactMutation } from '@organization/graphql/createContact.generated';
import { ContactCard } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/ContactCard';
import { useOrganizationPeoplePanelQuery } from '@organization/graphql/organizationPeoplePanel.generated';
import { PeoplePanelSkeleton } from '@organization/components/Tabs/panels/PeoplePanel/PeoplePanelSkeleton';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { useAddOrganizationToContactMutation } from '@organization/graphql/addContactToOrganization.generated';

export const PeoplePanel = () => {
  const id = useParams()?.id as string;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data, isLoading: isInitialLoading } = useOrganizationPeoplePanelQuery(
    client,
    {
      id,
    },
    {
      staleTime: 1000,
      refetchOnWindowFocus: 'always',
      refetchOnReconnect: 'always',
      refetchOnMount: 'always',
    },
  );
  const createContact = useCreateContactMutation(client);
  const addContactToOrganization = useAddOrganizationToContactMutation(client, {
    onSuccess: () => invalidateQuery(queryClient, id),
  });
  const isLoading =
    createContact.isPending || addContactToOrganization.isPending;

  const contacts = data?.organization?.contacts.content.map((c) => c) ?? [];

  const handleAddContact = (e: Event & MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    e.stopPropagation();
    createContact.mutate(
      { input: {} },
      {
        onSuccess: (data) => {
          const contactId = data.contact_Create.id;
          addContactToOrganization.mutate({
            input: { contactId, organizationId: id },
          });
        },
      },
    );
  };

  if (isInitialLoading) {
    return <PeoplePanelSkeleton />;
  }

  return (
    <OrganizationPanel
      title='People'
      withFade
      bgImage={
        !contacts?.length
          ? '/backgrounds/organization/half-circle-pattern.svg'
          : ''
      }
      actionItem={
        !!contacts.length && (
          <Button
            size='sm'
            variant='outline'
            className='text-gray-500'
            loadingText='Adding'
            isLoading={isLoading}
            spinner={
              <Spinner
                className='text-gray-300 fill-gray-400'
                size='sm'
                label='adding'
              />
            }
            onClick={handleAddContact}
            leftIcon={<UsersPlus className='text-gray-500' />}
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
            from {data?.organization?.name}.
          </span>
          <div>
            <Button
              variant='outline'
              loadingText='Adding'
              isLoading={isLoading}
              onClick={handleAddContact}
              spinner={
                <Spinner
                  className='text-gray-300 fill-gray-400'
                  size='sm'
                  label='adding'
                />
              }
            >
              Add someone
            </Button>
          </div>
        </div>
      )}
      {!!contacts.length &&
        contacts.map((contact, index) => (
          <div key={contact.id} style={{ width: '100%' }}>
            <ContactCard
              contact={contact as Contact}
              organizationName={data?.organization?.name}
            />
          </div>
        ))}
    </OrganizationPanel>
  );
};

'use client';
import { MouseEvent } from 'react';
import { useParams } from 'next/navigation';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Box } from '@ui/layout/Box';
import { Text } from '@ui/typography/Text';
import { Button } from '@ui/form/Button';
import { Icons } from '@ui/media/Icon';
import { Fade } from '@ui/transitions/Fade';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCreateContactMutation } from '@organization/graphql/createContact.generated';
import { useAddOrganizationToContactMutation } from '@organization/graphql/addContactToOrganization.generated';

import { ContactCard } from '@organization/components/Tabs/panels/PeoplePanel/ContactCard/ContactCard';
import { useOrganizationPeoplePanelQuery } from '@organization/graphql/organizationPeoplePanel.generated';
import { invalidateQuery } from '@organization/components/Tabs/panels/PeoplePanel/util';
import { Contact } from '@graphql/types';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { PeoplePanelSkeleton } from '@organization/components/Tabs/panels/PeoplePanel/PeoplePanelSkeleton';

export const PeoplePanel = () => {
  const id = useParams()?.id as string;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data, isInitialLoading } = useOrganizationPeoplePanelQuery(client, {
    id,
  });
  const createContact = useCreateContactMutation(client);
  const addContactToOrganization = useAddOrganizationToContactMutation(client, {
    onSuccess: () => invalidateQuery(queryClient, id),
  });
  const isLoading =
    createContact.isLoading || addContactToOrganization.isLoading;

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
            loadingText='Adding'
            isLoading={isLoading}
            onClick={handleAddContact}
            leftIcon={<Icons.UsersPlus />}
            type='button'
          >
            Add
          </Button>
        )
      }
    >
      {!contacts.length && (
        <Flex direction='column' alignItems='center' mt='4'>
          <Box
            border='1px solid'
            borderColor='gray.200'
            padding={3}
            borderRadius='md'
            mb={6}
          >
            <Icons.Users2 color='gray.700' boxSize='6' />
          </Box>
          <Text color='gray.700' fontWeight={600}>
            Let’s add some people
          </Text>
          <Text color='gray.500' mt={1} mb={6} textAlign='center'>
            With the right people, you&apos;ll create meaningful interactions
            and results. Start by adding yourself, your colleagues or anyone
            from {data?.organization?.name}.
          </Text>
          <div>
            <Button
              variant='outline'
              loadingText='Adding'
              isLoading={isLoading}
              onClick={handleAddContact}
            >
              Add someone
            </Button>
          </div>
        </Flex>
      )}

      {!!contacts.length &&
        contacts.map((contact, index) => (
          <Fade key={contact.id} in style={{ width: '100%' }}>
            <ContactCard index={index} contact={contact as Contact} />
          </Fade>
        ))}
    </OrganizationPanel>
  );
};

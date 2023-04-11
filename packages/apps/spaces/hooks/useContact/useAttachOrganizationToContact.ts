import {
  useAttachOrganizationToContactMutation,
  ContactOrganizationInput,
  GetContactCommunicationChannelsQuery,
} from './types';
import {
  AttachOrganizationToContactMutation,
  GetContactCommunicationChannelsDocument,
  GetContactPersonalDetailsDocument,
  GetContactPersonalDetailsWithOrganizationsDocument,
} from '../../graphQL/__generated__/generated';
import { gql, useApolloClient } from '@apollo/client';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';

interface Result {
  onAttachOrganizationToContact: (
    input: ContactOrganizationInput,
  ) => Promise<
    AttachOrganizationToContactMutation['contact_AddOrganizationById'] | null
  >;
}
export const useAttachOrganizationToContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const client = useApolloClient();
  const [attachOrganizationToContactMutation, { loading, error, data }] =
    useAttachOrganizationToContactMutation();

  const handleAttachOrganizationToContact: Result['onAttachOrganizationToContact'] =
    async (contactOrg) => {
      try {
        const response = await attachOrganizationToContactMutation({
          variables: {
            input: {
              contactId,
              organizationId: contactOrg.organizationId,
            },
          },
          refetchQueries: ['useGetContactPersonalDetailsWithOrganizations'],
        });

        return response.data?.contact_AddOrganizationById ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onAttachOrganizationToContact: handleAttachOrganizationToContact,
  };
};

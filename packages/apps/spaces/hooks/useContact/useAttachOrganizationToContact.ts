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
          variables: { input: contactOrg },
        });
        // const data = client.cache.readQuery({
        //   query: GetContactPersonalDetailsWithOrganizationsDocument,
        //   variables: { id: contactOrg.contactId },
        // });
        //
        // const newData = data?.contact?.organizations?.content?.map(
        //   (data, i) => i === index,
        // );
        // console.log('🏷️ ----- data: ', data);
        //
        // client.cache.writeFragment({
        //   id: `Contact:${contactId}`,
        //   fragment: gql`
        //     fragment organizationsInContact on Contact {
        //       id
        //       organizations
        //     }
        //   `,
        //   data: {
        //     // @ts-expect-error revisit
        //     ...data.contact,
        //     organizations: [
        //       // @ts-expect-error revisit
        //       ...data.contact.organizations,
        //       { ...response.data?.contact_AddOrganizationById },
        //     ],
        //   },
        // });
        // Update the cache with the new object
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

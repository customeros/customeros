import {
  AddTagToContactMutation,
  ContactTagInput,
  useAddTagToContactMutation,
} from './types';
import {
  GetContactPersonalDetailsWithOrganizationsDocument,
  GetContactTagsDocument,
  GetContactTimelineQuery,
} from '../../graphQL/__generated__/generated';
import { gql, useApolloClient } from '@apollo/client';
import { ApolloCache } from 'apollo-cache';

interface Result {
  onAddTagToContact: (
    input: ContactTagInput,
  ) => Promise<AddTagToContactMutation['contact_AddTagById'] | null>;
}
export const useAddTagToContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const client = useApolloClient();
  const [addTagToContactMutation, { loading, error, data }] =
    useAddTagToContactMutation();
  const handleUpdateCacheAfterTag = (
    cache: ApolloCache<any>,
    { data: { contact_AddTagById } }: any,
  ) => {
    console.log('🏷️ ----- contact_AddTagById: ', contact_AddTagById);
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTagsDocument,
      variables: {
        id: contactId,
      },
    });
    // @ts-expect-error fix function type
    const normalizedId = cache.identify({
      id: contactId,
      __typename: 'Contact',
    });
    const contactData = client.readFragment({
      id: normalizedId,
      fragment: gql`
        fragment ContactName on Contact {
          id
          name
          firstName
          lastName
          tags {
            id
            name
          }
        }
      `,
    });

    const contactWithNewJobrole = {
      ...contactData,
      jobRoles: [...contactData.jobRoles, newrole],
    };
    console.log('🏷️ ----- data: ', data);
    if (data === null) {
      client.writeQuery({
        query: GetContactPersonalDetailsWithOrganizationsDocument,
        data: {
          contact: {
            id: contactId,
            tags,
          },
          variables: { id: contactId },
        },
      });
    }

    const newData = {
      contact: {
        ...data,
        ...contactWithNewJobrole,
      },
    };
    client.cache.writeFragment({
      id: `Contact:${contactId}`,
      fragment: gql`
        fragment Tags on Contact {
          id
          tags
        }
      `,
      data: {
        tags: [...contact_AddTagById.tags],
      },
    });
    client.writeQuery({
      query: GetContactPersonalDetailsWithOrganizationsDocument,
      data: newData,
      variables: {
        id: contactId,
      },
    });
  };
  const handleAddTagToContact: Result['onAddTagToContact'] = async (
    contactTagInput,
  ) => {
    try {
      const response = await addTagToContactMutation({
        variables: { input: contactTagInput },
      });

      if (response?.data?.contact_AddTagById?.tags) {
        const data = client.cache.readQuery({
          query: GetContactTagsDocument,
          variables: { id: contactTagInput.contactId },
        });

        client.cache.writeFragment({
          id: `Contact:${contactId}`,
          fragment: gql`
            fragment Tags on Contact {
              id
              tags
            }
          `,
          data: {
            tags: [...response.data.contact_AddTagById.tags],
          },
        });
      }

      // Update the cache with the new object
      return response.data?.contact_AddTagById ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onAddTagToContact: handleAddTagToContact,
  };
};

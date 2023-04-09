import {
  CreateContactJobRoleMutation,
  GetContactPersonalDetailsWithOrganizationsDocument,
  GetContactTimelineQuery,
  JobRoleInput,
  useCreateContactJobRoleMutation,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';

interface Props {
  contactId: string;
}

interface Result {
  onCreateContactJobRole: (
    input: JobRoleInput,
  ) => Promise<CreateContactJobRoleMutation['jobRole_Create'] | null>;
}
export const useCreateContactJobRole = ({ contactId }: Props): Result => {
  const [createContactJobRoleMutation, { loading, error, data }] =
    useCreateContactJobRoleMutation();

  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: { jobRole_Create } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactPersonalDetailsWithOrganizationsDocument,
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
          source
          jobRoles {
            id
            jobTitle
            organization {
              id
              name
            }
          }
          organizations(pagination: { limit: 99999, page: 1 }) {
            content {
              id
              name
            }
          }
        }
      `,
    });

    const newrole = {
      id: '',
      jobTitle: '',
      organization: {
        id: '',
        name: '',
      },
      ...jobRole_Create,
    };
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
            ...contactWithNewJobrole,
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

    client.writeQuery({
      query: GetContactPersonalDetailsWithOrganizationsDocument,
      data: newData,
      variables: {
        id: contactId,
      },
    });
  };

  const handleCreateContactJobRole: Result['onCreateContactJobRole'] = async (
    jobRole,
  ) => {
    try {
      const response = await createContactJobRoleMutation({
        variables: { contactId, input: jobRole },
        refetchQueries: ['GetContactPersonalDetailsWithOrganizations'],
        // @ts-expect-error this should not result in error, debug later
        update: handleUpdateCacheAfterAddingNote,
      });
      return response.data?.jobRole_Create ?? null;
    } catch (err) {
      return null;
    }
  };

  return {
    onCreateContactJobRole: handleCreateContactJobRole,
  };
};

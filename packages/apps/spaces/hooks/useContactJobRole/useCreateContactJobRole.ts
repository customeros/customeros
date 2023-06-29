import {
  CreateContactJobRoleMutation, GetContactDocument,
  GetContactPersonalDetailsWithOrganizationsDocument, GetContactQuery,
  GetContactTimelineQuery,
  JobRoleInput,
  useCreateContactJobRoleMutation,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

interface Result {
  saving: boolean;
  onCreateContactJobRole: (
    input: JobRoleInput,
  ) => Promise<CreateContactJobRoleMutation['jobRole_Create'] | null>;
}
export const useCreateContactJobRole = ({ contactId }: Props): Result => {
  const [createContactJobRoleMutation, { loading, error, data }] =
    useCreateContactJobRoleMutation();

  const handleUpdateCacheAfterAddingJobRole = (
    cache: ApolloCache<any>,
    { data: { jobRole_Create } }: any,
  ) => {
    const data: GetContactQuery | null = client.readQuery({
      query: GetContactDocument,
      variables: {
        id: contactId,
      },
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

    if (data === null) {
      client.writeQuery({
        query: GetContactDocument,
        data: {
          contact: {
            id: contactId,
            jobRoles: [newrole],

          },
          variables: { id: contactId },
        },
      });
    }

    const newData = {
      contact: {
        ...data?.contact,
        jobRoles: [...(data?.contact?.jobRoles || []), newrole],
      },
    };

    client.writeQuery({
      query: GetContactDocument,
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
        update: handleUpdateCacheAfterAddingJobRole,
      });
      return response.data?.jobRole_Create ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong while adding new job role! Please contact us or try again later',
        {
          toastId: `create-contact-job-role-error--${contactId}`,
        },
      );
      return null;
    }
  };

  return {
    onCreateContactJobRole: handleCreateContactJobRole,
    saving: loading,
  };
};

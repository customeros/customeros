import {
  JobRoleUpdateInput,
  UpdateJobRoleMutation,
  useUpdateJobRoleMutation,
} from './types';
import {
  GetContactDocument,
  GetContactPersonalDetailsWithOrganizationsDocument
} from '../../graphQL/__generated__/generated';
import { gql, useApolloClient } from '@apollo/client';
import { toast } from 'react-toastify';

interface Result {
  onUpdateContactJobRole: (
    input: JobRoleUpdateInput,
  ) => Promise<UpdateJobRoleMutation['jobRole_Update'] | null>;
}
export const useUpdateContactJobRole = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [updateContactJobRoleMutation, { loading, error, data }] =
    useUpdateJobRoleMutation();
  const client = useApolloClient();

  const handleUpdateContactJobRole: Result['onUpdateContactJobRole'] = async (
    input,
  ) => {
    try {
      const response = await updateContactJobRoleMutation({
        variables: { input, contactId },
      });
      const data = client.cache.readQuery({
        query: GetContactDocument,
        variables: { id: contactId },
      });

      if (response.data?.jobRole_Update.primary) {
        // @ts-expect-error revisit
        const updatedJobRoles = data.contact.jobRoles.map((role) =>
          role.id !== response.data?.jobRole_Update.id
            ? {
                ...role,
                primary: response.data?.jobRole_Update.primary
                  ? false
                  : role.primary,
              }
            : { ...role, ...response.data?.jobRole_Update },
        );

        client.cache.writeFragment({
          id: `Contact:${contactId}`,
          fragment: gql`
            fragment JobRoles on Contact {
              id
              jobRoles {
                id
                jobTitle
                primary
                organization {
                  id
                  name
                }
              }
            }
          `,
          data: {
            // @ts-expect-error revisit
            ...(data?.contact ?? {}),
            jobRoles: updatedJobRoles,
          },
        });
      }

      return response.data?.jobRole_Update ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong while updating job role! Please contact us or try again later',
        {
          toastId: `update-contact-job-role-error-${input.id}-${contactId}`,
        },
      );
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateContactJobRole: handleUpdateContactJobRole,
  };
};

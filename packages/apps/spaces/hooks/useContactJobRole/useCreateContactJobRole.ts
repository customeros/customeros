import {
  CreateContactJobRoleMutation,
  JobRoleInput,
  useCreateContactJobRoleMutation,
} from '../../graphQL/__generated__/generated';

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

  const handleCreateContactJobRole: Result['onCreateContactJobRole'] = async (
    jobRole,
  ) => {
    try {
      const response = await createContactJobRoleMutation({
        variables: { contactId, input: jobRole },
        refetchQueries: ['GetContactPersonalDetails'],
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

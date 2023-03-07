import {
  JobRoleUpdateInput,
  UpdateJobRoleMutation,
  useUpdateJobRoleMutation,
} from './types';

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

  const handleUpdateContactJobRole: Result['onUpdateContactJobRole'] = async (
    input,
  ) => {
    console.log('🏷️ ----- input: ', input);
    try {
      const response = await updateContactJobRoleMutation({
        variables: { input, contactId },
        refetchQueries: ['GetContactPersonalDetails'],
      });

      return response.data?.jobRole_Update ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateContactJobRole: handleUpdateContactJobRole,
  };
};

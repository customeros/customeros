import {
  useRemoveContactJobRoleMutation,
  RemoveContactJobRoleMutation,
} from './types';

interface Props {
  contactId: string;
}
interface Result {
  onRemoveContactJobRole: (
    jobRoleId: string,
  ) => Promise<RemoveContactJobRoleMutation['jobRole_Delete'] | null>;
}
export const useRemoveJobRoleFromContactJobRole = ({
  contactId,
}: Props): Result => {
  const [removeJobRoleFromContactMutation, { loading, error, data }] =
    useRemoveContactJobRoleMutation();

  const handleRemoveJobRoleFromContact: Result['onRemoveContactJobRole'] =
    async (roleId) => {
      try {
        const response = await removeJobRoleFromContactMutation({
          variables: { contactId, roleId },
          refetchQueries: ['GetContactPersonalDetails'],
        });
        return response.data?.jobRole_Delete ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onRemoveContactJobRole: handleRemoveJobRoleFromContact,
  };
};

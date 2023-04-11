import {
  useRemoveContactJobRoleMutation,
  RemoveContactJobRoleMutation,
} from './types';
import { toast } from 'react-toastify';

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
          update(cache) {
            const normalizedId = cache.identify({
              id: roleId,
              __typename: 'JobRole',
            });
            cache.evict({ id: normalizedId });
            cache.gc();
          },
        });
        return response.data?.jobRole_Delete ?? null;
      } catch (err) {
        toast.error(
          'Something went wrong while deleting job role. Please contact us or try again later',
          {
            toastId: `contact-${roleId}-delete-error`,
          },
        );
        return null;
      }
    };

  return {
    onRemoveContactJobRole: handleRemoveJobRoleFromContact,
  };
};

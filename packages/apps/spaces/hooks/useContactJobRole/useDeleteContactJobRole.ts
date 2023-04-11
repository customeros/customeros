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
          refetchQueries: ['useGetContactPersonalDetailsWithOrganizations'],
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
        console.error(err);
        return null;
      }
    };

  return {
    onRemoveContactJobRole: handleRemoveJobRoleFromContact,
  };
};

import { ContactOrganizationInput } from './types';
import {
  RemoveOrganizationFromContactMutation,
  useRemoveOrganizationFromContactMutation,
} from '../../graphQL/__generated__/generated';

interface Result {
  onRemoveOrganizationFromContact: (
    input: ContactOrganizationInput,
  ) => Promise<
    | RemoveOrganizationFromContactMutation['contact_RemoveOrganizationById']
    | null
  >;
}
export const useRemoveOrganizationFromContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [removeOrganizationFromContactMutation, { loading, error, data }] =
    useRemoveOrganizationFromContactMutation();

  const handleRemoveOrganizationFromContact: Result['onRemoveOrganizationFromContact'] =
    async (contactOrg) => {
      try {
        const response = await removeOrganizationFromContactMutation({
          variables: {
            input: { contactId, organizationId: contactOrg.organizationId },
          },
          awaitRefetchQueries: true,
          refetchQueries: ['useGetContactPersonalDetailsWithOrganizations'],
        });
        return response.data?.contact_RemoveOrganizationById ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onRemoveOrganizationFromContact: handleRemoveOrganizationFromContact,
  };
};

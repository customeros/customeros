import { ContactOrganizationInput } from './types';
import {
  RemoveOrganizationFromContactMutation,
  useRemoveOrganizationFromContactMutation,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

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
          update(cache) {
            const normalizedId = cache.identify({
              id: contactOrg.organizationId,
              __typename: 'Organization',
            });
            cache.evict({ id: normalizedId });
            cache.gc();
          },
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

import {
  ContactUpdateInput,
  UpdateContactPersonalDetailsMutation,
  useUpdateContactPersonalDetailsMutation,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

interface Result {
  onUpdateContactPersonalDetails: (
    input: Omit<ContactUpdateInput, 'id'>,
  ) => Promise<UpdateContactPersonalDetailsMutation['contact_Update'] | null>;
}
export const useUpdateContactPersonalDetails = ({
  contactId,
}: Props): Result => {
  const [updateContactPersonalDetails, { loading, error, data }] =
    useUpdateContactPersonalDetailsMutation();

  const handleUpdateContactPersonalDetails: Result['onUpdateContactPersonalDetails'] =
    async (input) => {
      try {
        const response = await updateContactPersonalDetails({
          variables: { input: { ...input, id: contactId } },
        });
        return response.data?.contact_Update ?? null;
      } catch (err) {
        console.error(err);
        toast.error(
          'Something went wrong while updating contact personal details. Please contact us or try again later',
          {
            toastId: `personal-details-${contactId}-update-error`,
          },
        );
        return null;
      }
    };

  return {
    onUpdateContactPersonalDetails: handleUpdateContactPersonalDetails,
  };
};

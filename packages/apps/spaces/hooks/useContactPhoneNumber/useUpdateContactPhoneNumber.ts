import { PhoneNumber, PhoneNumberUpdateInput } from '../../graphQL/generated';
import {
  UpdateContactPhoneNumberMutation,
  useUpdateContactPhoneNumberMutation,
} from '../../graphQL/generated';

interface Result {
  onUpdateContactPhoneNumber: (
    input: Omit<PhoneNumberUpdateInput, 'id'>,
    oldValue: PhoneNumber,
  ) => Promise<
    UpdateContactPhoneNumberMutation['phoneNumberUpdateInContact'] | null
  >;
}
export const useUpdateContactPhoneNumber = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [updateContactNoteMutation, { loading, error, data }] =
    useUpdateContactPhoneNumberMutation();

  const handleUpdateContactPhoneNumber: Result['onUpdateContactPhoneNumber'] =
    async (input, { label, primary = false, id, ...rest }) => {
      const payload = {
        primary,
        label,
        id,
        ...input,
      };
      try {
        const response = await updateContactNoteMutation({
          variables: { input: payload, contactId },
          refetchQueries: ['GetContactCommunicationChannels'],
          optimisticResponse: {
            phoneNumberUpdateInContact: {
              __typename: 'PhoneNumber',
              ...rest,
              ...payload,
              primary: input.primary || primary || false,
            },
          },
        });

        return response.data?.phoneNumberUpdateInContact ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateContactPhoneNumber: handleUpdateContactPhoneNumber,
  };
};

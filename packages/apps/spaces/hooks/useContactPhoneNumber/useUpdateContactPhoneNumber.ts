import {
  PhoneNumber,
  PhoneNumberUpdateInput,
} from '../../graphQL/__generated__/generated';
import {
  UpdateContactPhoneNumberMutation,
  useUpdateContactPhoneNumberMutation,
} from '../../graphQL/__generated__/generated';

interface Result {
  onUpdateContactPhoneNumber: (
    input: PhoneNumberUpdateInput,
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
    async (input) => {
      const payload = {
        ...input,
      };
      try {
        const response = await updateContactNoteMutation({
          variables: { input: payload, contactId },
          optimisticResponse: {
            phoneNumberUpdateInContact: {
              __typename: 'PhoneNumber',
              ...payload,
              e164: payload.phoneNumber,
              rawPhoneNumber: payload.phoneNumber,
              primary: input.primary || false,
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

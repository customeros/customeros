import {
  AddPhoneToContactMutation,
  useAddPhoneToContactMutation,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

interface Result {
  onCreateContactPhoneNumber: (
    input: any, //FIXME
  ) => Promise<AddPhoneToContactMutation['phoneNumberMergeToContact'] | null>;
  loading: boolean;
}
export const useCreateContactPhoneNumber = ({ contactId }: Props): Result => {
  const [createContactPhoneNumberMutation, { loading, error, data }] =
    useAddPhoneToContactMutation();

  const handleCreateContactPhoneNumber: Result['onCreateContactPhoneNumber'] =
    async (input) => {
      try {
        const response = await createContactPhoneNumberMutation({
          variables: { contactId, input },
          refetchQueries: ['GetContactCommunicationChannels'],
          awaitRefetchQueries: true,
          optimisticResponse: {
            phoneNumberMergeToContact: {
              __typename: 'PhoneNumber',
              ...input,
              id: 'optimistic-id',
              e164: input?.phoneNumber || '',
              rawPhoneNumber: input?.phoneNumber || '',
              primary: input?.primary || false,
            },
          },
        });
        return response.data?.phoneNumberMergeToContact ?? null;
      } catch (err) {
        toast.error(
          'Something went wrong while adding phone number. Please contact us or try again later',
          {
            toastId: `contact-phone-number-${contactId}-add-error`,
          },
        );
        console.error(err);
        return null;
      }
    };

  return {
    onCreateContactPhoneNumber: handleCreateContactPhoneNumber,
    loading,
  };
};

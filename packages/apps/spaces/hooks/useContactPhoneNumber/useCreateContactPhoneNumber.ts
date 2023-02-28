import {
  AddPhoneToContactMutation,
  useAddPhoneToContactMutation,
} from '../../graphQL/__generated__/generated';
import { PhoneNumberInput } from '../../graphQL/__generated__/generated';

interface Props {
  contactId: string;
}

interface Result {
  onCreateContactPhoneNumber: (
    input: any, //FIXME
  ) => Promise<AddPhoneToContactMutation['phoneNumberMergeToContact'] | null>;
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
        });
        return response.data?.phoneNumberMergeToContact ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onCreateContactPhoneNumber: handleCreateContactPhoneNumber,
  };
};

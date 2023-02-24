import {
  RemovePhoneNumberFromContactMutation,
  useRemovePhoneNumberFromContactMutation,
} from '../../graphQL/generated';

interface Result {
  onRemovePhoneNumberFromContact: (
    phoneNumberId: string,
  ) => Promise<
    | RemovePhoneNumberFromContactMutation['phoneNumberRemoveFromContactById']
    | null
  >;
}
export const useRemovePhoneNumberFromContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [removePhoneNumberFromContactMutation, { loading, error, data }] =
    useRemovePhoneNumberFromContactMutation();

  const handleRemovePhoneNumberFromContact = async (id: string) => {
    try {
      const response = await removePhoneNumberFromContactMutation({
        variables: { id: id, contactId },
        refetchQueries: ['GetContactCommunicationChannels'],
      });
      return response.data?.phoneNumberRemoveFromContactById ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onRemovePhoneNumberFromContact: handleRemovePhoneNumberFromContact,
  };
};

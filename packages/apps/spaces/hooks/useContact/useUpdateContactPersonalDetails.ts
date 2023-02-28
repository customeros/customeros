import {
  ContactUpdateInput,
  UpdateContactPersonalDetailsMutation,
  useUpdateContactPersonalDetailsMutation,
} from '../../graphQL/__generated__/generated';

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
          // optimisticResponse: {
          //   emailUpdateInContact: {
          //     __typename: 'Email',
          //     ...input,
          //     primary: input.primary || false,
          //   },
          // },
        });
        return response.data?.contact_Update ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateContactPersonalDetails: handleUpdateContactPersonalDetails,
  };
};

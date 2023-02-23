import {
  ContactInput,
  CreateContactMutation,
  useCreateContactMutation,
} from './types';

interface Props {
  contact?: ContactInput;
}

interface Result {
  onCreateEmptyContact: () => Promise<string | null>;
  onCreateContact: (
    input: ContactInput,
  ) => Promise<CreateContactMutation['contact_Create'] | null>;
}
export const useCreateContact = ({ contact }: Props): Result => {
  const [createContactMutation, { loading, error, data }] =
    useCreateContactMutation();

  const handleCreateContact: Result['onCreateContact'] = async (contact) => {
    try {
      const optimisticItem = { id: 'optimistic-id', ...contact };
      const response = await createContactMutation({
        variables: { input: contact },
        optimisticResponse: {
          contact_Create: {
            __typename: 'Contact',
            ...optimisticItem,
          },
        },
      });
      return response.data?.contact_Create ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  const handleCreateEmptyContact: Result['onCreateEmptyContact'] = async () => {
    try {
      const response = await createContactMutation({
        variables: { input: {} },
      });
      return response.data?.contact_Create.id ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onCreateEmptyContact: handleCreateEmptyContact,
    onCreateContact: handleCreateContact,
  };
};

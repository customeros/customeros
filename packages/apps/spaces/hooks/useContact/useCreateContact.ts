import {
  ContactInput,
  CreateContactMutation,
  useCreateContactMutation,
} from './types';
import { toast } from 'react-toastify';

interface Props {
  contact?: ContactInput;
}

interface Result {
  onCreateEmptyContact: () => Promise<string | null>;
  onCreateContact: (
    input: ContactInput,
  ) => Promise<CreateContactMutation['contact_Create'] | null>;
}
export const useCreateContact = (): Result => {
  const [createContactMutation, { loading, error, data }] =
    useCreateContactMutation();

  const handleCreateContact: Result['onCreateContact'] = async (contact) => {
    try {
      const optimisticItem = { id: 'optimistic-id', ...contact };
      const response = await createContactMutation({
        variables: { input: contact },
        refetchQueries: ['GetDashboardData'],
      });

      toast.success(
        `${contact?.firstName} ${contact?.lastName} was added to contacts`,
      );
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

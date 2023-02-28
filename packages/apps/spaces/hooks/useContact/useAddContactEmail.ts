import { useAddEmailToContactMutation } from './types';
import { EmailInput } from '../../graphQL/__generated__/generated';
import { AddEmailToContactMutation } from '../../graphQL/__generated__/generated';

interface Props {
  email?: EmailInput;
}

interface Result {
  onAddEmailToContact: (
    input: EmailInput,
  ) => Promise<AddEmailToContactMutation['emailMergeToContact'] | null>;
}
export const useAddEmailToContactEmail = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [addEmailToContactMutation, { loading, error, data }] =
    useAddEmailToContactMutation();

  const handleAddEmailToContact: Result['onAddEmailToContact'] = async (
    email,
  ) => {
    try {
      const optimisticItem = { id: 'optimistic-id', ...email };
      const response = await addEmailToContactMutation({
        variables: { contactId, input: email },
        refetchQueries: ['GetContactCommunicationChannels'],
        optimisticResponse: {
          emailMergeToContact: {
            __typename: 'Email',
            ...optimisticItem,
            primary: optimisticItem?.primary || false,
          },
        },
      });
      return response.data?.emailMergeToContact ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onAddEmailToContact: handleAddEmailToContact,
  };
};

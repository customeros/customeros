import { useRemoveEmailFromContactMutation } from './types';
import { RemoveEmailFromContactMutation } from '../../graphQL/generated';

interface Props {
  contactId: string;
}

interface Result {
  onRemoveEmailFromContact: (
    emailId: string,
  ) => Promise<
    RemoveEmailFromContactMutation['emailRemoveFromContactById'] | null
  >;
}
export const useRemoveEmailFromContactEmail = ({
  contactId,
}: Props): Result => {
  const [removeEmailFromContactMutation, { loading, error, data }] =
    useRemoveEmailFromContactMutation();

  const handleRemoveEmailFromContact: Result['onRemoveEmailFromContact'] =
    async (emailId) => {
      try {
        const response = await removeEmailFromContactMutation({
          variables: { contactId, id: emailId },
          refetchQueries: ['GetContactCommunicationChannels'],
          optimisticResponse: {
            emailRemoveFromContactById: {
              __typename: 'Result',
              result: true,
            },
          },
        });
        return response.data?.emailRemoveFromContactById ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onRemoveEmailFromContact: handleRemoveEmailFromContact,
  };
};

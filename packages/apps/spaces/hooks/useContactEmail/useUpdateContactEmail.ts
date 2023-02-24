import { Email, EmailUpdateInput } from '../../graphQL/generated';
import {
  UpdateContactEmailMutation,
  useUpdateContactEmailMutation,
} from '../../graphQL/generated';

interface Result {
  onUpdateContactEmail: (
    input: EmailUpdateInput,
    oldValue: Email,
  ) => Promise<UpdateContactEmailMutation['emailUpdateInContact'] | null>;
}
export const useUpdateContactEmail = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [updateContactNoteMutation, { loading, error, data }] =
    useUpdateContactEmailMutation();

  const handleUpdateContactEmail: Result['onUpdateContactEmail'] = async (
    input,
    { label, primary = false, id, ...rest },
  ) => {
    const payload = {
      primary,
      label,
      // @ts-expect-error revisit later, shouldn't happen
      id,
      ...input,
    };
    try {
      const response = await updateContactNoteMutation({
        variables: { input: payload, contactId },
        refetchQueries: ['GetContactCommunicationChannels'],
        optimisticResponse: {
          emailUpdateInContact: {
            __typename: 'Email',
            ...rest,
            ...input,
            primary: input.primary || primary || false,
          },
        },
      });

      return response.data?.emailUpdateInContact ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateContactEmail: handleUpdateContactEmail,
  };
};

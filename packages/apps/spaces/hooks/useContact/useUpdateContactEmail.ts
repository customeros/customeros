import {
  Email,
  UpdateContactEmailMutation,
  useUpdateContactEmailMutation,
} from './types';
import { EmailUpdateInput } from '../../graphQL/types';

interface Props {
  contactId: string;
}

interface Result {
  onUpdateContactEmail: (
    input: EmailUpdateInput,
  ) => Promise<UpdateContactEmailMutation['emailUpdateInContact'] | null>;
}
export const useUpdateContactEmail = ({ contactId }: Props): Result => {
  const [updateContactEmailMutation, { loading, error, data }] =
    useUpdateContactEmailMutation();

  const handleUpdateContactEmail: Result['onUpdateContactEmail'] = async (
    input,
  ) => {
    try {
      console.log('🏷️ ----- input: ', input);
      const response = await updateContactEmailMutation({
        variables: { input: { ...input }, contactId },
        // optimisticResponse: {
        //   emailUpdateInContact: {
        //     __typename: 'Email',
        //     ...input,
        //     primary: input.primary || false,
        //   },
        // },
      });

      console.log('🏷️ ----- response: ', response);
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

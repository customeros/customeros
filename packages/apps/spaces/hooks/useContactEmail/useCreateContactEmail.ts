import { NoteInput } from '../../graphQL/generated';
import {
  CreateContactNoteMutation,
  useCreateContactNoteMutation,
} from '../../graphQL/generated';

interface Props {
  contactId: string;
}

interface Result {
  onCreateContactNote: (
    input: NoteInput,
  ) => Promise<CreateContactNoteMutation['note_CreateForContact'] | null>;
}
export const useCreateContactEmail = ({ contactId }: Props): Result => {
  const [createContactNoteMutation, { loading, error, data }] =
    useCreateContactNoteMutation();

  const handleCreateContactNote: Result['onCreateContactNote'] = async (
    note,
  ) => {
    try {
      const response = await createContactNoteMutation({
        variables: { contactId, input: note },
      });
      return response.data?.note_CreateForContact ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onCreateContactNote: handleCreateContactNote,
  };
};

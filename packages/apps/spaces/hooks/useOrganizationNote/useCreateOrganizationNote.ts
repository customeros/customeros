import {
  NoteInput,
  CreateOrganizationNoteMutation,
  useCreateOrganizationNoteMutation,
} from './types';

interface Props {
  organizationId: string;
}

interface Result {
  onCreateOrganizationNote: (
    input: NoteInput,
  ) => Promise<
    CreateOrganizationNoteMutation['note_CreateForOrganization'] | null
  >;
}
export const useCreateOrganizationNote = ({
  organizationId,
}: Props): Result => {
  const [createOrganizationNoteMutation, { loading, error, data }] =
    useCreateOrganizationNoteMutation();

  const handleCreateOrganizationNote: Result['onCreateOrganizationNote'] =
    async (note) => {
      try {
        const response = await createOrganizationNoteMutation({
          variables: { organizationId, input: note },
        });
        return response.data?.note_CreateForOrganization ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onCreateOrganizationNote: handleCreateOrganizationNote,
  };
};

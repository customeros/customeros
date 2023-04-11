import {
  UpdateOrganizationEmailMutation,
  useUpdateOrganizationEmailMutation,
  EmailUpdateInput,
} from './types';

interface Props {
  organizationId: string;
}

interface Result {
  onUpdateOrganizationEmail: (
    input: EmailUpdateInput,
  ) => Promise<
    UpdateOrganizationEmailMutation['emailUpdateInOrganization'] | null
  >;
}
export const useUpdateOrganizationEmail = ({
  organizationId,
}: Props): Result => {
  const [updateOrganizationEmailMutation, { loading, error, data }] =
    useUpdateOrganizationEmailMutation();

  const handleUpdateOrganizationEmail: Result['onUpdateOrganizationEmail'] =
    async (input) => {
      try {
        const response = await updateOrganizationEmailMutation({
          variables: { input: { ...input }, organizationId },
        });

        return response.data?.emailUpdateInOrganization ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateOrganizationEmail: handleUpdateOrganizationEmail,
  };
};

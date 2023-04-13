import {
  UpdateOrganizationEmailMutation,
  useUpdateOrganizationEmailMutation,
  EmailUpdateInput,
} from './types';
import { ApolloCache } from 'apollo-cache';
import {
  GetOrganizationCommunicationChannelsQuery,
  GetOrganizationCommunicationChannelsDocument,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';

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
  const handleUpdateCacheAfterAddingEmail = (
    cache: ApolloCache<any>,
    { data: { emailUpdateInOrganization } }: any,
  ) => {
    const data: GetOrganizationCommunicationChannelsQuery | null =
      client.readQuery({
        query: GetOrganizationCommunicationChannelsDocument,
        variables: {
          id: organizationId,
        },
      });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationCommunicationChannelsDocument,
        variables: {
          id: organizationId,
        },
        data: {
          organization: {
            id: organizationId,
            emails: [emailUpdateInOrganization],
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        emails: (data.organization?.emails || []).map((e) =>
          e.id === emailUpdateInOrganization.id
            ? { ...e, ...emailUpdateInOrganization }
            : {
                ...e,
                primary: emailUpdateInOrganization.primary ? false : e.primary,
              },
        ),
      },
    };
    client.writeQuery({
      query: GetOrganizationCommunicationChannelsDocument,
      data: newData,
      variables: {
        id: organizationId,
      },
    });
  };
  const handleUpdateOrganizationEmail: Result['onUpdateOrganizationEmail'] =
    async (input) => {
      try {
        const response = await updateOrganizationEmailMutation({
          variables: { input: { ...input }, organizationId },
          //@ts-expect-error fixme
          update: handleUpdateCacheAfterAddingEmail,
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

import {
  ContactUpdateInput,
  GetContactPersonalDetailsDocument,
  GetContactTimelineQuery,
  UpdateContactPersonalDetailsMutation,
  useUpdateContactPersonalDetailsMutation,
} from '../../graphQL/__generated__/generated';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';

interface Props {
  contactId: string;
}

interface Result {
  onUpdateContactPersonalDetails: (
    input: Omit<ContactUpdateInput, 'id'>,
  ) => Promise<UpdateContactPersonalDetailsMutation['contact_Update'] | null>;
}
export const useUpdateContactPersonalDetails = ({
  contactId,
}: Props): Result => {
  const [updateContactPersonalDetails, { loading, error, data }] =
    useUpdateContactPersonalDetailsMutation();
  const handleUpdateCacheAfterAddingNote = (
    cache: ApolloCache<any>,
    { data: data1 }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactPersonalDetailsDocument,
      variables: {
        contactId,
      },
    });
    // @ts-expect-error fix function type
    const normalizedId = cache.identify({
      id: contactId,
      __typename: 'Contact',
    });
    const contactData = client.readFragment({
      id: normalizedId,
      fragment: gql`
        fragment ContactPersonalDetailsFragment on Contact {
          ...ContactPersonalDetails
        }
      `,
    });
    // const newNoteWithNoted = {
    //   ...note_CreateForContact,
    //   noted: [
    //     {
    //       ...contactData,
    //     },
    //   ],
    // };
    // if (data === null) {
    //   client.writeQuery({
    //     query: GetContactPersonalDetailsDocument,
    //     data: {
    //       contact: {
    //         contactId,
    //         timelineEvents: [newNoteWithNoted],
    //       },
    //       variables: { contactId, from: NOW_DATE, size: 10 },
    //     },
    //   });
    //   return;
    // }
    //
    // const newData = {
    //   contact: {
    //     ...data.contact,
    //     timelineEvents: [newNoteWithNoted],
    //   },
    // };
    //
    // client.writeQuery({
    //   query: GetContactPersonalDetailsDocument,
    //   data: newData,
    //   variables: {
    //     contactId,
    //     from: NOW_DATE,
    //     size: 10,
    //   },
    // });
  };

  const handleUpdateContactPersonalDetails: Result['onUpdateContactPersonalDetails'] =
    async (input) => {
      try {
        const response = await updateContactPersonalDetails({
          variables: { input: { ...input, id: contactId } },
        });
        return response.data?.contact_Update ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateContactPersonalDetails: handleUpdateContactPersonalDetails,
  };
};

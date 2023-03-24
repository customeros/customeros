import {
  CreatePhoneCallInteractionEventMutation, GetContactTimelineDocument, GetContactTimelineQuery,
  useCreatePhoneCallInteractionEventMutation,
} from '../../graphQL/__generated__/generated';
import {ApolloCache} from "apollo-cache";
import client from "../../apollo-client";
import {gql} from "@apollo/client";
import {toast} from "react-toastify";

interface Props {
  contactId: string;
}

interface Result {
  onCreatePhoneCallInteractionEvent: (
    input: any,
  ) => Promise<
    CreatePhoneCallInteractionEventMutation['interactionEvent_Create'] | null
  >;
}

const NOW_DATE = new Date().toISOString();

export const useCreatePhoneCallInteractionEvent = ({
  contactId,
}: Props): Result => {
  const [createPhoneCallInteractionEvent, { loading, error, data }] =
    useCreatePhoneCallInteractionEventMutation();

  const handleUpdateCacheAfterAddingPhoneCall = (
      cache: ApolloCache<any>,
      { data: { interactionEvent_Create } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTimelineDocument,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [interactionEvent_Create],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [
          interactionEvent_Create,
        ],
      },
    };

    client.writeQuery({
      query: GetContactTimelineDocument,
      data: newData,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });
  };

  const handleCreatePhoneCallInteractionEvent: Result['onCreatePhoneCallInteractionEvent'] = async (
      input,
  ) => {
    try {
      const response = await createPhoneCallInteractionEvent({
        variables: {
          contactId: contactId,
          content: input.content,
          contentType: input.contentType,
          sentBy: input.sentBy,
        },
        // @ts-expect-error this should not result in error, debug later
        update: handleUpdateCacheAfterAddingPhoneCall,
      });
      if (response.data) {
        toast.success('Phone call log added!', {
          toastId: `phone-call-added-${response.data?.interactionEvent_Create.id}`,
        });
      }
      return response.data?.interactionEvent_Create ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding a phone call', {
        toastId: `phone-call-add-error-${contactId}`,
      });
      return null;
    }
  };

  return {
    onCreatePhoneCallInteractionEvent: handleCreatePhoneCallInteractionEvent,
  };
};

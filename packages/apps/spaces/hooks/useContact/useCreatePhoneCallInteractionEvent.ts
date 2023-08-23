import {
  DataSource,
  GetContactTimelineDocument,
  GetContactTimelineQuery,
  useCreatePhoneCallInteractionEventMutation,
} from '@spaces/graphql';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
  onSuccess: (data: any) => void;
}

interface Result {
  onCreatePhoneCallInteractionEvent: (input: any) => void;
}

const NOW_DATE = new Date().toISOString();

export const useCreatePhoneCallInteractionEvent = ({
  contactId,
  onSuccess,
}: Props): Result => {
  const [createPhoneCallInteractionEvent] =
    useCreatePhoneCallInteractionEventMutation({
      onError: () => {
        toast.error('Something went wrong while adding a phone call', {
          toastId: `phone-call-add-error-${contactId}`,
        });
      },
      onCompleted: (res) => {
        onSuccess(res);
      },
    });

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

    const interactionEvent = {
      ...interactionEvent_Create,
      source: DataSource.Openline,
    };
    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [interactionEvent],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [interactionEvent],
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

  const handleCreatePhoneCallInteractionEvent: Result['onCreatePhoneCallInteractionEvent'] =
    (input) => {
      return createPhoneCallInteractionEvent({
        variables: {
          contactId: contactId,
          content: input.content,
          contentType: input.contentType,
          sentBy: input.sentBy,
        },
        update: handleUpdateCacheAfterAddingPhoneCall,
      });
    };

  return {
    onCreatePhoneCallInteractionEvent: handleCreatePhoneCallInteractionEvent,
  };
};

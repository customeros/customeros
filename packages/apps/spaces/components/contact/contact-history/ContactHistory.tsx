import React from "react";
import {useContactNotes} from "../../../hooks/useContactNote/useContactNotes";
import {Timeline} from "../../ui-kit/organisms";
import {useContactConversations} from "../../../hooks/useContactConversations";
import {useContactActions} from "../../../hooks/useContactActions";

export const ContactHistory = ({id}: { id: string }) => {
    const {data: notes, loading: notesLoading, error: notesError} = useContactNotes({id});
    // const {data: actions, loading: actionsLoading, error: actionsError} = useContactActions({id});
    // TODO add pagination support
    const {data: conversations, loading: conversationsLoading, error: conversationsError} = useContactConversations({id});

    console.log('in history')

    const noHistoryItemsAvailable = !notesLoading && notes?.notes?.content.length == 0 &&
        !conversationsLoading && conversations?.conversations?.content.length == 0;
        // !actionsLoading && actions?.actions.length == 0;

    const getSortedItems = (data1: Array<any> | undefined, data2: Array<any> | undefined, data3: Array<any> | undefined) => {
        var data = [...(data1 || []), ...(data2 || []), ...(data3 || [])];
        return data.sort((a, b) => {
            // @ts-ignore
            return Date.parse(a?.createdAt) - Date.parse(b?.createdAt);
        })
    }

    return (
        <Timeline
            loading={notesLoading || conversationsLoading}
            noActivity={noHistoryItemsAvailable}
            contactId={id}
            notifyChange={() => {
            }}
            loggedActivities={getSortedItems(notes?.notes.content, conversations?.conversations.content, [])}/>
            // loggedActivities={getSortedItems([], [], [])}/>

    );
}

export default ContactHistory

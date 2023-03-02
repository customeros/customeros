import React from 'react';
import {useContactNotes} from '../../../hooks/useContactNote/useContactNotes';
import {Timeline} from '../../ui-kit/organisms';
import {useContactConversations} from '../../../hooks/useContactConversations';
import {useContactTickets} from '../../../hooks/useContact/useContactTickets';

export const ContactHistory = ({id}: { id: string }) => {
    const {
        data: notes,
        loading: notesLoading,
        error: notesError,
    } = useContactNotes({id});
    const {
        data: tickets,
        loading: ticketsLoading,
        error: ticketsError,
    } = useContactTickets({id});
    // const {data: actions, loading: actionsLoading, error: actionsError} = useContactActions({id});
    // TODO add pagination support
    const {
        data: conversations,
        loading: conversationsLoading,
        error: conversationsError,
    } = useContactConversations({id});

    console.log('in history');

    const noHistoryItemsAvailable =
        !notesLoading &&
        notes?.notes?.content.length == 0 &&
        !conversationsLoading &&
        conversations?.conversations?.content.length == 0 &&
        !ticketsLoading &&
        tickets?.tickets?.length == 0;
    // !actionsLoading && actions?.actions.length == 0;

    const getSortedItems = (
        data1: Array<any> | undefined,
        data2: Array<any> | undefined,
        data3: Array<any> | undefined,
    ) => {
        const data = [...(data1 || []), ...(data2 || []), ...(data3 || [])];
        console.log(data);

        return data.sort((a, b) => {
            return Date.parse(a?.createdAt) - Date.parse(b?.createdAt);
        });
    };

    return (
        <Timeline
            loading={notesLoading || conversationsLoading || ticketsLoading}
            noActivity={noHistoryItemsAvailable}
            contactId={id}
            loggedActivities={getSortedItems(
                notes?.notes.content,
                conversations?.conversations.content,
                tickets?.tickets,
            )}
        />
    );
};

export default ContactHistory;

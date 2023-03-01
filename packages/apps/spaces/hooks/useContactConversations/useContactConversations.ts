import {ApolloError} from 'apollo-client';
import {GetContactConversationsQuery, Pagination, useGetContactConversationsQuery} from "../../graphQL/__generated__/generated";

interface Props {
    id: string;
}

interface Result {
    data: GetContactConversationsQuery['contact'] | null | undefined;
    loading: boolean;
    error: ApolloError | null;
}

export const useContactConversations = ({id}: Props): Result => {
    const {data, loading, error} = useGetContactConversationsQuery({
        variables: {id},
    });

    if (loading) {
        return {
            loading: true,
            error: null,
            data: null,
        };
    }

    if (error) {
        return {
            error,
            loading: false,
            data: null,
        };
    }

    console.log('data loaded for conversations')
    return {
        data: data?.contact ?? null,
        loading,
        error: null,
    };
};

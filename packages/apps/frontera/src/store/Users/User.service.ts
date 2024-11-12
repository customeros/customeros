import { gql } from 'graphql-request';
import { Transport } from '@store/transport';

import {
  Filter,
  UserPage,
  Pagination,
} from '@shared/types/__generated__/graphql.types';

class UserService {
  private static instance: UserService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): UserService {
    if (!UserService.instance) {
      UserService.instance = new UserService(transport);
    }

    return UserService.instance;
  }

  async getUsers(payload: USERS_QUERY_PAYLOAD): Promise<USERS_QUERY_RESPONSE> {
    return this.transport.graphql.request<
      USERS_QUERY_RESPONSE,
      USERS_QUERY_PAYLOAD
    >(USERS_QUERY, payload);
  }
}

type USERS_QUERY_PAYLOAD = {
  where?: Filter;
  pagination: Pagination;
};
type USERS_QUERY_RESPONSE = {
  users: UserPage;
};
const USERS_QUERY = gql`
  query getUsers($pagination: Pagination!, $where: Filter) {
    users(pagination: $pagination, where: $where) {
      content {
        id
        firstName
        lastName
        name
        profilePhotoUrl
        mailboxes
        bot
        timezone
        hasLinkedInToken
        emails {
          email
        }
      }
      totalElements
    }
  }
`;

export { UserService };

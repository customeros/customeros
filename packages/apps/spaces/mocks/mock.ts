import { setupServer } from 'msw/node';
import { gql } from 'graphql-tag';

const mocks = [
  {
    request: {
      query: gql`
        query GetUser($id: ID!) {
          user(id: $id) {
            id
            firstName
            lastName
          }
        }
      `,
      variables: {
        id: '1',
      },
    },
    result: {
      data: {
        user: {
          id: '123',
          firstName: 'John',
          lastName: 'John',
        },
      },
    },
  },
];
const server = setupServer();

export { server, mocks };

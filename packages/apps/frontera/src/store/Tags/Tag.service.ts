import { gql } from 'graphql-request';
import { Transport } from '@store/transport';

import { Tag, TagUpdateInput } from '@shared/types/__generated__/graphql.types';

class TagService {
  private static instance: TagService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): TagService {
    if (!TagService.instance) {
      TagService.instance = new TagService(transport);
    }

    return TagService.instance;
  }

  async updateTag(payload: UPDATE_TAG_PAYLOAD): Promise<UPDATE_TAG_RESPONSE> {
    return this.transport.graphql.request<
      UPDATE_TAG_RESPONSE,
      UPDATE_TAG_PAYLOAD
    >(UPDATE_TAG_MUTATION, payload);
  }

  async deleteTag(payload: DELETE_TAG_PAYLOAD): Promise<DELETE_TAG_RESPONSE> {
    return this.transport.graphql.request<
      DELETE_TAG_RESPONSE,
      DELETE_TAG_PAYLOAD
    >(DELETE_TAG_MUTATION, payload);
  }
}

type DELETE_TAG_PAYLOAD = {
  id: string;
};
type DELETE_TAG_RESPONSE = {
  tag_Delete: {
    result: boolean;
  };
};
const DELETE_TAG_MUTATION = gql`
  mutation deleteTag($id: ID!) {
    tag_Delete(id: $id) {
      result
    }
  }
`;

type UPDATE_TAG_PAYLOAD = {
  input: TagUpdateInput;
};

type UPDATE_TAG_RESPONSE = {
  data: Tag;
};

const UPDATE_TAG_MUTATION = gql`
  mutation UpdateTag($input: TagUpdateInput!) {
    updateTag(input: $input) {
      id
    }
  }
`;

export { TagService };

import { Transport } from '@store/transport';

import { EntityType } from '@graphql/types';

import GetTagsDocument from './getTags.graphql';
import CreateTagDocument from './createTag.graphql';
import UpdateTagDocument from './updateTag.graphql';
import DeleteTagDocument from './deleteTag.graphql';
import GetTagsByEntityTypeDocument from './getTagsByEntity.graphql';
import { GetTagsQuery, GetTagsQueryVariables } from './getTags.generated';
import {
  CreateTagMutation,
  CreateTagMutationVariables,
} from './createTag.generated';
import {
  DeleteTagMutation,
  DeleteTagMutationVariables,
} from './deleteTag.generated';
import {
  UpdateTagMutation,
  UpdateTagMutationVariables,
} from './updateTag.generated';
import {
  GetTagsByEntityTypeQuery,
  GetTagsByEntityTypeQueryVariables,
} from './getTagsByEntity.generated';

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

  async createTag(payload: CreateTagMutationVariables) {
    return this.transport.graphql.request<
      CreateTagMutation,
      CreateTagMutationVariables
    >(CreateTagDocument, payload);
  }

  async updateTag(payload: UpdateTagMutationVariables) {
    return this.transport.graphql.request<
      UpdateTagMutation,
      UpdateTagMutationVariables
    >(UpdateTagDocument, payload);
  }

  async deleteTag(payload: DeleteTagMutationVariables) {
    return this.transport.graphql.request<
      DeleteTagMutation,
      DeleteTagMutationVariables
    >(DeleteTagDocument, payload);
  }

  async getTags() {
    return this.transport.graphql.request<GetTagsQuery, GetTagsQueryVariables>(
      GetTagsDocument,
    );
  }

  async getTagsByEntityType(entityType: EntityType) {
    return this.transport.graphql.request<
      GetTagsByEntityTypeQuery,
      GetTagsByEntityTypeQueryVariables
    >(GetTagsByEntityTypeDocument, {
      entityType,
    });
  }
}

export { TagService };

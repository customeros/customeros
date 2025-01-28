import { Transport } from './../../transport';
import GetSkusDocument from './queries/skus.graphql';
import { SkusQuery } from './queries/skus.generated';
import SaveSkuDocument from './mutations/skuSave.graphql';
import ArchiveSkuDocument from './mutations/skuArchive.graphql';
import {
  SkuSaveMutation,
  SkuSaveMutationVariables,
} from './mutations/skuSave.generated.ts';
import {
  SkuArchiveMutation,
  SkuArchiveMutationVariables,
} from './mutations/skuArchive.generated.ts';

export class SkuRepository {
  static instance: SkuRepository | null = null;
  private transport = Transport.getInstance();

  public static getInstance() {
    if (!SkuRepository.instance) {
      SkuRepository.instance = new SkuRepository();
    }

    return SkuRepository.instance;
  }

  async getSkus() {
    return this.transport.graphql.request<SkusQuery>(GetSkusDocument);
  }

  async saveSku(payload: SkuSaveMutationVariables) {
    return this.transport.graphql.request<
      SkuSaveMutation,
      SkuSaveMutationVariables
    >(SaveSkuDocument, payload);
  }

  async archiveSku(payload: SkuArchiveMutationVariables) {
    return this.transport.graphql.request<
      SkuArchiveMutation,
      SkuArchiveMutationVariables
    >(ArchiveSkuDocument, payload);
  }
}

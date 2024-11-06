import { match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { CustomFieldsStore } from '@store/Settings/CustomFields.store';

import { CustomFieldTemplateInput } from '@shared/types/__generated__/graphql.types';

import GetCustomFieldTemplatesDocument from './customFieldTemplate.graphql';
import SaveCustomFieldTemplateDocument from './customFieldTemplate_save.graphql';
import DeleteCustomFieldTemplateDocument from './customFieldTemplate_delete.graphql';
import {
  GetCustomFieldTemplatesQuery,
  GetCustomFieldTemplatesQueryVariables,
} from './customFieldTemplate.generated';
import {
  SaveCustomFieldTemplateMutation,
  SaveCustomFieldTemplateMutationVariables,
} from './customFieldTemplate_save.generated';
import {
  DeleteCustomFieldTemplateMutation,
  DeleteCustomFieldTemplateMutationVariables,
} from './customFieldTemplate_delete.generated';

export class CustomFieldsService {
  private static instance: CustomFieldsService;

  private constructor(private transport: Transport) {}

  static getInstance(transport: Transport) {
    if (!CustomFieldsService.instance) {
      CustomFieldsService.instance = new CustomFieldsService(transport);
    }

    return CustomFieldsService.instance;
  }

  async getCustomFields() {
    return this.transport.graphql.request<
      GetCustomFieldTemplatesQuery,
      GetCustomFieldTemplatesQueryVariables
    >(GetCustomFieldTemplatesDocument);
  }

  async saveCustomField(payload: CustomFieldTemplateInput) {
    return this.transport.graphql.request<
      SaveCustomFieldTemplateMutation,
      SaveCustomFieldTemplateMutationVariables
    >(SaveCustomFieldTemplateDocument, { input: payload });
  }

  async deleteCustomField(id: string) {
    return this.transport.graphql.request<
      DeleteCustomFieldTemplateMutation,
      DeleteCustomFieldTemplateMutationVariables
    >(DeleteCustomFieldTemplateDocument, { id });
  }

  public async mutateOperation(
    operation: Operation,
    _store: CustomFieldsStore,
  ) {
    const diff = operation.diff?.[0];
    const path = diff?.path;
    const customFieldId = operation?.entityId;

    if (!operation.diff.length) {
      return;
    }

    if (!customFieldId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }

    return match(path).otherwise(async () => {
      const payload = makePayload<CustomFieldTemplateInput>(operation);

      return await this.saveCustomField({
        ...payload,
        id: customFieldId,
      });
    });
  }
}

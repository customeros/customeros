import { RootStore } from '@store/root';
import { Syncable } from '@store/syncable';
import { Transport } from '@store/transport';
import { action, override, makeObservable } from 'mobx';

import {
  CustomField,
  CustomFieldInput,
  CustomFieldDataType,
} from '@shared/types/__generated__/graphql.types';

import { CustomFieldsService } from './__service__/customFields/CustomFields.service';

export class CustomFieldStore extends Syncable<CustomField> {
  private service: CustomFieldsService;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: CustomField,
  ) {
    super(root, transport, data ?? getDefaultValue());
    this.service = CustomFieldsService.getInstance(transport);

    makeObservable<CustomFieldStore>(this, {
      id: override,
      save: override,
      getId: override,
      setId: override,
      invalidate: action,
      getCustomFields: action,
      getChannelName: override,
    });
  }

  getCustomFields() {
    return this.service.getCustomFields();
  }

  getId() {
    return this.value.id;
  }

  setId(id: string) {
    this.value.id = id;
  }

  getChannelName(): string {
    return 'customFields';
  }

  static getDefaultValue(): CustomFieldInput {
    return {
      id: crypto.randomUUID(),
      name: '',
      value: '',
      templateId: '',
      datatype: CustomFieldDataType.Text,
    };
  }
}

export const getDefaultValue = (): CustomFieldInput => ({
  id: crypto.randomUUID(),
  name: '',
  value: '',
  templateId: '',
  datatype: CustomFieldDataType.Text,
});

import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import { action, override, runInAction, makeObservable } from 'mobx';

import {
  CustomField,
  CustomFieldTemplateInput,
} from '@shared/types/__generated__/graphql.types';

import { CustomFieldStore } from './CustomField.store';
import { CustomFieldsService } from './__service__/customFields/CustomFields.service';

export class CustomFieldsStore extends SyncableGroup<
  CustomField,
  CustomFieldStore
> {
  private service: CustomFieldsService;

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, CustomFieldStore);
    this.service = CustomFieldsService.getInstance(transport);

    makeObservable<CustomFieldsStore>(this, {
      save: action,
      deleteCustomField: action,
      channelName: override,
    });
  }

  async bootstrap() {
    try {
      const { customFieldTemplate_List } = await this.service.getCustomFields();

      const customFields = customFieldTemplate_List.map((template) => ({
        id: template.id,
        name: template.name,
        value: '',
        createdAt: new Date(),
        updatedAt: new Date(),
        source: 'customerOs',
        template: template,
      }));

      this.load(customFields as CustomField[], {
        getId: (data: CustomField) => data.id,
      });

      runInAction(() => {
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  async save(payload: CustomFieldTemplateInput) {
    try {
      await this.service.saveCustomField({
        validValues: payload.validValues,
        name: payload.name,
        entityType: payload.entityType,
        type: payload.type,
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  async deleteCustomField(id: string) {
    try {
      await this.service.deleteCustomField(id);
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  get channelName() {
    return 'customFields';
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray<T extends CustomFieldStore>(
    compute: (arr: CustomFieldStore[]) => T[],
  ) {
    const arr = this.toArray();

    return compute(arr);
  }
}

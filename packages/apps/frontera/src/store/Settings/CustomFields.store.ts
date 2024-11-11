import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { SyncableGroup } from '@store/syncable-group';
import { action, override, runInAction, makeObservable } from 'mobx';

import {
  DataSource,
  EntityType,
  CustomField,
  CustomFieldDataType,
  CustomFieldTemplateType,
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
      const response = await this.service.saveCustomField({
        validValues: payload.validValues,
        name: payload.name,
        entityType: payload.entityType,
        type: payload.type,
      });

      const customField: CustomFieldTemplateInput = {
        id: response.customFieldTemplate_Save.id,
        name: payload.name,
        entityType: payload.entityType,
        type: payload?.type,
        validValues: payload.validValues,
      };

      runInAction(() => {
        this.load(
          [
            {
              id: customField.id || '',
              name: customField.name || '',
              value: '',
              createdAt: new Date(),
              updatedAt: new Date(),
              source: DataSource.Openline,
              datatype: CustomFieldDataType.Text,
              template: {
                name: customField.name || '',
                id: customField.id || '',
                validValues: customField.validValues || [],
                createdAt: new Date(),
                updatedAt: new Date(),
                entityType: customField.entityType || EntityType.Organization,
                type: customField.type || CustomFieldTemplateType.FreeText,
              },
            },
          ],
          {
            getId: (data: CustomField) => data.id,
          },
        );
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  async deleteCustomField(id: string) {
    const customField = this.value.get(id);

    try {
      await this.service.deleteCustomField(id);

      if (customField) {
        runInAction(() => {
          this.value.delete(id);
        });
      }
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  get channelName() {
    return 'customFields';
  }

  get persisterKey() {
    return 'CustomFields';
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

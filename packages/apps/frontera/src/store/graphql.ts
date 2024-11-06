import { match } from 'ts-pattern';

import type { Operation } from './types';
import type { Transport } from './transport';

import { RootStore } from './root';
import { OrganizationsService } from './Organizations/__service__/Organizations.service';
import { CustomFieldsService } from './Settings/__service__/customFields/CustomFields.service';

export class GraphqlService {
  private organizationsService: OrganizationsService;
  private customFieldsService: CustomFieldsService;

  constructor(private root: RootStore, private transport: Transport) {
    this.organizationsService = OrganizationsService.getInstance(
      this.transport,
    );
    this.customFieldsService = CustomFieldsService.getInstance(this.transport);
    this.getStore = this.getStore.bind(this);
  }

  public async mutate(operation: Operation) {
    if (!operation.entityId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }

    return match(operation.entity)
      .with('Organizations', async () => {
        const store = this.getStore(operation, 'organizations');

        if (!store) return;

        return await this.organizationsService.mutateOperation(
          operation,
          store,
        );
      })
      .with('customFields', async () => {
        const store = this.getStore(operation, 'customFields');

        if (!store) return;

        return await this.customFieldsService.mutateOperation(operation, store);
      })
      .otherwise(() => {});
  }

  private getStore(operation: Operation, storePath: keyof RootStore) {
    // @ts-expect-error no issue
    const store = this.root[storePath]?.value?.get(operation.entityId);

    if (!store) {
      console.error(
        `Store with id ${operation.entityId} not found. Mutations will not fire`,
      );

      return null;
    }

    return store;
  }
}

import omit from 'lodash/omit';
import { match } from 'ts-pattern';
import { Operation } from '@store/types';
import { Transport } from '@infra/transport';

import { TableViewDef } from '../TableViewDef.dto';
import GetTableViewDefsDocument from './getTableViewDefs.graphql';
import { TableViewDefsQuery } from './getTableViewDefs.generated';
import CreateTableViewDefDocument from './createTableViewDef.graphql';
import UpdateTableViewDefDocument from './updateTableViewDef.graphql';
import ArchiveTableViewDefDocument from './archiveTableViewDef.graphql';
import UpdateTableViewDefSharedDocument from './updateTableViewDefShared.graphql';
import {
  UpdateTableViewDefMutation,
  UpdateTableViewDefMutationVariables,
} from './updateTableViewDef.generated';
import {
  CreateTableViewDefMutation,
  CreateTableViewDefMutationVariables,
} from './createTableViewDef.generated';
import {
  ArchiveTableViewDefMutation,
  ArchiveTableViewDefMutationVariables,
} from './archiveTableViewDef.generated';
import {
  UpdateTableViewDefSharedMutation,
  UpdateTableViewDefSharedMutationVariables,
} from './updateTableViewDefShared.generated';

export class TableViewDefsService {
  private static instance: TableViewDefsService | null = null;
  private transport = Transport.getInstance();

  constructor() {}

  public static getInstance(): TableViewDefsService {
    if (!TableViewDefsService.instance) {
      TableViewDefsService.instance = new TableViewDefsService();
    }

    return TableViewDefsService.instance;
  }

  async getTableViewDefs(): Promise<TableViewDefsQuery> {
    return this.transport.graphql.request<TableViewDefsQuery>(
      GetTableViewDefsDocument,
    );
  }

  async createTableViewDef(
    variables: CreateTableViewDefMutationVariables,
  ): Promise<CreateTableViewDefMutation> {
    return this.transport.graphql.request<
      CreateTableViewDefMutation,
      CreateTableViewDefMutationVariables
    >(CreateTableViewDefDocument, variables);
  }

  async archiveTableViewDef(
    variables: ArchiveTableViewDefMutationVariables,
  ): Promise<ArchiveTableViewDefMutation> {
    return this.transport.graphql.request<
      ArchiveTableViewDefMutation,
      ArchiveTableViewDefMutationVariables
    >(ArchiveTableViewDefDocument, variables);
  }

  async updateTableViewDef(
    variables: UpdateTableViewDefMutationVariables,
  ): Promise<UpdateTableViewDefMutation> {
    return this.transport.graphql.request<
      UpdateTableViewDefMutation,
      UpdateTableViewDefMutationVariables
    >(UpdateTableViewDefDocument, variables);
  }

  async updateTableViewDefShared(
    variables: UpdateTableViewDefSharedMutationVariables,
  ): Promise<UpdateTableViewDefSharedMutation> {
    return this.transport.graphql.request<
      UpdateTableViewDefSharedMutation,
      UpdateTableViewDefSharedMutationVariables
    >(UpdateTableViewDefSharedDocument, variables);
  }

  public async mutateOperation(operation: Operation, store: TableViewDef) {
    if (!operation.diff.length) {
      return;
    }

    if (!operation.entityId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }

    const payload = {
      input: omit(
        store.value,
        'updatedAt',
        'createdAt',
        'tableType',
        'tableId',
        'isPreset',
        'isShared',
        'defaultFilters',
      ),
    };

    return match(store.value.isShared)
      .with(true, async () => this.updateTableViewDefShared(payload))
      .otherwise(async () => this.updateTableViewDef(payload));
  }
}

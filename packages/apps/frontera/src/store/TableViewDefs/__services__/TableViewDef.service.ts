import { Transport } from '@store/transport';
import GetTableViewDefs from '@store/TableViewDefs/__services__/getTableViewDefs.graphql';
import { TableViewDefsQuery } from '@store/TableViewDefs/__services__/getTableViewDefs.generated.ts';
import CreateTableViewDefDocument from '@store/TableViewDefs/__services__/createTableViewDef.graphql';
import UpdateTableViewDefDocument from '@store/TableViewDefs/__services__/updateTableViewDef.graphql';
import ArchiveTableViewDefDocument from '@store/TableViewDefs/__services__/archiveTableViewDef.graphql';
import {
  UpdateTableViewDefMutation,
  UpdateTableViewDefMutationVariables,
} from '@store/TableViewDefs/__services__/updateTableViewDef.generated.ts';
import {
  CreateTableViewDefMutation,
  CreateTableViewDefMutationVariables,
} from '@store/TableViewDefs/__services__/createTableViewDef.generated.ts';
import {
  ArchiveTableViewDefMutation,
  ArchiveTableViewDefMutationVariables,
} from '@store/TableViewDefs/__services__/archiveTableViewDef.generated.ts';

export class TableViewDefsService {
  private static instance: TableViewDefsService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  public static getInstance(transport: Transport): TableViewDefsService {
    if (!TableViewDefsService.instance) {
      TableViewDefsService.instance = new TableViewDefsService(transport);
    }

    return TableViewDefsService.instance;
  }

  async getTableViewDefs(): Promise<TableViewDefsQuery> {
    return this.transport.graphql.request<TableViewDefsQuery>(GetTableViewDefs);
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
}

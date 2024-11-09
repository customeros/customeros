import { match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';

import { InvoiceUpdateInput } from '@shared/types/__generated__/graphql.types';
import {
  GetInvoiceQuery,
  GetInvoiceQueryVariables,
} from '@shared/graphql/getInvoice.generated';

import { InvoiceStore } from '../Invoice.store';
import GetInvoiceDocument from './getInvoice.graphql';
import GetInvoicesDocument from './getInvoices.graphql';
import UpdateInvoiceStatusDocument from './updateInvoiceStatus.graphql';
import {
  GetInvoicesQuery,
  GetInvoicesQueryVariables,
} from './getInvoices.generated';
import {
  UpdateInvoiceStatusMutation,
  UpdateInvoiceStatusMutationVariables,
} from './updateInvoiceStatus.generated';

export class InvoicesService {
  private static instance: InvoicesService;
  private transport: Transport;

  private constructor(transport: Transport) {
    this.transport = transport;
  }

  public static getInstance(transport: Transport) {
    if (!InvoicesService.instance) {
      InvoicesService.instance = new InvoicesService(transport);
    }

    return InvoicesService.instance;
  }

  async getInvoice(invoiceNumber: string) {
    return this.transport.graphql.request<
      GetInvoiceQuery,
      GetInvoiceQueryVariables
    >(GetInvoiceDocument, { id: invoiceNumber });
  }

  async getInvoices(payload: GetInvoicesQueryVariables) {
    return this.transport.graphql.request<
      GetInvoicesQuery,
      GetInvoicesQueryVariables
    >(GetInvoicesDocument, payload);
  }

  async updateInvoiceStatus(payload: UpdateInvoiceStatusMutationVariables) {
    return this.transport.graphql.request<
      UpdateInvoiceStatusMutation,
      UpdateInvoiceStatusMutationVariables
    >(UpdateInvoiceStatusDocument, {
      input: {
        ...payload.input,
        id: payload.input.id,
        patch: true,
      },
    });
  }

  public async mutateOperation(operation: Operation, store: InvoiceStore) {
    const diff = operation.diff?.[0];
    const path = diff?.path;
    const invoiceNumber = operation.entityId;
    const invoiceId = store.value.metadata.id;

    if (!operation.diff.length) {
      return;
    }

    if (!invoiceNumber) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }
    match(path)
      .with(['status'], () => {
        const payload = makePayload<InvoiceUpdateInput>(operation);

        this.updateInvoiceStatus({ input: { ...payload, id: invoiceId } });
      })

      .otherwise(() => {});
  }
}

import { UTCDate } from '@date-fns/utc';

import {
  BilledType,
  ServiceLineItem,
  ServiceLineItemBulkUpdateItem,
} from '@graphql/types';

export interface BulkUpdateServiceLineItem {
  type?: string;
  name?: string;
  price: number;
  quantity: number;
  isDeleted: boolean;
  vatRate?: number | null;
  billed?: BilledType | null;
  serviceStarted?: Date | null;
  serviceLineItemId?: string | null;
}

export class ServiceLineItemsDTO implements ServiceLineItemBulkUpdateItem {
  name?: string;
  quantity: number;
  price: number;

  billed?: BilledType | null;
  type?: string | null;
  isDeleted?: boolean | null;
  serviceLineItemId?: string | null;
  serviceStarted?: Date | null;
  vatRate?: number | null;

  constructor(data?: ServiceLineItem | null) {
    this.serviceLineItemId = data?.metadata?.id ?? '';
    this.name = data?.description ?? '';
    this.quantity = data?.quantity ?? 1;
    this.price = data?.price ?? 0;
    this.billed = data?.billingCycle ?? BilledType.Monthly;
    this.isDeleted = data?.serviceEnded ?? false;
    this.serviceStarted = data?.serviceStarted;
    this.vatRate = data?.tax?.taxRate ?? 0;
    this.type = [
      BilledType.Quarterly,
      BilledType.Monthly,
      BilledType.Annually,
    ].includes(data?.billingCycle as BilledType)
      ? 'RECURRING'
      : data?.billingCycle;
  }

  static toPayload(data: ServiceLineItem): BulkUpdateServiceLineItem {
    return {
      serviceLineItemId: data?.metadata?.id ?? '',
      name: data?.description ?? '',
      quantity: data?.quantity ?? 1,
      price: data?.price ?? 0,
      billed: data?.billingCycle ?? BilledType.Monthly,
      isDeleted: data?.serviceEnded ?? false,
      serviceStarted: data?.serviceStarted
        ? new UTCDate(data.serviceStarted)
        : null,
      vatRate: data?.tax?.taxRate ?? 0,
      type: [
        BilledType.Quarterly,
        BilledType.Monthly,
        BilledType.Annually,
      ].includes(data?.billingCycle)
        ? 'RECURRING'
        : data?.billingCycle,
    };
  }
}

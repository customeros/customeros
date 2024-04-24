import { makeAutoObservable } from 'mobx';

import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import {
  BilledType,
  ServiceLineItem,
  ServiceLineItemBulkUpdateItem,
  InvoiceSimulateServiceLineInput,
} from '@graphql/types';

interface IServiceLineItem extends ServiceLineItem {
  isNew: boolean;
  isDeleted: boolean;
  newVersion: boolean;
  closedVersion: boolean;
  frontendMetadata?: null | {
    color: string;
    shapeVariant: string | number;
  };
}

class ServiceLineItemStore {
  serviceLineItem: IServiceLineItem | null = null;
  public revisedFields: string[] = [];
  public lastRevisedFields: string[] = [];
  public isNewlyAdded: boolean = false;
  public lastRevisionIsNewlyAdded: boolean = false;

  constructor() {
    makeAutoObservable(this);
  }

  get billingValue() {
    return billedTypeOptions.find(
      (option) => option.value === this.serviceLineItem?.billingCycle,
    );
  }

  resetStateFields() {
    this.revisedFields = [];
    this.lastRevisedFields = [];
    this.isNewlyAdded = false;
    this.lastRevisionIsNewlyAdded = false;
  }
  setServiceLineItem(item: IServiceLineItem) {
    this.serviceLineItem = { ...item };
    if (!item.parentId) {
      this.isNewlyAdded = true;
    }
  }

  updateBilledType(billingCycle: BilledType) {
    if (
      this.serviceLineItem &&
      this.serviceLineItem.billingCycle !== billingCycle
    ) {
      this.revisedFields = Array.from(
        new Set([...this.revisedFields, 'billingCycle']),
      );
      this.serviceLineItem.billingCycle = billingCycle;
    }
  }

  updateQuantity(quantity: string) {
    if (this.serviceLineItem && this.serviceLineItem.quantity !== quantity) {
      this.revisedFields = Array.from(
        new Set([...this.revisedFields, 'quantity']),
      );
      this.serviceLineItem.quantity = parseFloat(quantity);
    }
  }

  updatePrice(price: string) {
    if (
      this.serviceLineItem &&
      this.serviceLineItem.price !== parseFloat(price)
    ) {
      this.revisedFields = Array.from(
        new Set([...this.revisedFields, 'price']),
      );
      this.serviceLineItem.price = parseFloat(price);
    }
  }
  updateDescription(desc: string) {
    if (this.serviceLineItem && this.serviceLineItem.description !== desc) {
      this.revisedFields = Array.from(
        new Set([...this.revisedFields, 'description']),
      );
      this.serviceLineItem.description = desc;
    }
  }

  updateTaxRate(taxRate: number) {
    if (
      this.serviceLineItem &&
      this.serviceLineItem.tax &&
      this.serviceLineItem.tax.taxRate !== taxRate
    ) {
      this.revisedFields = Array.from(
        new Set([...this.revisedFields, 'taxRate']),
      );
      this.serviceLineItem.tax.taxRate = taxRate;
    }
  }

  updateStartDate(date: Date | null) {
    if (this.serviceLineItem && this.serviceLineItem.serviceStarted !== date) {
      this.revisedFields = Array.from(
        new Set([...this.revisedFields, 'serviceStarted']),
      );
      this.serviceLineItem.serviceStarted = date;
    }
  }

  addComment(comment: string) {
    if (this.serviceLineItem) {
      this.serviceLineItem.comments += ` ${comment}`;
    }
  }
  setIsDeleted(isDeleted: boolean) {
    if (this.serviceLineItem && this.serviceLineItem.isDeleted !== isDeleted) {
      this.serviceLineItem.isDeleted = isDeleted;
    }
  }

  setIsClosedVersion(isDeleted: boolean) {
    if (this.serviceLineItem) {
      this.serviceLineItem.closedVersion = isDeleted;
    }
  }
  setIsEnded() {
    if (this.serviceLineItem) {
      this.serviceLineItem.serviceEnded = new Date();
    }
  }

  shouldReactToRevisedFields(): boolean {
    if (!this.revisedFields.length) {
      return false;
    }
    if (this.revisedFields.length !== this.lastRevisedFields.length) {
      this.lastRevisedFields = this.revisedFields;

      return true;
    }
    if (
      this.isNewlyAdded &&
      this.isNewlyAdded !== this.lastRevisionIsNewlyAdded
    ) {
      this.lastRevisionIsNewlyAdded = this.isNewlyAdded;

      return true;
    }

    return false;
  }

  get serviceLineItemValues() {
    return {
      billingCycle: this.serviceLineItem?.billingCycle,
      quantity: this.serviceLineItem?.quantity,
      price: this.serviceLineItem?.price,
      description: this.serviceLineItem?.description,
      taxRate: this.serviceLineItem?.tax.taxRate,
      serviceStarted: this.serviceLineItem?.serviceStarted,
      comments: this.serviceLineItem?.comments,
      serviceEnded: this.serviceLineItem?.serviceEnded,
      closedVersion: this.serviceLineItem?.closedVersion,
      newVersion: this.serviceLineItem?.newVersion,
    };
  }

  getServiceLineItemBulkUpdateItem(): ServiceLineItemBulkUpdateItem | null {
    // Do not save if no fields were revised nor if it is a newly added service line item
    if (!this.revisedFields.length && !this.serviceLineItem?.isNew) {
      return null;
    }

    if (!this.serviceLineItem?.isNew) {
      return {
        serviceLineItemId: this.serviceLineItem?.parentId || '',
        name: this.serviceLineItem?.description,
        billed: this.serviceLineItem?.billingCycle,
        price: this.serviceLineItem?.price,
        quantity: this.serviceLineItem?.quantity,
        vatRate: this.serviceLineItem?.tax.taxRate,
        comments: this.serviceLineItem?.comments,
        serviceStarted: this.serviceLineItem?.serviceStarted,
        closeVersion: this.serviceLineItem?.closedVersion,
      };
    }

    return {
      name: this.serviceLineItem?.description,
      billed: this.serviceLineItem?.billingCycle,
      price: this.serviceLineItem?.price,
      quantity: this.serviceLineItem?.quantity,
      vatRate: this.serviceLineItem?.tax.taxRate,
      comments: this.serviceLineItem?.comments,
      serviceStarted: this.serviceLineItem?.serviceStarted,
      closeVersion: this.serviceLineItem?.closedVersion,
    };
  }
  getInvoiceSimulationServiceLineItem(): InvoiceSimulateServiceLineInput | null {
    if (this.serviceLineItem === null) {
      throw new Error('Service line item is not set.');
    }

    return {
      key: this.serviceLineItem.metadata.id,
      parentId: this.serviceLineItem.parentId,
      description: this.serviceLineItem.description ?? 'Unnamed',
      billingCycle: this.serviceLineItem.billingCycle,
      price: this.serviceLineItem.price,
      quantity: parseFloat(this.serviceLineItem.quantity),
      taxRate: this.serviceLineItem.tax.taxRate,
      serviceStarted: this.serviceLineItem.serviceStarted,

      serviceLineItemId: this.serviceLineItem.isNew
        ? undefined
        : this.serviceLineItem.metadata.id,
    };
  }
}

export default ServiceLineItemStore;

import { makeAutoObservable } from 'mobx';

import { billedTypeOptions } from '@organization/components/Tabs/panels/AccountPanel/utils';
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
  isModification: boolean;
  nextBilling?: Date | null;
  frontendMetadata?: null | {
    color: string;
    shapeVariant: string | number;
  };
}
// todo move to context
class ServiceLineItemStore {
  serviceLineItem: IServiceLineItem | null = null;
  private revisedFields: Map<string, number> = new Map();
  private lastRevision: Map<string, number> = new Map();
  public isNewlyAdded: boolean = false;
  public lastRevisionIsNewlyAdded: boolean = false;
  public revisedFieldsSummary: string[] = [];

  constructor() {
    makeAutoObservable(this);
  }

  get billingValue() {
    return billedTypeOptions.find(
      (option) => option.value === this.serviceLineItem?.billingCycle,
    );
  }

  resetStateFields() {
    this.revisedFields = new Map();
    this.lastRevision = new Map();
    this.isNewlyAdded = false;
    this.lastRevisionIsNewlyAdded = false;
  }
  // Utilize this method to record changes to fields
  private markFieldAsRevised(field: string) {
    const currentCount = this.revisedFields.get(field) || 0;
    this.revisedFields.set(field, currentCount + 1);
  }

  setServiceLineItem(item: IServiceLineItem & { parentId?: string | null }) {
    this.serviceLineItem = item;
    if (!item.isModification && item.isNew) {
      this.isNewlyAdded = true;
    }
  }

  updateBilledType(billingCycle: BilledType) {
    if (
      this.serviceLineItem &&
      this.serviceLineItem.billingCycle !== billingCycle
    ) {
      this.markFieldAsRevised('billingCycle');

      this.serviceLineItem.billingCycle = billingCycle;
    }
  }

  updateQuantity(quantity: string) {
    if (this.serviceLineItem && this.serviceLineItem.quantity !== quantity) {
      this.markFieldAsRevised('quantity');

      this.serviceLineItem.quantity = parseFloat(quantity);
    }
  }

  updatePrice(price: string) {
    if (
      this.serviceLineItem &&
      this.serviceLineItem.price !== parseFloat(price)
    ) {
      this.markFieldAsRevised('price');

      this.serviceLineItem.price = parseFloat(price);
    }
  }
  updateDescription(desc: string) {
    if (this.serviceLineItem && this.serviceLineItem.description !== desc) {
      this.markFieldAsRevised('description');

      this.serviceLineItem.description = desc;
    }
  }

  updateTaxRate(taxRate: number) {
    if (
      this.serviceLineItem &&
      this.serviceLineItem.tax &&
      this.serviceLineItem.tax.taxRate !== taxRate
    ) {
      this.markFieldAsRevised('taxRate');
      this.serviceLineItem.tax.taxRate = taxRate;
    }
  }

  updateStartDate(date: Date | null) {
    if (this.serviceLineItem && this.serviceLineItem.serviceStarted !== date) {
      this.markFieldAsRevised('serviceStarted');

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
      this.markFieldAsRevised('isDeleted');
    }
  }

  setIsClosedVersion(isDeleted: boolean) {
    if (this.serviceLineItem) {
      this.serviceLineItem.closedVersion = isDeleted;
      this.markFieldAsRevised('closedVersion');
    }
  }

  private updateRevisedFieldsSummary() {
    this.revisedFieldsSummary = Array.from(this.revisedFields.entries()).map(
      ([field, count]) => `${field} (changed ${count} times)`,
    );
  }
  shouldReactToRevisedFields(): boolean {
    let shouldReact = false;
    this.revisedFields.forEach((count, field) => {
      const lastCount = this.lastRevision.get(field) || 0;
      if (count !== lastCount) {
        shouldReact = true;
        this.lastRevision.set(field, count);
      }
    });

    if (!shouldReact && this.revisedFields.size !== this.lastRevision.size) {
      shouldReact = true;
    }

    if (shouldReact) {
      this.updateRevisedFieldsSummary();
    }

    return shouldReact;
  }

  isFieldRevised(fieldName: string): boolean {
    return this.revisedFields.has(fieldName);
  }

  get uiMetadata(): {
    color: string;
    shapeVariant: string | number;
  } {
    return {
      color: this.serviceLineItem?.frontendMetadata?.color || '',
      shapeVariant: this.serviceLineItem?.frontendMetadata?.shapeVariant || '',
    };
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
    const hasRevisedFields = this.revisedFields.size > 0;

    if (
      (!hasRevisedFields && !this.serviceLineItem?.isNew) ||
      !this.serviceLineItem
    ) {
      return null;
    }

    return {
      serviceLineItemId: this.serviceLineItem.parentId.length
        ? this.serviceLineItem.parentId
        : undefined,
      name: this.serviceLineItem.description,
      billed: this.serviceLineItem.billingCycle,
      price: this.serviceLineItem.price,
      quantity: this.serviceLineItem.quantity,
      vatRate: this.serviceLineItem.tax?.taxRate,
      comments: this.serviceLineItem.comments,
      serviceStarted: this.serviceLineItem.serviceStarted,
      closeVersion:
        this.serviceLineItem.closedVersion || this.serviceLineItem.isDeleted,
      newVersion:
        (this.serviceLineItem.isModification &&
          !this.revisedFields.has('description')) ||
        false,
    };
  }
  getInvoiceSimulationServiceLineItem(): InvoiceSimulateServiceLineInput | null {
    if (this.serviceLineItem === null) {
      throw new Error('Service line item is not set.');
    }

    return {
      key: this.serviceLineItem.metadata.id,
      parentId:
        this.serviceLineItem.isModification || !this.isNewlyAdded
          ? this.serviceLineItem.parentId
          : undefined,
      description: this.serviceLineItem.description ?? 'Unnamed',
      billingCycle: this.serviceLineItem.billingCycle,
      price: this.serviceLineItem.price,
      quantity: parseFloat(this.serviceLineItem.quantity),
      taxRate: this.serviceLineItem.tax.taxRate,
      serviceStarted: this.serviceLineItem.serviceStarted,
      closeVersion:
        this.serviceLineItem.closedVersion || this.serviceLineItem.isDeleted,
      serviceLineItemId:
        !this.serviceLineItem.isModification && !this.isNewlyAdded
          ? this.serviceLineItem.metadata.id
          : undefined,
    };
  }
}

export default ServiceLineItemStore;

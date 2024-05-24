import { makeAutoObservable } from 'mobx';
import { GraphQLClient } from 'graphql-request';

import { uuidv4 } from '@spaces/utils/generateUuid';
import { DateTimeUtils } from '@spaces/utils/date.ts';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  BilledType,
  DataSource,
  InvoiceSimulate,
  ServiceLineItem,
} from '@graphql/types';
import InvoiceListStore from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/InvoicePreviewList.store.ts';
import { getVersionFromUUID } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/highlighters';

import ServiceLineItemStore from './Service.store';

const defaultValue = {
  billingCycle: BilledType.Monthly,
  closed: false,
  comments: '',
  createdBy: '',
  description: 'Unnamed',
  externalLinks: [],
  metadata: {
    id: 'default-meta-id',
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    sourceOfTruth: DataSource.Na,
  },
  parentId: 'default-parent-id',
  price: 0,
  quantity: 0,
  serviceEnded: null,
  serviceStarted: new Date(),
  isNew: false,
  isDeleted: false,
  tax: { salesTax: false, vat: false, taxRate: 0 },
};
interface IServiceLineItem extends ServiceLineItem {
  isNew: boolean;
  isDeleted: boolean;
  frontendMetadata?: null | {
    color: string;
    shapeVariant: string;
  };
}
function groupServicesByParentId(
  services: ServiceLineItemStore[],
): Array<ServiceLineItemStore[]> {
  const grouped: Record<string, ServiceLineItemStore[]> = {};

  services.forEach((service) => {
    const parentId = service.serviceLineItem?.parentId;
    if (parentId) {
      if (!grouped[parentId]) {
        grouped[parentId] = [];
      }
      grouped[parentId].push(service);
    }
  });

  const sortedGroups = Object.values(grouped).map((group) =>
    group.sort(
      (a, b) =>
        new Date(a.serviceLineItem?.serviceStarted).getTime() -
        new Date(b.serviceLineItem?.serviceStarted).getTime(),
    ),
  );

  // Filtering groups to exclude those where all items have 'serviceEnded' as null
  const filtered = sortedGroups.filter((group) =>
    group.some((service) => service.serviceLineItem?.serviceEnded === null),
  );

  return filtered;
}
// todo move to context
export class ServiceFormStore {
  oneTimeServices: ServiceLineItemStore[][] = [];
  subscriptionServices: ServiceLineItemStore[][] = [];
  public isSimulationRunning: boolean = false;
  public isSaving: boolean = false;
  public keyColorPairs: Record<string, string> = {};
  private contractId: string = '';
  private readonly graphQLClient: GraphQLClient = getGraphQLClient();

  constructor() {
    makeAutoObservable(this);
  }

  set contractIdValue(contractId: string) {
    this.contractId = contractId;
  }

  async runSimulation(invoiceListStore: InvoiceListStore): Promise<void> {
    const payload = this.getInvoiceSimulationInput();
    if (!payload.length) {
      return;
    }
    this.isSimulationRunning = true;
    try {
      const response = await this.graphQLClient.request<TSimulateQueryResult>(
        SimulateInvoiceDocument,
        {
          input: {
            contractId: this.contractId,
            serviceLines: payload,
          },
        },
      );
      invoiceListStore.initializeSimulatedInvoices(
        response?.invoice_Simulate.map((e) => ({
          ...e,
          invoiceLineItems: e.invoiceLineItems.map((lineItem) => ({
            serviceLineItemStore: this.getServiceLineItemById(lineItem.key),

            ...lineItem,
          })),
        })),
      );
    } catch (error) {
      console.error(`Simulation failed: ${error}`);
    } finally {
      this.isSimulationRunning = false;
    }
  }

  getServiceLineItemById(
    serviceLineItemId: string,
  ): ServiceLineItemStore | null {
    const allServices = [
      ...this.oneTimeServices,
      ...this.subscriptionServices,
    ].flat();

    return (
      allServices.find(
        (service) => service.serviceLineItem?.metadata.id === serviceLineItemId,
      ) || null
    );
  }

  shouldReact(): boolean {
    let hasChanges = false;
    const checkAndReact = (services: ServiceLineItemStore[]) => {
      services.forEach((service) => {
        if (service.shouldReactToRevisedFields()) {
          hasChanges = true;
        }
      });
    };

    this.oneTimeServices.forEach(checkAndReact);
    this.subscriptionServices.forEach(checkAndReact);

    return hasChanges;
  }

  clearUsedColors() {
    this.subscriptionServices = [];
    this.oneTimeServices = [];
  }
  initializeServices(contractLineItems?: ServiceLineItem[]) {
    if (contractLineItems?.length) {
      const { subscription, once } = contractLineItems.reduce<{
        once: ServiceLineItemStore[];
        subscription: ServiceLineItemStore[];
      }>(
        (acc, item) => {
          const key: 'subscription' | 'once' = [
            BilledType.Monthly,
            BilledType.Quarterly,
            BilledType.Annually,
          ].includes(item.billingCycle)
            ? 'subscription'
            : 'once';
          const newItemStore = new ServiceLineItemStore();
          newItemStore.setServiceLineItem({
            ...item,
            isNew: false,
            isModification: false,
            closedVersion: item.serviceEnded !== null,
            newVersion: false,
            isDeleted: false,
            frontendMetadata: null,
          });
          acc[key].push(newItemStore);

          return acc;
        },
        { subscription: [], once: [] },
      );

      this.oneTimeServices = groupServicesByParentId(once);
      this.subscriptionServices = groupServicesByParentId(subscription);
    } else {
      this.oneTimeServices = [];
      this.subscriptionServices = [];
    }
  }

  private createServiceLineItem(
    id: string | null,
    modification: Partial<IServiceLineItem> | ServiceLineItem,
    isSubscription: boolean,
  ) {
    const newItemStore = new ServiceLineItemStore();
    const newItemId = uuidv4();
    // const backgroundColor = getColorByUUID(newItemId, this.usedColors);
    const highlightVersion = getVersionFromUUID(newItemId);
    const getNextBillingDate = (): Date | undefined => {
      const today = new Date().toString();
      if (!modification.serviceStarted)
        return DateTimeUtils.addDays(new Date().toString(), 1);
      let nextBillingDate = new Date(modification.serviceStarted);
      switch (modification.billingCycle) {
        case BilledType.Monthly: {
          const diffMonths = DateTimeUtils.differenceInMonths(
            modification.serviceStarted,
            today,
          );

          nextBillingDate = DateTimeUtils.addMonth(
            modification.serviceStarted,
            diffMonths + 1,
          );
          break;
        }
        case BilledType.Annually: {
          const diffYears = DateTimeUtils.differenceInYears(
            modification.serviceStarted,
            today,
          );
          nextBillingDate = DateTimeUtils.addYears(
            modification.serviceStarted,
            diffYears + 1,
          );
          break;
        }
        case BilledType.Quarterly: {
          const diffMonths = DateTimeUtils.differenceInMonths(
            modification.serviceStarted,
            today,
          );
          nextBillingDate = DateTimeUtils.addMonth(
            modification.serviceStarted,
            diffMonths + 3,
          );
          break;
        }
        default:
          return undefined; // or throw an error if an unknown billing cycle is not allowed
      }

      return nextBillingDate;
    };

    // this.usedColors.push(backgroundColor);
    const newServiceData = {
      ...defaultValue,
      ...modification,

      serviceEnded: null,
      closedVersion: false,
      newVersion: !!id,
      metadata: { ...defaultValue.metadata, id: newItemId },
      parentId: id ?? '',
      isModification: !!id,
      serviceStarted: id
        ? DateTimeUtils.addDays(modification.serviceStarted, 1)
        : DateTimeUtils.addDays(new Date().toString(), 1),
      nextBilling:
        id && modification.billingCycle !== BilledType.Once
          ? getNextBillingDate()
          : null,
      isNew: true,
      frontendMetadata: {
        color: 'transparent',
        shapeVariant: highlightVersion,
      },
      createdBy: null,
    };
    newItemStore.setServiceLineItem(newServiceData);
    const targetArray = isSubscription
      ? this.subscriptionServices
      : this.oneTimeServices;
    const serviceGroup = targetArray.find(
      (group) => group[0]?.serviceLineItem?.parentId === id,
    );

    if (serviceGroup) {
      serviceGroup.push(newItemStore);
    } else {
      targetArray.push([newItemStore]);
    }
  }

  addService(serviceId: string | null, isSubscription?: boolean) {
    const serviceArray = [
      ...this.subscriptionServices,
      ...this.oneTimeServices,
    ];

    const prevValue = serviceArray
      .flat()
      .filter(
        (e) =>
          e.serviceLineItem?.parentId === serviceId ||
          e.serviceLineItem?.metadata.id === serviceId,
      )
      .reduce((maxDate: null | ServiceLineItemStore, currentObj) => {
        if (!maxDate) return currentObj;

        return DateTimeUtils.isBefore(
          maxDate?.serviceLineItem?.serviceStarted || 0,
          currentObj?.serviceLineItem?.serviceStarted || 0,
        )
          ? currentObj
          : maxDate;
      }, null);

    this.createServiceLineItem(
      serviceId,
      prevValue?.serviceLineItemValues ?? {
        billingCycle: isSubscription ? BilledType.Monthly : BilledType.Once,
      },
      !!isSubscription,
    );
  }

  getInvoiceSimulationInput() {
    const allServiceStores = [
      ...this.oneTimeServices.flat(),
      ...this.subscriptionServices.flat(),
    ];

    return allServiceStores
      .map((store) => {
        return store.getInvoiceSimulationServiceLineItem();
      })
      .filter((e) => e);
  }

  getServiceLineItemsBulkUpdateInput() {
    const allServiceStores = [
      ...this.oneTimeServices.flat(),
      ...this.subscriptionServices.flat(),
    ];

    return allServiceStores
      .map((store) => {
        return store.getServiceLineItemBulkUpdateItem();
      })
      .filter((e) => e);
  }

  async saveServiceLineItems() {
    if (!this.subscriptionServices.length && !this.oneTimeServices.length) {
      return;
    }
    const payload = this.getServiceLineItemsBulkUpdateInput();
    if (!payload.length) {
      return;
    }

    this.isSaving = true;
    try {
      const response =
        await this.graphQLClient?.request<TBulkUpdateServicesResult>(
          UpdateServicesDocument,
          {
            input: {
              contractId: this.contractId,
              invoiceNote: '',
              serviceLineItems: payload,
            },
          },
        );

      return response;
    } catch (error) {
      throw new Error(`Failed to save service line items: ${error}`);
    } finally {
      this.isSaving = false;
    }
  }
}
type TSimulateQueryResult = { invoice_Simulate: Array<InvoiceSimulate> };

const SimulateInvoiceDocument = `
mutation simulateInvoice($input: InvoiceSimulateInput!) {
  invoice_Simulate(input: $input) {
    amount
    currency
    due
    invoiceNumber
    invoicePeriodEnd
    invoicePeriodStart
    issued
    note
    offCycle
    postpaid
    subtotal
    taxDue
    total
    invoiceLineItems {
      key
      description
      price
      quantity
      subtotal
      taxDue
      total
    }
    customer {
      name
      email
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
      addressRegion
    }
    provider {
      name
      addressLine1
      addressLine2
      addressZip
      addressLocality
      addressCountry
      addressRegion
    }
  }
}
    `;

type TBulkUpdateServicesResult = { serviceLineItem_BulkUpdate: Array<string> };

export const UpdateServicesDocument = `
    mutation updateServices($input: ServiceLineItemBulkUpdateInput!) {
  serviceLineItem_BulkUpdate(input: $input) 
}
    `;

import { makeAutoObservable } from 'mobx';
import { GraphQLClient } from 'graphql-request';

import { uuidv4 } from '@spaces/utils/generateUuid';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  BilledType,
  DataSource,
  InvoiceSimulate,
  ServiceLineItem,
} from '@graphql/types';
import InvoiceListStore from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/stores/InvoicePreviewList.store';
import { HighlightColor } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/highlighters/utils';
import {
  getColorByUUID,
  getVersionFromUUID,
} from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/highlighters';

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
    group.every((service) => service.serviceLineItem?.serviceEnded === null),
  );

  return filtered;
}

export class ServiceFormStore {
  oneTimeServices: ServiceLineItemStore[][] = [];
  subscriptionServices: ServiceLineItemStore[][] = [];
  public isSimulationRunning: boolean = false;
  public isSaving: boolean = false;
  public keyColorPairs: Record<string, string> = {};
  private usedColors: Array<HighlightColor> = [];
  private contractId: string = '';
  private readonly graphQLClient: GraphQLClient = getGraphQLClient();
  private lastOneTimeSnapshot: string = '';
  private lastSubscriptionSnapshot: string = '';

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

  private serializeServices(services: ServiceLineItemStore[][]): string {
    return JSON.stringify(
      services.flat().map((sli) => ({
        ...sli.serviceLineItem,
      })),
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
    this.usedColors = [];
    this.subscriptionServices = [];
    this.oneTimeServices = [];
  }
  initializeServices(contractLineItems?: ServiceLineItem[]) {
    this.usedColors = [];
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
    modification: Partial<IServiceLineItem>,
    isSubscription: boolean,
  ) {
    const newItemStore = new ServiceLineItemStore();
    const newItemId = uuidv4();
    const backgroundColor = getColorByUUID(newItemId, this.usedColors);
    const highlightVersion = getVersionFromUUID(newItemId);
    this.keyColorPairs = {
      ...this.keyColorPairs,
      [newItemId]: backgroundColor,
    };

    this.usedColors.push(backgroundColor);
    newItemStore.setServiceLineItem({
      closedVersion: false,
      newVersion: false,
      ...defaultValue,
      ...modification,
      metadata: { ...defaultValue.metadata, id: newItemId },
      parentId: id ?? newItemId,
      isModification: !!id,
      isNew: true,
      frontendMetadata: {
        color: backgroundColor,
        shapeVariant: highlightVersion,
      },
      createdBy: null,
    });

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
    const serviceArray = isSubscription
      ? this.subscriptionServices
      : this.oneTimeServices;
    const prevValue = serviceArray
      .flat()
      .find(
        (e) => e.serviceLineItem?.metadata.id === serviceId,
      )?.serviceLineItemValues;

    this.createServiceLineItem(
      serviceId,
      prevValue ?? {
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

    this.isSaving = true;
    try {
      const response =
        await this.graphQLClient?.request<TBulkUpdateServicesResult>(
          UpdateServicesDocument,
          {
            input: {
              contractId: this.contractId,
              invoiceNote: '',
              serviceLineItems: this.getServiceLineItemsBulkUpdateInput(),
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

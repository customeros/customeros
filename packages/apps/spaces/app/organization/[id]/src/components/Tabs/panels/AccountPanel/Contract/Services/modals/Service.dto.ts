import { SelectOption } from '@shared/types/SelectOptions';
import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import {
  BilledType,
  ServiceLineItem,
  ContractRenewalCycle,
  ServiceLineItemInput,
  ServiceLineItemUpdateInput,
} from '@graphql/types';

export interface ServiceForm {
  name?: string;
  price?: number;
  quantity?: number;

  appSource?: string;
  serviceStartedAt?: Date;
  externalReference?: string;
  renewalCycle?: ContractRenewalCycle;
  billed?: SelectOption<BilledType> | null;
}

type ServiceItem = Omit<ServiceLineItem, 'id'>;

export class ServiceDTO implements ServiceForm {
  name?: string;
  price?: number;
  quantity?: number;
  appSource?: string;
  billed?: SelectOption<BilledType> | null;
  serviceStartedAt?: Date;
  externalReference?: string;
  renewalCycle?: ContractRenewalCycle;

  constructor(data?: ServiceItem) {
    this.quantity = data?.quantity;
    this.name = data?.name;
    this.price = data?.price;
    this.appSource = data?.appSource;
    this.billed =
      billedTypeOptions.find((o) => o.value === data?.billed) ?? null;
  }

  static toForm(data?: ServiceItem): ServiceForm {
    return new ServiceDTO(data);
  }

  static toPayload(
    data: ServiceForm,
    contractId: string,
  ): ServiceLineItemInput {
    return {
      contractId,
      quantity: data?.quantity,
      name: data?.name,
      price: data?.price,
      appSource: data?.appSource,
      billed: data?.billed?.value,
    };
  }
  static toUpdatePayload(
    data: ServiceForm,
  ): Omit<ServiceLineItemUpdateInput, 'serviceLineItemId'> {
    return {
      quantity: data?.quantity,
      name: data?.name,
      price: data?.price,
      appSource: data?.appSource,
      billed: data?.billed?.value,
    };
  }
}

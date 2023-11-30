import { SelectOption } from '@shared/types/SelectOptions';
import {
  RenewalCycle,
  BillingDetails,
  BillingDetailsInput,
} from '@graphql/types';

import { frequencyOptions } from '../utils';

export interface BillingDetailsForm {
  amount: string | null;
  frequency: SelectOption<RenewalCycle> | null;
}

export class BillingDetailsDTO implements BillingDetailsForm {
  amount: string | null;
  frequency: SelectOption<RenewalCycle> | null;

  constructor(data?: BillingDetails | null) {
    this.amount = data?.amount ? data.amount.toString() : '';
    this.frequency =
      frequencyOptions.find(({ value }) => value === data?.frequency) ?? null;
  }

  static toForm(data?: BillingDetails | null) {
    return new BillingDetailsDTO(data);
  }

  static toPayload(
    data: BillingDetailsForm & { organizationId: string } & Pick<
        BillingDetails,
        'renewalCycle' | 'renewalCycleStart'
      >,
  ): BillingDetailsInput {
    return {
      id: data.organizationId,
      amount: data.amount ? parseFloat(data.amount) : undefined,
      frequency: data.frequency?.value ?? undefined,
      renewalCycle: data.renewalCycle,
      renewalCycleStart: data.renewalCycleStart,
    };
  }
}

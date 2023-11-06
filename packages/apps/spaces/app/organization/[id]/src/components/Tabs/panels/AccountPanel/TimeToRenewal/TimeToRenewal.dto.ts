import { SelectOption } from '@shared/types/SelectOptions';
import {
  RenewalCycle,
  BillingDetails,
  BillingDetailsInput,
} from '@graphql/types';

import { frequencyOptions } from '../utils';

interface TimeToRenewalForm {
  renewalCycleStart: Date | null;
  renewalCycle: SelectOption<RenewalCycle> | null;
}

export class TimeToRenewalDTO implements TimeToRenewalForm {
  renewalCycle: SelectOption<RenewalCycle> | null;
  renewalCycleStart: Date | null;

  constructor(data?: BillingDetails | null) {
    this.renewalCycle =
      frequencyOptions.find((o) => o.value === data?.renewalCycle) ?? null;
    this.renewalCycleStart = data?.renewalCycleStart
      ? new Date(data.renewalCycleStart)
      : null;
  }

  static toForm(data?: BillingDetails | null): TimeToRenewalForm {
    return new TimeToRenewalDTO(data);
  }

  static toPayload(
    data: TimeToRenewalForm & { id: string } & Pick<
        BillingDetails,
        'amount' | 'frequency'
      >,
  ): BillingDetailsInput {
    return {
      id: data.id,
      amount: data.amount,
      frequency: data.frequency,
      renewalCycle: data.renewalCycle?.value,
      renewalCycleStart: data.renewalCycleStart?.toISOString(),
    };
  }
}

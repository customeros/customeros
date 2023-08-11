import { SelectOption } from '@shared/types/SelectOptions';

export enum FrequencyOptions {
  WEEKLY = 'WEEKLY',
  BIWEEKLY = 'BIWEEKLY',
  MONTHLY = 'MONTHLY',
  QUARTERLY = 'QUARTERLY',
  BIANNUALLY = 'BIANNUALLY',
  ANNUALLY = 'ANNUALLY',
}

export const frequencyOptions: SelectOption[] = [
  { label: 'Weekly', value: FrequencyOptions.WEEKLY },
  { label: 'Biweekly', value: FrequencyOptions.BIWEEKLY },
  { label: 'Monthly', value: FrequencyOptions.MONTHLY },
  { label: 'Quarterly', value: FrequencyOptions.QUARTERLY },
  { label: 'Biannually', value: FrequencyOptions.BIANNUALLY },
  { label: 'Annually', value: FrequencyOptions.ANNUALLY },
];

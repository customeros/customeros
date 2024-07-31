import { DateTimeUtils } from '@utils/date.ts';
import { SelectOption } from '@shared/types/SelectOptions';
import {
  Maybe,
  Contract,
  BilledType,
  ServiceLineItem,
  ContractBillingCycle,
  ContractRenewalCycle,
  OpportunityRenewalLikelihood,
} from '@graphql/types';

export function getRenewalLikelihoodColor(
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood> | undefined,
) {
  switch (renewalLikelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'greenLight';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'yellow';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'orangeDark';
    case OpportunityRenewalLikelihood.ZeroRenewal:
      return 'gray';
    default:
      return 'gray';
  }
}

export function getRenewalLikelihoodLabel(
  renewalLikelihood?: Maybe<OpportunityRenewalLikelihood> | undefined,
) {
  switch (renewalLikelihood) {
    case OpportunityRenewalLikelihood.HighRenewal:
      return 'High';
    case OpportunityRenewalLikelihood.MediumRenewal:
      return 'Medium';
    case OpportunityRenewalLikelihood.LowRenewal:
      return 'Low';
    case OpportunityRenewalLikelihood.ZeroRenewal:
      return 'Zero';
    default:
      return '';
  }
}
export const contractRenewalCycle: SelectOption<
  ContractRenewalCycle | 'MULTI_YEAR'
>[] = [
  { label: 'Monthly', value: ContractRenewalCycle.MonthlyRenewal },
  { label: 'Quarterly', value: ContractRenewalCycle.QuarterlyRenewal },
  { label: 'Annual', value: ContractRenewalCycle.AnnualRenewal },
  { label: 'Multi-year', value: 'MULTI_YEAR' },
];
export const billingFrequencyOptions: SelectOption<ContractBillingCycle>[] = [
  { label: 'Monthly', value: ContractBillingCycle.MonthlyBilling },
  { label: 'Quarterly', value: ContractBillingCycle.QuarterlyBilling },
  { label: 'Annual', value: ContractBillingCycle.AnnualBilling },
];
export const contractBillingCycleOptions: SelectOption<number>[] = [
  { label: 'monthly', value: 1 },
  { label: 'quarterly', value: 3 },
  { label: 'annually', value: 12 },
];

export const paymentDueOptions: SelectOption<number>[] = [
  { label: '0 days', value: 0 },
  { label: '15 days', value: 15 },
  { label: '30 days', value: 30 },
  { label: '45 days', value: 45 },
  { label: '60 days', value: 60 },
  { label: '90 days', value: 90 },
];

export function calculateMaxArr(
  serviceLineItems: ServiceLineItem[],
  contract: Contract,
) {
  if (!serviceLineItems || !contract) return 0;

  const totalAnnualPrice = serviceLineItems?.reduce((acc, sli) => {
    if (
      sli.closed ||
      (sli.serviceEnded && DateTimeUtils.isPast(sli.serviceEnded))
    ) {
      return acc;
    }
    const annualPrice = calculateAnnualPrice(sli);

    return acc + annualPrice;
  }, 0);
  const proratedArr = contract.contractEnded
    ? prorateArr(
        totalAnnualPrice,
        monthsUntilContractEnd(new Date(), contract.contractEnded),
      )
    : totalAnnualPrice;

  return roundHalfUpFloat(proratedArr, 2);
}

function calculateAnnualPrice(sli: ServiceLineItem) {
  switch (sli.billingCycle) {
    case BilledType.Annually:
      return sli.price * sli.quantity;
    case BilledType.Monthly:
      return sli.price * sli.quantity * 12;
    case BilledType.Quarterly:
      return sli.price * sli.quantity * 3;
    default:
      return 0;
  }
}

function prorateArr(arr: number, monthsUntilEnd: number) {
  const monthsInYear = 12;

  return (arr / monthsInYear) * monthsUntilEnd;
}

function monthsUntilContractEnd(currentDate: Date, endDate: Date) {
  const months = (endDate.getFullYear() - currentDate.getFullYear()) * 12;

  return months - currentDate.getMonth() + endDate.getMonth();
}

function roundHalfUpFloat(num: number, decimalPlaces: number) {
  const factor = Math.pow(10, decimalPlaces);

  return Math.round(num * factor) / factor;
}

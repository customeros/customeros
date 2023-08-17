import {
  addWeeks,
  addMonths,
  addYears,
  differenceInWeeks,
  differenceInMonths,
  differenceInYears,
  formatDistanceToNowStrict,
} from 'date-fns';

export function formatSocialUrl(value = '') {
  let url = value;

  if (url.startsWith('http')) {
    url = url.replace('https://', '');
  }
  if (url.startsWith('www')) {
    url = url.replace('www.', '');
  }
  if (url.includes('twitter')) {
    url = url.replace('twitter.com', '');
  }
  if (url.includes('linkedin.com/in')) {
    url = url.replace('linkedin.com/in', '');
  }
  if (url.includes('linkedin.com/company')) {
    url = url.replace('linkedin.com/company', '');
  }

  return url;
}

export type RenewalFrequency =
  | 'WEEKLY'
  | 'BIWEEKLY'
  | 'MONTHLY'
  | 'QUARTERLY'
  | 'BIANNUALLY'
  | 'ANNUALLY';

export function getTimeToNextRenewal(
  renewalStart: Date,
  renewalFrequency: RenewalFrequency,
): string[] {
  const now = new Date();
  let totalRenewals: number;
  let nextRenewalDate: Date;

  switch (renewalFrequency) {
    case 'WEEKLY':
      totalRenewals = differenceInWeeks(now, renewalStart);
      nextRenewalDate = addWeeks(renewalStart, totalRenewals + 1);
      break;
    case 'BIWEEKLY':
      totalRenewals = differenceInWeeks(now, renewalStart) / 2;
      nextRenewalDate = addWeeks(renewalStart, 2 * (totalRenewals + 1));
      break;
    case 'MONTHLY':
      totalRenewals = differenceInMonths(now, renewalStart);
      nextRenewalDate = addMonths(renewalStart, totalRenewals + 1);
      break;
    case 'QUARTERLY':
      totalRenewals = differenceInMonths(now, renewalStart) / 3;
      nextRenewalDate = addMonths(renewalStart, 3 * (totalRenewals + 1));
      break;
    case 'BIANNUALLY':
      totalRenewals = differenceInMonths(now, renewalStart) / 6;
      nextRenewalDate = addMonths(renewalStart, 6 * (totalRenewals + 1));
      break;
    case 'ANNUALLY':
      totalRenewals = differenceInYears(now, renewalStart);
      nextRenewalDate = addYears(renewalStart, totalRenewals + 1);
      break;
    default:
      throw new Error('Unrecognized renewal frequency');
  }

  const distanceToNextRenewal = formatDistanceToNowStrict(nextRenewalDate);

  return distanceToNextRenewal?.split(' ');
}

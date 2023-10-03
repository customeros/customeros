import differenceInMonths from 'date-fns/differenceInMonths';
import differenceInDays from 'date-fns/differenceInDays';
import differenceInWeeks from 'date-fns/differenceInWeeks';

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

export function getTimeToRenewal(renewalDate: string) {
  const now = new Date();
  const next = new Date(renewalDate);

  const months = differenceInMonths(next, now);
  const weeks = differenceInWeeks(next, now);
  const days = differenceInDays(next, now);

  if (days === 0) return ['0', 'days'];

  if (days === 1) return [days, 'day'];
  if (days < 7 && days !== 1) return [days, 'days'];

  if (weeks === 1) return [weeks, 'week'];
  if (weeks <= 4 && weeks !== 1 && months === 0) return [weeks, 'weeks'];
  if (weeks % 4 === 0 && weeks / 4 !== 1) return [weeks / 4, 'months'];

  if (months === 1 && weeks % 4 === 0) return [months, 'month'];

  const roundedMonths = weeks % 4 > 2 ? months + 1 : months;
  return [roundedMonths, 'months'];
}

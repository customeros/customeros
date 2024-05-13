import { set } from 'date-fns/set';
import { differenceInDays } from 'date-fns/differenceInDays';
import { differenceInYears } from 'date-fns/differenceInYears';
import { differenceInWeeks } from 'date-fns/differenceInWeeks';
import { differenceInHours } from 'date-fns/differenceInHours';
import { differenceInMonths } from 'date-fns/differenceInMonths';
import { differenceInMinutes } from 'date-fns/differenceInMinutes';

export function getDifferenceFromNow(targetDate: string) {
  const now = set(new Date(), { hours: 0, minutes: 0, seconds: 0 });
  const next = set(new Date(targetDate), { hours: 0, minutes: 0, seconds: 1 });

  const years = differenceInYears(next, now);
  const months = differenceInMonths(next, now);
  const monthsAfterYears = months - years * 12;
  const weeks = differenceInWeeks(next, now);
  const days = differenceInDays(next, now);

  if (days === 0) return ['', 'today'];
  if (days === 1) return ['1', 'day'];
  if (days < 7) return [`${days}`, 'days'];

  if (weeks === 1) return ['1', 'week'];
  if (weeks < 4) return [`${weeks}`, 'weeks'];

  if (years === 0) {
    if (monthsAfterYears === 1 || weeks === 4) return ['1', 'month'];

    return [months, 'months'];
  }

  if (years === 1) {
    if (monthsAfterYears === 0) return ['1', 'year'];

    return ['1', 'year', `${monthsAfterYears}`, 'months'];
  }

  if (monthsAfterYears === 0) return [`${years}`, 'years'];

  return [`${years}`, 'years', `${monthsAfterYears}`, 'months'];
}

export function getDifferenceInMinutesOrHours(targetDate: string) {
  const now = new Date();
  const next = new Date(targetDate);

  const minutes = Math.abs(differenceInMinutes(next, now));
  const hours = Math.abs(differenceInHours(next, now));

  if (minutes === 0) return ['1', 'minute'];
  if (minutes === 1) return [minutes, 'minute'];
  if (minutes < 60 && minutes > 1) return [minutes, 'minutes'];

  if (hours === 1) return [hours, 'hour'];
  if (hours <= 24 && hours > 1) return [hours, 'hours'];

  return [hours, 'hours'];
}

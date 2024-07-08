import { format, toZonedTime } from 'date-fns-tz';
import { differenceInHours } from 'date-fns/differenceInHours';
import {
  set,
  formatRFC3339,
  formatDistanceToNow,
  differenceInMinutes,
  FormatDurationOptions,
  isPast as isPastDateFns,
  isToday as isTodayDateFns,
  addDays as addDaysDateFns,
  formatDistanceToNowStrict,
  isBefore as isBeforeDateFns,
  isFuture as isFutureDateFns,
  addYears as addYearsDateFns,
  addMonths as addMonthsDateFns,
  isSameDay as isSameDayDateFns,
  isTomorrow as isTomorrowDateFns,
  formatDuration as formatDurationDateFns,
  differenceInDays as differenceInDaysDateFns,
  differenceInWeeks as differenceInWeeksDateFns,
  differenceInYears as differenceInYearsDateFns,
  differenceInMonths as differenceInMonthsDateFns,
} from 'date-fns';

export type ReturnDifference =
  | [null, 'today']
  | [number, 'day']
  | [number, 'days']
  | [number, 'week']
  | [number, 'weeks']
  | [number, 'month']
  | [number, 'months']
  | [number, 'year']
  | [number, 'years']
  | [number, 'year', string, 'months']
  | [1, 'year', string, 'months']
  | [number, 'years', string, 'months'];

export class DateTimeUtils {
  private static defaultFormatString = "EEE dd MMM - HH'h' mm zzz"; // Output: "Wed 08 Mar - 14h30CET"
  public static dateWithFullMonth = 'd MMMM yyyy'; // Output: "1 August 2024"
  public static defaultFormatShortString = 'dd MMM `yy'; // Output: "Wed 08 Mar - 14h30CET"
  public static dateWithHour = 'd MMM yyyy • HH:mm'; // Output: "19 Jun 2023 • 14:34"
  public static date = 'd MMM yyyy'; // Output: "19 Jun 2023"
  public static dateWithAbreviatedMonth = 'd MMM yyyy'; // Output: "1 Aug 2024"
  public static iso8601 = 'yyyy-MM-dd'; // Output:  "2024-07-08"
  public static dateWithShortYear = 'd MMM yy'; // Output: "1 Aug '24"
  public static dateDayAndMonth = 'd MMM'; // Output: "1 Aug"
  public static abreviatedMonth = 'MMM'; // Output: "Aug"
  public static shortWeekday = 'iiiiii'; // Output: "We"
  public static longWeekday = 'iiii'; // Output: "Wednesday"
  public static defaultTimeFormatString = 'HH:mm';
  public static dateTimeWithGMT = 'd MMM yyyy • Kbbb (z)'; // Output: "19 Jun 2023 • 2pm GMT"
  public static timeWithGMT = 'Kbbb (z)'; // Output: "2pm GMT"
  public static usaTimeFormatString = 'Kbbb';
  public static defaultDurationFormat = { format: ['minutes'] };

  private static getDate(date: string | number): Date {
    return new Date(new Date(date).toUTCString());
  }

  public static format(date: string | number, formatString?: string): string {
    const formatStr = formatString || this.defaultFormatString;

    return date ? format(this.getDate(date), formatStr) : '';
  }

  public static formatTime(
    date: string | number,
    formatString?: string,
  ): string {
    const formatStr = formatString || this.defaultTimeFormatString;

    return date ? format(this.getDate(date), formatStr) : '';
  }

  public static timeAgo(
    date: string | number,
    options?: {
      strict?: boolean;
      addSuffix?: boolean;
      includeMin?: boolean;
      includeSeconds?: boolean;
    },
  ): string {
    const isToday = this.isToday(this.getDate(date).toISOString());
    if (isToday && !options?.includeMin) {
      return 'today';
    }

    const formatter = options?.strict
      ? formatDistanceToNowStrict
      : formatDistanceToNow;

    return formatter(this.getDate(date), options);
  }

  public static isBeforeNow(date: string | number): boolean {
    return isBeforeDateFns(new Date(), new Date(date));
  }

  public static isBefore(dateLeft: string, dateRight: string): boolean {
    return isBeforeDateFns(this.getDate(dateLeft), this.getDate(dateRight));
  }

  public static toHoursAndMinutes(totalSeconds: number) {
    const totalMinutes = Math.floor(totalSeconds / 60);

    const seconds = totalSeconds % 60;
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    return { hours, minutes, seconds };
  }

  public static formatSecondsDuration(
    seconds: number,
    options?: { format: string[] },
  ): string {
    if (seconds === 0) {
      return '0 seconds';
    }

    const duration = this.toHoursAndMinutes(seconds);

    return formatDurationDateFns(duration, options as FormatDurationOptions);
  }

  public static isFuture(date: string): boolean {
    return isFutureDateFns(this.getDate(date));
  }

  public static isPast(date: string): boolean {
    return isPastDateFns(this.getDate(date));
  }
  public static isToday(date: string): boolean {
    return isTodayDateFns(this.getDate(date));
  }

  public static isTomorrow(date: string): boolean {
    return isTomorrowDateFns(this.getDate(date));
  }

  public static addYears(date: string, yearsCount: number): Date {
    return addYearsDateFns(this.getDate(date), yearsCount);
  }

  public static addMonth(date: string, monthsCount: number): Date {
    return addMonthsDateFns(this.getDate(date), monthsCount);
  }
  public static addDays(date: string, daysCount: number): Date {
    return addDaysDateFns(this.getDate(date), daysCount);
  }

  public static isSameDay(dateLeft: string, dateRight: string): boolean {
    return isSameDayDateFns(this.getDate(dateLeft), this.getDate(dateRight));
  }
  public static differenceInMins(dateLeft: string, dateRight: string): number {
    return differenceInMinutes(this.getDate(dateLeft), this.getDate(dateRight));
  }

  public static differenceInMonths(
    dateLeft: string,
    dateRight: string,
  ): number {
    return differenceInMonthsDateFns(
      this.getDate(dateLeft),
      this.getDate(dateRight),
    );
  }

  public static differenceInDays(dateLeft: string, dateRight: string): number {
    return differenceInDaysDateFns(
      this.getDate(dateLeft),
      this.getDate(dateRight),
    );
  }
  public static differenceInYears(dateLeft: string, dateRight: string): number {
    return differenceInYearsDateFns(
      this.getDate(dateLeft),
      this.getDate(dateRight),
    );
  }

  public static convertToTimeZone(
    date: string | Date,
    formatString: string,
    timeZone?: string,
  ) {
    const _date = typeof date === 'string' ? new Date(date) : date;
    const zonedDateStr = timeZone ? toZonedTime(date, timeZone ?? '') : null;

    return format(zonedDateStr ?? _date, formatString, {
      timeZone: timeZone || undefined,
    });
  }

  public static toISOMidnight(date: string | Date): string {
    const dateAtMidnight = set(
      typeof date === 'string' ? this.getDate(date) : date,
      {
        hours: 0,
        minutes: 0,
        seconds: 0,
        milliseconds: 0,
      },
    );

    return formatRFC3339(dateAtMidnight);
  }

  public static getUTCDateAtMidnight(date: Date): string {
    const utcDate = new Date(
      Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()),
    );

    return utcDate.toISOString().split('T')[0] + 'T00:00:00.000Z';
  }

  public static getDifferenceFromNow(targetDate: string): ReturnDifference {
    const now = set(new Date(), { hours: 0, minutes: 0, seconds: 0 });
    const next = set(new Date(targetDate), {
      hours: 0,
      minutes: 0,
      seconds: 1,
    });

    const years = differenceInYearsDateFns(next, now);
    const months = differenceInMonthsDateFns(next, now);
    const monthsAfterYears = months - years * 12;
    const weeks = differenceInWeeksDateFns(next, now);
    const days = differenceInDaysDateFns(next, now);

    if (days === 0) return [null, 'today'];
    if (days === 1) return [1, 'day'];
    if (days < 7) return [days, 'days'];

    if (weeks === 1) return [1, 'week'];
    if (weeks < 4) return [weeks, 'weeks'];

    if (years === 0) {
      if (monthsAfterYears === 1 || weeks === 4) return [1, 'month'];

      return [months, 'months'];
    }

    if (years === 1) {
      if (monthsAfterYears === 0) return [1, 'year'];

      return [1, 'year', `${monthsAfterYears}`, 'months'];
    }

    if (monthsAfterYears === 0) return [years, 'years'];

    return [years, 'years', `${monthsAfterYears}`, 'months'];
  }
  public static getDifferenceInMinutesOrHours(targetDate: string) {
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
}

import { format, utcToZonedTime } from 'date-fns-tz';
import {
  set,
  formatRFC3339,
  formatDistanceToNow,
  differenceInMinutes,
  isPast as isPastDateFns,
  isToday as isTodayDateFns,
  addDays as addDaysDateFns,
  isBefore as isBeforeDateFns,
  isFuture as isFutureDateFns,
  addYears as addYearsDateFns,
  addMonths as addMonthsDateFns,
  isSameDay as isSameDayDateFns,
  isTomorrow as isTomorrowDateFns,
  formatDuration as formatDurationDateFns,
  differenceInDays as differenceInDaysDateFns,
  differenceInMonths as differenceInMonthsDateFns,
} from 'date-fns';

export class DateTimeUtils {
  private static defaultFormatString = "EEE dd MMM - HH'h' mm zzz"; // Output: "Wed 08 Mar - 14h30CET"
  public static dateWithFullMonth = 'd MMMM yyyy'; // Output: "1 August 2024"
  public static defaultFormatShortString = 'dd MMM `yy'; // Output: "Wed 08 Mar - 14h30CET"
  public static dateWithHour = 'd MMM yyyy • HH:mm'; // Output: "19 Jun 2023 • 14:34"
  public static date = 'd MMM yyyy'; // Output: "19 Jun 2023"
  public static dateWithAbreviatedMonth = 'd MMM yyyy'; // Output: "1 Aug 2024"
  public static dateWithShortYear = 'd MMM yy'; // Output: "1 Aug '24"
  public static dateDayAndMonth = 'd MMM'; // Output: "1 Aug"
  public static abreviatedMonth = 'MMM'; // Output: "Aug"
  public static shortWeekday = 'iiiiii'; // Output: "We"
  public static longWeekday = 'iiii'; // Output: "Wednesday"
  public static defaultTimeFormatString = 'HH:mm';
  public static dateTimeWithGMT = 'd MMM yyyy • Kbbb (z)'; // Output: "19 Jun 2023 • 2pm GMT"
  public static timeWithGMT = 'Kbbb (z)'; // Output: "2pm GMT"
  public static usaTimeFormatString = 'Kbbb';
  private static defaultDurationFormat = { format: ['minutes'] };

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
      addSuffix?: boolean;
      includeMin?: boolean;
      includeSeconds?: boolean;
    },
  ): string {
    const isToday = this.isToday(this.getDate(date).toISOString());
    if (isToday && !options?.includeMin) {
      return 'today';
    }

    return formatDistanceToNow(this.getDate(date), options);
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

    return formatDurationDateFns(duration, options);
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

  public static addMonth(date: string, yearsCount: number): Date {
    return addMonthsDateFns(this.getDate(date), yearsCount);
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

  public static convertToTimeZone(
    date: string | Date,
    formatString: string,
    timeZone?: string,
  ) {
    const _date = typeof date === 'string' ? new Date(date) : date;
    const zonedDateStr = timeZone ? utcToZonedTime(date, timeZone ?? '') : null;

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
}

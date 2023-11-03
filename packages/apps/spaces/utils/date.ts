import {
  differenceInMinutes,
  formatDistanceToNow,
  formatDuration as formatDurationDateFns,
  isBefore,
  isSameDay as isSameDayDateFns,
} from 'date-fns';
import { format, utcToZonedTime } from 'date-fns-tz';

export class DateTimeUtils {
  private static defaultFormatString = "EEE dd MMM - HH'h' mm zzz"; // Output: "Wed 08 Mar - 14h30CET"
  public static dateWithFullMonth = 'd MMMM yyyy'; // Output: "1 August 2024"
  public static defaultFormatShortString = 'dd MMM `yy'; // Output: "Wed 08 Mar - 14h30CET"
  public static dateWithHour = 'd MMM yyyy • HH:mm'; // Output: "19 Jun 2023 • 14:34"
  public static date = 'd MMM yyyy'; // Output: "19 Jun 2023"
  public static dateWithAbreviatedMonth = 'd MMM yyyy'; // Output: "1 Aug 2024"
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
    options?: { includeSeconds?: boolean; addSuffix?: boolean },
  ): string {
    return formatDistanceToNow(this.getDate(date), options);
  }

  public static isBeforeNow(date: string | number): boolean {
    return isBefore(new Date(), new Date(date));
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

  public static isSameDay(dateLeft: string, dateRight: string): boolean {
    return isSameDayDateFns(this.getDate(dateLeft), this.getDate(dateRight));
  }
  public static differenceInMins(dateLeft: string, dateRight: string): number {
    return differenceInMinutes(this.getDate(dateLeft), this.getDate(dateRight));
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
}

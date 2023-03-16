import { formatDistanceToNow } from 'date-fns';
import { format } from 'date-fns-tz';

export class DateTimeUtils {
  private static defaultFormatString = "EEE dd MMM - HH'h' mm zzz"; // Output: "Wed 08 Mar - 14h30CET"

  private static getDate(date: string | number): Date {
    return new Date(new Date(date).toUTCString());
  }
  public static format(date: string | number, formatString?: string): string {
    const formatStr = formatString || this.defaultFormatString;

    return format(this.getDate(date), formatStr);
  }

  public static timeAgo(
    date: string | number,
    options?: { includeSeconds?: boolean; addSuffix?: boolean },
  ): string {
    return formatDistanceToNow(this.getDate(date), options);
  }
}

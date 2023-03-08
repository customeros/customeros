import { formatDistanceToNow } from 'date-fns';
import { format } from 'date-fns-tz';

export class DateTimeUtils {
  private static defaultFormatString = "EEE dd MMM - HH'h'mm zzz"; // Output: "Wed 08 Mar - 14h30CET"

  public static format(date: Date, formatString?: string): string {
    const formatStr = formatString || this.defaultFormatString;

    return format(date, formatStr);
  }

  public static timeAgo(
    date: Date,
    options?: { includeSeconds?: boolean; addSuffix?: boolean },
  ): string {
    return formatDistanceToNow(date, options);
  }
}

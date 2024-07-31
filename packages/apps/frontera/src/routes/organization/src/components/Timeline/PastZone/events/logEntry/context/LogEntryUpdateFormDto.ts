import { DateTimeUtils } from '@utils/date';
import { LogEntryUpdateInput } from '@graphql/types';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';

export interface LogEntryUpdateFormDtoI {
  time: string;
  content: string;
  date: Date | string;
  contentType: string;
}

export interface LogEntryUpdateForm {
  time: string;
  content: string;
  date: Date | string;
  contentType: string;
}

export class LogEntryUpdateFormDto implements LogEntryUpdateForm {
  date: Date | string;
  time: string;
  content: string;
  contentType: string;

  constructor(data?: LogEntryWithAliases) {
    this.content = data?.content ?? '';
    this.date = data?.logEntryStartedAt || new Date();
    this.contentType = data?.contentType ?? '';
    this.time = data?.logEntryStartedAt
      ? DateTimeUtils.formatTime(data?.logEntryStartedAt)
      : DateTimeUtils.formatTime(new Date().toISOString());
  }

  private static applyHourAndMinuteToDate(
    date: Date | string,
    time: string,
  ): Date {
    const timeArray = time?.split(':');
    const newDate = new Date(date); // Create a new Date object to maintain immutability

    newDate.setHours(Number(timeArray?.[0] ?? '00'));
    newDate.setMinutes(Number(timeArray?.[1] ?? '00'));

    return newDate;
  }

  static toPayload(
    data: LogEntryUpdateFormDtoI & { contentType: string },
  ): LogEntryUpdateInput {
    return {
      startedAt: this.applyHourAndMinuteToDate(data.date, data.time),
      contentType: data.contentType,
    };
  }
}

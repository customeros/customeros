import { LogEntryUpdateInput } from '@graphql/types';
import { DateTimeUtils } from '@spaces/utils/date';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';

export interface LogEntryUpdateFormDtoI {
  date: Date | string;
  time: string;
}

export interface LogEntryUpdateForm {
  date: Date | string;
  time: string;
}

export class LogEntryUpdateFormDto implements LogEntryUpdateForm {
  date: Date | string;
  time: string;

  constructor(data?: LogEntryWithAliases) {
    this.date = data?.logEntryStartedAt || new Date();
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
    data: LogEntryUpdateForm,
  ): Pick<LogEntryUpdateInput, 'startedAt'> {
    return {
      startedAt: this.applyHourAndMinuteToDate(data.date, data.time),
    };
  }
}

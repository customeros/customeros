import { LogEntry, LogEntryUpdateInput } from '@graphql/types';

export interface LogEntryFormDtoI {
  content: string;
  appSource: string;
  contentType: string;
  startedAt: string | Date;
  tags: Array<{ label: string; value: string }>;
}

export interface LogEntryForm {
  content: string;
  contentType: string;
  tags: Array<{ label: string; value: string }>;
}

export class LogEntryFormDto implements LogEntryForm {
  tags: Array<{ label: string; value: string }>;
  content: string;
  contentType: string;
  appSource: string;
  startedAt: string | Date;

  constructor(data?: LogEntry) {
    this.content = data?.content || '';
    this.contentType = 'text/html';
    this.tags = [];
    this.appSource = 'customerOS';
    this.startedAt = new Date().toISOString();
  }

  // Commented out as at this moment user cannot update existing data besides started at. Therefore this is not needed for now
  // static toForm(data: LogEntryDtoI) {
  //   return new LogEntryDto(data);
  // }

  static toPayload(data: LogEntryForm): LogEntryUpdateInput {
    return {
      tags: data.tags.map((data) => ({ name: data?.label })),
      content: data.content,
      contentType: data.contentType,
      appSource: 'customerOS',
      startedAt: new Date().toISOString(),
    } as unknown as LogEntryUpdateInput;
  }
}

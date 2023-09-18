import { TagIdOrNameInput } from '@graphql/types';

export interface LogEntryDtoI {
  tags: Array<TagIdOrNameInput>;
  content: string;
  contentType: string;
  appSource: string;
  startedAt: string | Date;
}

export interface LogEntryDtoFormI {
  tags: Array<{ label: string; value: string }>;
  content: string;
  contentType: string;
}

export class LogEntryDto implements LogEntryDtoI {
  tags: Array<TagIdOrNameInput>;
  content: string;
  contentType: string;
  appSource: string;
  startedAt: string | Date;

  constructor(data?: any) {
    this.content = data?.content || '';
    this.contentType = 'text/html';
    this.tags = [];
    this.appSource = 'customerOS';
    this.startedAt = new Date().toISOString();
  }

  static toForm(data: any) {
    return new LogEntryDto(data);
  }

  static toPayload(data: LogEntryDtoFormI) {
    return {
      tags: data.tags.map(({ label, value }) => ({ name: label })),
      content: data.content,
      contentType: data.contentType,
      appSource: 'customerOS',
      startedAt: new Date().toISOString(),
    } as any;
  }
}

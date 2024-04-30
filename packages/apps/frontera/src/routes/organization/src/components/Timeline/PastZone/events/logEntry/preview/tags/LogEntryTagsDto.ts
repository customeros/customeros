import { LogEntry, TagIdOrNameInput } from '@graphql/types';

export interface LogEntryTagsFormDtoI {
  tags: Array<{ label: string; value: string }>;
}

export interface LogEntryTagsForm {
  tags: Array<{ label: string; value: string }>;
}

export interface UpdateTagsInput {
  tags: Array<TagIdOrNameInput>;
}

export class LogEntryTagsDto implements LogEntryTagsForm {
  tags: Array<{ label: string; value: string }>;

  constructor(data?: Pick<LogEntry, 'tags'>) {
    this.tags = (data?.tags ?? [])?.map((e) => ({
      label: e.name,
      value: e.id,
    }));
  }

  static toPayload(data: LogEntryTagsFormDtoI): UpdateTagsInput {
    return {
      tags: data.tags.map((data) => ({ name: data?.label })),
    };
  }
}

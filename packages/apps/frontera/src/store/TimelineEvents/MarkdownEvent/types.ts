import { Metadata, MarkdownEvent } from '@graphql/types';

export type MarkdownEventType = MarkdownEvent & {
  markdownEventMetadata: Metadata;
};

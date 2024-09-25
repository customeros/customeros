import { LinkNode, AutoLinkNode } from '@lexical/link';
import { QuoteNode, HeadingNode } from '@lexical/rich-text';

import { MentionNode } from './MentionNode';
import { HashtagNode } from './HashtagNode';

export const nodes = [
  LinkNode,
  AutoLinkNode,
  HashtagNode,
  MentionNode,
  HashtagNode,
  HeadingNode,
  QuoteNode,
];

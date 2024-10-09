import { TextNode } from 'lexical';
import { CodeNode } from '@lexical/code';
import { LinkNode, AutoLinkNode } from '@lexical/link';
import { ListNode, ListItemNode } from '@lexical/list';
import { QuoteNode, HeadingNode } from '@lexical/rich-text';

import { MentionNode } from './MentionNode';
import { HashtagNode } from './HashtagNode';
import { VariableNode } from './VariableNode';
import { ExtendedTextNode } from './ExtendedTextNode';
import { ExtendedQuoteNode } from './ExtendedQuoteNode';

export const nodes = [
  LinkNode,
  AutoLinkNode,
  HashtagNode,
  VariableNode,
  MentionNode,
  HashtagNode,
  ExtendedTextNode,
  HeadingNode,
  CodeNode,
  ListNode,
  ListItemNode,
  ExtendedQuoteNode,
  {
    replace: TextNode,
    with: (node: TextNode) => new ExtendedTextNode(node.__text),
    withKlass: ExtendedTextNode,
  },
  {
    replace: QuoteNode,
    with: () => new ExtendedQuoteNode(),
    withKlass: ExtendedQuoteNode,
  },
];

import data from 'svgmoji/emoji-slack.json';
import { IdentifierSchemaAttributes } from 'remirror';
import {
  BoldExtension,
  LinkExtension,
  EmojiExtension,
  ItalicExtension,
  StrikeExtension,
  EventsExtension,
  HeadingExtension,
  FontSizeExtension,
  MarkdownExtension,
  UnderlineExtension,
  BlockquoteExtension,
  BulletListExtension,
  FontFamilyExtension,
  OrderedListExtension,
  MentionAtomExtension,
  NodeFormattingExtension,
} from 'remirror/extensions';

export const basicEditorExtensions = () => [
  new ItalicExtension(),
  new BoldExtension(),
  new StrikeExtension(),
  new UnderlineExtension(),
  new OrderedListExtension(),
  new BulletListExtension(),
  new FontSizeExtension(),
  new FontFamilyExtension(),
  new BlockquoteExtension(),
  new HeadingExtension(),
  new MarkdownExtension(),
  new NodeFormattingExtension(),
  new EventsExtension(),
  new EmojiExtension({ data, moji: 'noto', fallback: '', plainText: true }),
  new MentionAtomExtension({
    matchers: [{ name: 'tag', char: '#', mentionClassName: 'customeros-tag' }],
  }),
  new LinkExtension({ autoLink: true }),
];

export const extraAttributes: IdentifierSchemaAttributes[] = [
  {
    identifiers: ['mention', 'emoji'],
    attributes: { role: { default: 'presentation' } },
  },
  { identifiers: ['mention'], attributes: { href: { default: `/` } } },
];

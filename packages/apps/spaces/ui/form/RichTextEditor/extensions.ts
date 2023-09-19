import {
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  FontFamilyExtension,
  FontSizeExtension,
  HeadingExtension,
  ItalicExtension,
  LinkExtension,
  NodeFormattingExtension,
  OrderedListExtension,
  StrikeExtension,
  UnderlineExtension,
  MentionAtomExtension,
  EmojiExtension,
  MarkdownExtension,
  EventsExtension,
} from 'remirror/extensions';
import data from 'svgmoji/emoji-slack.json';
import { IdentifierSchemaAttributes } from 'remirror';

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
    matchers: [
      { name: 'at', char: '@' },
      { name: 'tag', char: '#', mentionClassName: 'customeros-tag' },
    ],
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

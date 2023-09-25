import {
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  EmojiExtension,
  EventsExtension,
  ItalicExtension,
  LinkExtension,
  MarkdownExtension,
  MentionAtomExtension,
  NodeFormattingExtension,
  OrderedListExtension,
  StrikeExtension,
  UnderlineExtension,
} from 'remirror/extensions';
import data from 'svgmoji/emoji-slack.json';

export const logEntryEditorExtensions = () => [
  new ItalicExtension(),
  new BoldExtension(),
  new StrikeExtension(),
  new UnderlineExtension(),
  new OrderedListExtension(),
  new BulletListExtension(),
  new BlockquoteExtension(),
  new MarkdownExtension(),
  new NodeFormattingExtension(),
  new EventsExtension(),
  new EmojiExtension({ data, moji: 'noto', fallback: '', plainText: true }),
  new MentionAtomExtension({
    matchers: [{ name: 'tag', char: '#', mentionClassName: 'customeros-tag' }],
  }),
  new LinkExtension({ autoLink: true }),
];

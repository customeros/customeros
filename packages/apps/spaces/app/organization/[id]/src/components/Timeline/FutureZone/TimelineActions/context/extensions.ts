import data from 'svgmoji/emoji-slack.json';
import {
  BoldExtension,
  LinkExtension,
  EmojiExtension,
  EventsExtension,
  ItalicExtension,
  StrikeExtension,
  MarkdownExtension,
  UnderlineExtension,
  BlockquoteExtension,
  BulletListExtension,
  MentionAtomExtension,
  OrderedListExtension,
  NodeFormattingExtension,
} from 'remirror/extensions';

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
    matchers: [
      { name: 'tag', char: '#', mentionClassName: 'customeros-tag' },
      { name: 'at', char: '@', mentionClassName: 'customeros-mention' },
    ],
  }),
  new LinkExtension({ autoLink: true }),
];

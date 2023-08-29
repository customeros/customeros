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
  new NodeFormattingExtension(),
  new LinkExtension({ autoLink: true }),
];

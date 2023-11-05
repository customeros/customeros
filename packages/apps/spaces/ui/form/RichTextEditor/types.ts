import { AnyExtension } from 'remirror';
import { ReactExtensions, UseRemirrorReturn } from '@remirror/react';
import {
  BoldExtension,
  LinkExtension,
  EmojiExtension,
  ItalicExtension,
  StrikeExtension,
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

export type BasicEditorExtentions =
  | ItalicExtension
  | BoldExtension
  | StrikeExtension
  | UnderlineExtension
  | OrderedListExtension
  | NodeFormattingExtension
  | BlockquoteExtension
  | BulletListExtension
  | FontFamilyExtension
  | FontSizeExtension
  | LinkExtension
  | HeadingExtension
  | MentionAtomExtension
  | EmojiExtension
  | MarkdownExtension;

export type RemirrorProps<T extends AnyExtension> = UseRemirrorReturn<
  ReactExtensions<T>
>;

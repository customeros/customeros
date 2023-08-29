import { ReactExtensions, UseRemirrorReturn } from '@remirror/react';
import {
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  FontFamilyExtension,
  FontSizeExtension,
  HeadingExtension,
  ItalicExtension,
  NodeFormattingExtension,
  OrderedListExtension,
  StrikeExtension,
  UnderlineExtension,
} from 'remirror/extensions';
import { AnyExtension } from 'remirror';

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
  | HeadingExtension;

export type RemirrorProps<T extends AnyExtension> = UseRemirrorReturn<
  ReactExtensions<T>
>;

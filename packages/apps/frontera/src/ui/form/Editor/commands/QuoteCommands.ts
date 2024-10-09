import { $setBlocksType } from '@lexical/selection';
import {
  $insertNodes,
  $getSelection,
  LexicalEditor,
  createCommand,
  $isRangeSelection,
  KEY_ENTER_COMMAND,
  $createParagraphNode,
  COMMAND_PRIORITY_LOW,
  COMMAND_PRIORITY_CRITICAL,
} from 'lexical';

import {
  $isExtendedQuoteNode,
  $createExtendedQuoteNode,
} from './../nodes/ExtendedQuoteNode';

function toggleQuote(editor: LexicalEditor): boolean {
  let wasToggled = false;

  editor.update(() => {
    const selection = $getSelection();

    if (!$isRangeSelection(selection)) {
      return;
    }

    const anchorNode = selection.anchor.getNode();
    const focusNode = selection.focus.getNode();

    const isEntirelyWithinQuote =
      $isExtendedQuoteNode(anchorNode.getParent()) &&
      $isExtendedQuoteNode(focusNode.getParent()) &&
      anchorNode.getParent() === focusNode.getParent();

    if (isEntirelyWithinQuote) {
      $setBlocksType(selection, $createParagraphNode);
      wasToggled = true;
    } else {
      $setBlocksType(selection, $createExtendedQuoteNode);
      wasToggled = true;
    }
  });

  return wasToggled;
}
export const TOGGLE_BLOCKQUOTE_COMMAND = createCommand(
  'TOGGLE_BLOCKQUOTE_COMMAND',
);

export function registerToggleQuoteCommand(editor: LexicalEditor): () => void {
  return editor.registerCommand(
    TOGGLE_BLOCKQUOTE_COMMAND,
    () => {
      return toggleQuote(editor);
    },
    COMMAND_PRIORITY_LOW,
  );
}

export function registerEnterQuoteCommand(editor: LexicalEditor): () => void {
  return editor.registerCommand(
    KEY_ENTER_COMMAND,
    (event) => {
      const selection = $getSelection();

      if (!$isRangeSelection(selection) || !event) {
        return false;
      }

      const anchorNode = selection.anchor.getNode();
      const topLevelElement = anchorNode.getTopLevelElement();

      if ($isExtendedQuoteNode(topLevelElement)) {
        event.preventDefault();

        if (topLevelElement.getTextContent().trim() === '') {
          editor.update(() => {
            const paragraph = $createParagraphNode();

            topLevelElement.replace(paragraph);
            paragraph.select();
          });
        } else if (
          anchorNode.getTextContent().length === selection.anchor.offset
        ) {
          editor.update(() => {
            const newQuote = $createExtendedQuoteNode();

            topLevelElement.insertAfter(newQuote);
            newQuote.select();
          });
        } else {
          editor.update(() => {
            const newQuote = $createExtendedQuoteNode();

            $insertNodes([newQuote]);

            newQuote.select();
          });
        }

        return true;
      }

      return false;
    },
    COMMAND_PRIORITY_CRITICAL,
  );
}

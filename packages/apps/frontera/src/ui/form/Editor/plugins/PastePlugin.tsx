import { useEffect, ClipboardEvent } from 'react';

import { $createLinkNode } from '@lexical/link';
import { $generateNodesFromDOM } from '@lexical/html';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $getSelection,
  PASTE_COMMAND,
  $createTextNode,
  $isRangeSelection,
} from 'lexical';

function isValidUrl(string: string) {
  try {
    new URL(string);

    return true;
  } catch (_) {
    return false;
  }
}

export function LinkPastePlugin() {
  const [editor] = useLexicalComposerContext();

  useEffect(() => {
    const handlePaste = (event: ClipboardEvent) => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        const clipboardData = event.clipboardData;
        const pastedData = clipboardData.getData('text/plain');
        const selectedText = selection.getTextContent().trim();

        if (selectedText.length && isValidUrl(pastedData)) {
          editor.update(() => {
            const linkNode = $createLinkNode(pastedData);

            // Get the selected text
            const selectedText = selection.getTextContent();

            // If there's selected text, use it; otherwise, use the URL
            const textNode = $createTextNode(selectedText || pastedData);

            linkNode.append(textNode);

            // Replace the selection with the new link node
            selection.insertNodes([linkNode]);
          });
        } else {
          editor.update(() => {
            const htmlData = clipboardData.getData('text/html');

            if (htmlData) {
              const parser = new DOMParser();
              const doc = parser.parseFromString(htmlData, 'text/html');

              const nodes = $generateNodesFromDOM(editor, doc);

              selection.insertNodes(nodes);
            } else {
              const textNode = $createTextNode(pastedData);

              selection.insertNodes([textNode]);
            }
          });
        }
      }
    };

    editor.registerCommand(
      PASTE_COMMAND,
      (event: ClipboardEvent) => {
        handlePaste(event);

        return true;
      },
      1,
    );
  }, [editor]);

  return null;
}

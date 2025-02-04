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

import { isValidUrl } from '@utils/urlValidation';
import { convertPlainTextToHtml } from '@ui/form/Editor/utils/convertPlainTextToHtml';

export function LinkPastePlugin() {
  const [editor] = useLexicalComposerContext();

  useEffect(() => {
    const handlePaste = (event: ClipboardEvent) => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        const clipboardData = event.clipboardData;
        const pastedData = clipboardData?.getData('text/plain');
        const selectedText = selection.getTextContent().trim();

        if (!pastedData) return;

        if (selectedText.length && isValidUrl(pastedData)) {
          editor.update(() => {
            const linkNode = $createLinkNode(pastedData);
            const textNode = $createTextNode(selectedText || pastedData);

            linkNode.append(textNode);
            selection.insertNodes([linkNode]);
          });
        } else {
          editor.update(() => {
            const htmlData = clipboardData?.getData('text/html');

            if (htmlData) {
              const parser = new DOMParser();
              const doc = parser.parseFromString(htmlData, 'text/html');

              // strip all formatting enforced by pre tags
              doc.querySelectorAll('pre')?.forEach((preEl) => {
                if (!preEl.querySelector('code')) {
                  const span = doc.createElement('div');

                  // Move all child nodes into the new div
                  while (preEl.firstChild) {
                    span.appendChild(preEl.firstChild);
                  }
                  preEl.replaceWith(span);
                }
              });

              const nodes = $generateNodesFromDOM(editor, doc);

              selection.insertNodes(nodes);
            } else {
              const htmlData = convertPlainTextToHtml(pastedData);
              const parser = new DOMParser();
              const doc = parser.parseFromString(htmlData, 'text/html');
              const nodes = $generateNodesFromDOM(editor, doc);

              selection.insertNodes(nodes);
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

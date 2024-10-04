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
const ALLOWED_TAGS = [
  'ul',
  'ol',
  'p',
  'div',
  'span',
  'a',
  'strong',
  'em',
  's',
  'blockquote',
  'i',
  'b',
  'u',
  'li', // Additional text formatting tags
];

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
            const textNode = $createTextNode(selectedText || pastedData);

            linkNode.append(textNode);
            selection.insertNodes([linkNode]);
          });
        } else {
          editor.update(() => {
            const htmlData = clipboardData.getData('text/html');

            if (htmlData) {
              const parser = new DOMParser();
              const doc = parser.parseFromString(htmlData, 'text/html');

              // Filter out unsupported elements
              const filterNode = (node: Node): Node | null => {
                if (node.nodeType === Node.ELEMENT_NODE) {
                  const element = node as Element;

                  if (!ALLOWED_TAGS.includes(element.tagName.toLowerCase())) {
                    return document.createTextNode(element.textContent || '');
                  }
                  Array.from(element.childNodes).forEach((child) => {
                    const filteredChild = filterNode(child);

                    if (filteredChild) {
                      element.replaceChild(filteredChild, child);
                    } else {
                      element.removeChild(child);
                    }
                  });
                }

                return node;
              };

              doc.body.childNodes.forEach((child) => {
                filterNode(child);
              });

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

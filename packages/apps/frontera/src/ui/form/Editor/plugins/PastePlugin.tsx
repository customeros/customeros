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

              const filterNode = (node: Node): DocumentFragment => {
                const fragment = document.createDocumentFragment();

                if (node.nodeType === Node.TEXT_NODE) {
                  fragment.appendChild(node.cloneNode());
                } else if (node.nodeType === Node.ELEMENT_NODE) {
                  const element = node as HTMLElement;
                  const tagName = element.tagName.toLowerCase();

                  if (ALLOWED_TAGS.includes(tagName)) {
                    const newElement = document.createElement(tagName);

                    if (tagName === 'a') {
                      const href = element.getAttribute('href');

                      if (href) newElement.setAttribute('href', href);
                    }

                    Array.from(element.childNodes).forEach((child) => {
                      newElement.appendChild(filterNode(child));
                    });

                    fragment.appendChild(newElement);
                  } else {
                    Array.from(element.childNodes).forEach((child) => {
                      fragment.appendChild(filterNode(child));
                    });
                  }
                }

                return fragment;
              };

              const filteredBody = filterNode(doc.body);

              // Create a new document to hold our filtered content
              const newDoc = document.implementation.createHTMLDocument();

              newDoc.body.appendChild(filteredBody);

              const nodes = $generateNodesFromDOM(editor, newDoc);

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

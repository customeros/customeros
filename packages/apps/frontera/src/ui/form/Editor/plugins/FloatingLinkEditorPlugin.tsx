import * as React from 'react';
import { createPortal } from 'react-dom';
import { useRef, Dispatch, useState, useEffect, useCallback } from 'react';

import { mergeRegister, $findMatchingParent } from '@lexical/utils';
import { flip, shift, offset, computePosition } from '@floating-ui/dom';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $isLinkNode,
  $createLinkNode,
  TOGGLE_LINK_COMMAND,
} from '@lexical/link';
import {
  $getSelection,
  BaseSelection,
  CLICK_COMMAND,
  LexicalEditor,
  $setSelection,
  $isRangeSelection,
  KEY_ESCAPE_COMMAND,
  COMMAND_PRIORITY_LOW,
  COMMAND_PRIORITY_HIGH,
  $createRangeSelection,
  SELECTION_CHANGE_COMMAND,
  COMMAND_PRIORITY_NORMAL as NORMAL_PRIORITY,
  SELECTION_CHANGE_COMMAND as ON_SELECTION_CHANGE,
} from 'lexical';

import { Input } from '@ui/form/Input';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Divider } from '@ui/presentation/Divider';
import { getExternalUrl } from '@utils/getExternalLink.ts';
import { FloatingToolbarButton } from '@ui/form/Editor/components';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02.tsx';

import { sanitizeUrl } from '../utils/url';
import { getSelectedNode } from '../utils/getSelectedNode';

function FloatingLinkEditor({
  editor,
  isLink,
  setIsLink,
  anchorElem,
}: {
  isLink: boolean;
  editor: LexicalEditor;
  anchorElem: HTMLElement;
  setIsLink: Dispatch<boolean>;
}): JSX.Element {
  const editorRef = useRef<HTMLDivElement | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [linkUrl, setLinkUrl] = useState('');
  const [lastSelection, setLastSelection] = useState<BaseSelection | null>(
    null,
  );

  const updateLinkEditor = useCallback(() => {
    const selection = $getSelection();

    if ($isRangeSelection(selection)) {
      const node = getSelectedNode(selection);
      const linkParent = $findMatchingParent(node, $isLinkNode);

      if (linkParent) {
        setLinkUrl(linkParent.getURL());
      } else if ($isLinkNode(node)) {
        setLinkUrl(node.getURL());
      } else {
        setLinkUrl('');
      }
    }

    const editorElem = editorRef.current;
    const nativeSelection = window.getSelection();

    if (editorElem === null || !isLink) {
      return;
    }

    const rootElement = editor.getRootElement();

    if (
      selection !== null &&
      nativeSelection !== null &&
      rootElement !== null &&
      rootElement.contains(nativeSelection.anchorNode) &&
      editor.isEditable()
    ) {
      computePosition(anchorElem, editorElem, {
        placement: 'top-start',
        middleware: [offset(4), flip(), shift()],
      }).then(({ x, y }) => {
        editorElem.style.top = `${y}px`;
        editorElem.style.left = `${x}px`;
      });

      setLastSelection(selection);
    }
  }, [anchorElem, editor, isLink]);

  useEffect(() => {
    const scrollerElem = anchorElem.parentElement;

    const update = () => {
      editor.getEditorState().read(() => {
        updateLinkEditor();
      });
    };

    window.addEventListener('resize', update);

    if (scrollerElem) {
      scrollerElem.addEventListener('scroll', update);
    }

    return () => {
      window.removeEventListener('resize', update);

      if (scrollerElem) {
        scrollerElem.removeEventListener('scroll', update);
      }
    };
  }, [anchorElem.parentElement, editor, updateLinkEditor]);

  useEffect(() => {
    return mergeRegister(
      editor.registerUpdateListener(({ editorState }) => {
        editorState.read(() => {
          updateLinkEditor();
        });
      }),

      editor.registerCommand(
        SELECTION_CHANGE_COMMAND,
        () => {
          updateLinkEditor();

          return true;
        },
        COMMAND_PRIORITY_LOW,
      ),
      editor.registerCommand(
        KEY_ESCAPE_COMMAND,
        () => {
          if (isLink) {
            setIsLink(false);

            return true;
          }

          return false;
        },
        COMMAND_PRIORITY_HIGH,
      ),
    );
  }, [editor, updateLinkEditor, setIsLink, isLink]);

  useEffect(() => {
    editor.getEditorState().read(() => {
      updateLinkEditor();
    });
  }, [editor, updateLinkEditor]);
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        editorRef.current &&
        !editorRef.current.contains(event.target as Node)
      ) {
        setIsLink(false);
      }
    };

    if (isLink) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isLink, setIsLink]);

  const handleLinkSubmission = useCallback(() => {
    if (lastSelection !== null) {
      editor.update(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          const node = getSelectedNode(selection);
          const parent = node.getParent();

          if (linkUrl.trim() === '') {
            // Remove the link if URL is empty
            if ($isLinkNode(parent)) {
              const children = parent.getChildren();

              for (const child of children) {
                parent.insertBefore(child);
              }
              parent.remove();
            } else if ($isLinkNode(node)) {
              node.remove();
            }
          } else {
            // Existing logic for creating or updating links
            if ($isLinkNode(parent)) {
              parent.setURL(sanitizeUrl(linkUrl));
            } else if ($isLinkNode(node)) {
              node.setURL(sanitizeUrl(linkUrl));
            } else {
              const linkNode = $createLinkNode(sanitizeUrl(linkUrl));

              selection.insertNodes([linkNode]);
            }
          }
        }
      });
      setIsLink(false);
    }
  }, [editor, linkUrl, lastSelection, setIsLink]);

  const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
    if (event.key === 'Tab') {
      event.preventDefault();

      const focusableElements =
        editorRef.current?.querySelectorAll('input, button') || [];
      const firstElement = focusableElements[0] as HTMLElement;
      const lastElement = focusableElements[
        focusableElements.length - 1
      ] as HTMLElement;

      if (event.shiftKey) {
        if (document.activeElement === firstElement) {
          lastElement.focus();
        } else {
          (
            document.activeElement?.previousElementSibling as HTMLElement
          )?.focus();
        }
      } else {
        if (document.activeElement === lastElement) {
          firstElement.focus();
        } else {
          (document.activeElement?.nextElementSibling as HTMLElement)?.focus();
        }
      }
    }
  }, []);

  return createPortal(
    <div
      ref={editorRef}
      onKeyDown={handleKeyDown}
      className='absolute top-0 left-0 z-[99999] pointer-events-auto'
    >
      {isLink && (
        <div className='bg-gray-700 flex items-center min-w-[auto] max-w-[800px] p-1 pl-3 shadow-lg rounded-md'>
          <Input
            size='sm'
            ref={inputRef}
            value={linkUrl}
            variant='unstyled'
            placeholder='Enter a URL'
            className='leading-none min-h-0 pointer-events-auto text-gray-25 overflow-ellipsis'
            onChange={(event) => {
              setLinkUrl(event.target.value);
            }}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                event.preventDefault();
                handleLinkSubmission();
              }
            }}
          />

          <Divider className='w-[1px] h-3 border-b-0 border-l-[1px] border-gray-500 mx-2' />

          <FloatingToolbarButton
            aria-label='Open link'
            icon={<LinkExternal02 className='text-gray-500 stroke-1' />}
            onClick={() => {
              const link = getExternalUrl(sanitizeUrl(linkUrl));

              window.open(link, '_blank', 'noopener,noreferrer');
            }}
          />
          <FloatingToolbarButton
            aria-label='Delete link'
            icon={<Trash01 className='text-gray-500' />}
            onClick={() => {
              editor.dispatchCommand(TOGGLE_LINK_COMMAND, null);
              setIsLink(false);
            }}
          />
        </div>
      )}
    </div>,
    anchorElem,
  );
}

function useFloatingLinkEditorToolbar(
  editor: LexicalEditor,
  anchorElem: HTMLElement,
): JSX.Element | null {
  const [isLink, setIsLink] = useState(false);

  useEffect(() => {
    return mergeRegister(
      editor.registerCommand(
        CLICK_COMMAND,
        (e) => {
          const selection = $getSelection();

          if ($isRangeSelection(selection)) {
            const node = selection.getNodes()[0];
            const linkNode = $findMatchingParent(node, $isLinkNode);

            if ($isLinkNode(linkNode)) {
              const domSelection = window.getSelection();
              const domRange = domSelection?.getRangeAt(0);
              const clickOffset = domRange?.startOffset;

              // Check if the click is at the end of the link
              if (clickOffset === linkNode.getTextContent().length) {
                // Click is at the end, don't select or open toolbar
                return false;
              }

              // Click is not at the end, proceed with link selection
              e.preventDefault();
              setIsLink(true);

              const newSelection = $createRangeSelection();

              newSelection.anchor.set(linkNode.getKey(), 0, 'element');
              newSelection.focus.set(
                linkNode.getKey(),
                linkNode.getChildrenSize(),
                'element',
              );

              editor.update(() => {
                newSelection.dirty = true;
                $setSelection(newSelection);
              });

              return true;
            }
          }

          return false;
        },
        COMMAND_PRIORITY_LOW,
      ),
      editor.registerCommand(
        ON_SELECTION_CHANGE,
        () => {
          const selection = $getSelection();

          if (
            $isRangeSelection(selection) &&
            !selection.anchor.is(selection.focus)
          ) {
            const node = selection.getNodes()[0];
            const linkNode = $findMatchingParent(node, $isLinkNode);

            if ($isLinkNode(linkNode)) {
              return true;
            }
          }

          return false;
        },
        NORMAL_PRIORITY,
      ),
    );
  }, [editor]);

  return (
    <FloatingLinkEditor
      isLink={isLink}
      editor={editor}
      setIsLink={setIsLink}
      anchorElem={anchorElem}
    />
  );
}

export default function FloatingLinkEditorPlugin({
  anchorElem = document.body,
}: {
  anchorElem?: HTMLElement;
}): JSX.Element | null {
  const [editor] = useLexicalComposerContext();

  return useFloatingLinkEditorToolbar(editor, anchorElem);
}

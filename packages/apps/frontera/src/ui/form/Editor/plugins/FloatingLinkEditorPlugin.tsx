import * as React from 'react';
import { createPortal } from 'react-dom';
import { useRef, Dispatch, useState, useEffect, useCallback } from 'react';

import { mergeRegister, $findMatchingParent } from '@lexical/utils';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $isLinkNode,
  $createLinkNode,
  $isAutoLinkNode,
  TOGGLE_LINK_COMMAND,
} from '@lexical/link';
import {
  $getSelection,
  BaseSelection,
  CLICK_COMMAND,
  LexicalEditor,
  $isLineBreakNode,
  $isRangeSelection,
  KEY_ESCAPE_COMMAND,
  FORMAT_TEXT_COMMAND,
  COMMAND_PRIORITY_LOW,
  COMMAND_PRIORITY_HIGH,
  SELECTION_CHANGE_COMMAND,
  COMMAND_PRIORITY_CRITICAL,
} from 'lexical';

import { Input } from '@ui/form/Input';
import { Check } from '@ui/media/icons/Check';
import { Edit01 } from '@ui/media/icons/Edit01';
import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Bold01 } from '@ui/media/icons/Bold01.tsx';
import { FloatingToolbarButton } from '@ui/form/Editor/components';

import { sanitizeUrl } from '../utils/url';
import { getSelectedNode } from '../utils/getSelectedNode';
import { setFloatingElemPositionForLinkEditor } from '../utils/setFloatingElemPositionForLinkEditor';

function FloatingLinkEditor({
  editor,
  isLink,
  setIsLink,
  anchorElem,
  isLinkEditMode,
  setIsLinkEditMode,
}: {
  isLink: boolean;
  editor: LexicalEditor;
  anchorElem: HTMLElement;
  isLinkEditMode: boolean;
  setIsLink: Dispatch<boolean>;
  setIsLinkEditMode: Dispatch<boolean>;
}): JSX.Element {
  const editorRef = useRef<HTMLDivElement | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [linkUrl, setLinkUrl] = useState('');
  const [editedLinkUrl, setEditedLinkUrl] = useState('https://');
  const [lastSelection, setLastSelection] = useState<BaseSelection | null>(
    null,
  );

  const $updateLinkEditor = useCallback(() => {
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

      if (isLinkEditMode) {
        setEditedLinkUrl(linkUrl);
      }
    }
    const editorElem = editorRef.current;
    const nativeSelection = window.getSelection();
    const activeElement = document.activeElement;

    if (editorElem === null) {
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
      const domRect: DOMRect | undefined =
        nativeSelection.focusNode?.parentElement?.getBoundingClientRect();

      if (domRect) {
        domRect.y += 40;
        setFloatingElemPositionForLinkEditor(domRect, editorElem, anchorElem);
      }
      setLastSelection(selection);
    } else if (!activeElement || activeElement.className !== 'link-input') {
      if (rootElement !== null) {
        setFloatingElemPositionForLinkEditor(null, editorElem, anchorElem);
      }
      setLastSelection(null);
      setIsLinkEditMode(false);
      setLinkUrl('');
    }

    return true;
  }, [
    anchorElem,
    editor,
    editorRef,
    setIsLinkEditMode,
    isLink,
    isLinkEditMode,
    linkUrl,
  ]);

  useEffect(() => {
    const scrollerElem = anchorElem.parentElement;

    const update = () => {
      editor.getEditorState().read(() => {
        $updateLinkEditor();
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
  }, [anchorElem.parentElement, editor, $updateLinkEditor]);

  useEffect(() => {
    return mergeRegister(
      editor.registerUpdateListener(({ editorState }) => {
        editorState.read(() => {
          $updateLinkEditor();
        });
      }),

      editor.registerCommand(
        SELECTION_CHANGE_COMMAND,
        () => {
          $updateLinkEditor();

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
  }, [editor, $updateLinkEditor, setIsLink, isLink]);

  useEffect(() => {
    editor.getEditorState().read(() => {
      $updateLinkEditor();
    });
  }, [editor, $updateLinkEditor]);

  useEffect(() => {
    if (isLinkEditMode && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isLinkEditMode, isLink]);

  const monitorInputInteraction = (
    event: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (event.key === 'Enter') {
      event.preventDefault();
      handleLinkSubmission();
    } else if (event.key === 'Escape') {
      event.preventDefault();
      setIsLinkEditMode(false);
    }
  };

  const handleLinkSubmission = () => {
    if (lastSelection !== null) {
      if (editedLinkUrl !== '') {
        editor.dispatchCommand(TOGGLE_LINK_COMMAND, sanitizeUrl(editedLinkUrl));
        editor.update(() => {
          const selection = $getSelection();

          if ($isRangeSelection(selection)) {
            const parent = getSelectedNode(selection).getParent();

            if ($isAutoLinkNode(parent)) {
              const linkNode = $createLinkNode(parent.getURL(), {
                rel: parent.__rel,
                target: parent.__target,
                title: parent.__title,
              });

              parent.replace(linkNode, true);
            }
          }
        });
      }
      setEditedLinkUrl('https://');
      setIsLinkEditMode(false);
    }
  };

  return createPortal(
    <div
      ref={editorRef}
      style={{ pointerEvents: 'auto' }}
      className='absolute top-0 left-0 z-[99999] pointer-events-auto'
    >
      {(isLink || isLinkEditMode) && (
        <div
          data-side='bottom'
          onMouseDown={(e) => {
            e.preventDefault();
            e.preventDefault();
          }}
          className='bg-gray-700 flex items-center gap-2 min-w-[400px] max-w-[800px] py-1.5 px-[6px] shadow-lg  rounded-md data-[side=top]:animate-slideDownAndFade data-[side=right]:animate-slideLeftAndFade data-[side=bottom]:animate-slideUpAndFade data-[side=left]:animate-slideRightAndFade'
        >
          {isLinkEditMode ? (
            <>
              <Input
                ref={inputRef}
                variant='unstyled'
                value={editedLinkUrl}
                className='leading-none min-h-0 pointer-events-auto text-gray-25'
                onKeyDown={(event) => {
                  monitorInputInteraction(event);
                }}
                onChange={(event) => {
                  setEditedLinkUrl(event.target.value);
                }}
              />
              <FloatingToolbarButton
                label='Cancel'
                icon={<XClose className='text-gray-500' />}
                onMouseDown={(event) => event.preventDefault()}
                onClick={() => {
                  setIsLinkEditMode(false);
                }}
              />

              <FloatingToolbarButton
                label='Confirm'
                onClick={handleLinkSubmission}
                icon={<Check className='text-gray-500' />}
                onMouseDown={(event) => event.preventDefault()}
              />
            </>
          ) : (
            <>
              <a
                target='_blank'
                rel='noopener noreferrer'
                href={sanitizeUrl(linkUrl)}
              >
                {linkUrl}
              </a>
              <IconButton
                size='xs'
                tabIndex={0}
                variant='ghost'
                aria-label='edit'
                icon={<Edit01 />}
                className='pointer-events-auto'
                onMouseDown={(event) => event.preventDefault()}
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  setEditedLinkUrl(linkUrl);
                  setIsLink(true);
                  setIsLinkEditMode(true);
                }}
              />
              <IconButton
                size='xs'
                tabIndex={0}
                variant='ghost'
                icon={<Trash01 />}
                aria-label='delete'
                className='pointer-events-auto'
                onMouseDown={(event) => event.preventDefault()}
                onClick={() => {
                  editor.dispatchCommand(TOGGLE_LINK_COMMAND, null);
                }}
              />
            </>
          )}
        </div>
      )}
    </div>,
    document.body,
  );
}

function useFloatingLinkEditorToolbar(
  editor: LexicalEditor,
  anchorElem: HTMLElement,
  isLinkEditMode: boolean,
  setIsLinkEditMode: Dispatch<boolean>,
): JSX.Element | null {
  const [activeEditor, setActiveEditor] = useState(editor);
  const [isLink, setIsLink] = useState(false);
  const [isUserTyping, setIsUserTyping] = useState(false);

  useEffect(() => {
    function $updateToolbar() {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        const focusNode = getSelectedNode(selection);
        const focusLinkNode = $findMatchingParent(focusNode, $isLinkNode);
        const focusAutoLinkNode = $findMatchingParent(
          focusNode,
          $isAutoLinkNode,
        );

        if (!(focusLinkNode || focusAutoLinkNode)) {
          setIsLink(false);

          return;
        }

        const badNode = selection
          .getNodes()
          .filter((node) => !$isLineBreakNode(node))
          .find((node) => {
            const linkNode = $findMatchingParent(node, $isLinkNode);
            const autoLinkNode = $findMatchingParent(node, $isAutoLinkNode);

            return (
              (focusLinkNode && !focusLinkNode.is(linkNode)) ||
              (linkNode && !linkNode.is(focusLinkNode)) ||
              (focusAutoLinkNode && !focusAutoLinkNode.is(autoLinkNode)) ||
              (autoLinkNode && !autoLinkNode.is(focusAutoLinkNode))
            );
          });

        if (!badNode && isUserTyping) {
          setIsLink(true);
        } else {
          setIsLink(false);
        }
      }
    }

    function onKeyDown(event: KeyboardEvent) {
      if (
        event.key !== 'Meta' &&
        event.key !== 'Shift' &&
        event.key !== 'Alt' &&
        event.key !== 'Control'
      ) {
        setIsUserTyping(true);
      }
    }

    function onKeyUp(event: KeyboardEvent) {
      if (
        event.key !== 'Meta' &&
        event.key !== 'Shift' &&
        event.key !== 'Alt' &&
        event.key !== 'Control'
      ) {
        // Reset the typing state after a short delay
        setTimeout(() => setIsUserTyping(false), 500);
      }
    }

    document.addEventListener('keydown', onKeyDown);
    document.addEventListener('keyup', onKeyUp);

    return mergeRegister(
      editor.registerUpdateListener(({ editorState }) => {
        editorState.read(() => {
          $updateToolbar();
        });
      }),
      editor.registerCommand(
        SELECTION_CHANGE_COMMAND,
        (_payload, newEditor) => {
          $updateToolbar();
          setActiveEditor(newEditor);

          return false;
        },
        COMMAND_PRIORITY_CRITICAL,
      ),
      editor.registerCommand(
        CLICK_COMMAND,
        (payload) => {
          const selection = $getSelection();

          if ($isRangeSelection(selection)) {
            const node = getSelectedNode(selection);
            const linkNode = $findMatchingParent(node, $isLinkNode);

            if ($isLinkNode(linkNode)) {
              setIsLink(true);

              if (payload.metaKey || payload.ctrlKey) {
                window.open(linkNode.getURL(), '_blank');

                return true;
              }
            }
          }

          return false;
        },
        COMMAND_PRIORITY_LOW,
      ),
      () => {
        document.removeEventListener('keydown', onKeyDown);
        document.removeEventListener('keyup', onKeyUp);
      },
    );
  }, [editor, isUserTyping]);

  return createPortal(
    <FloatingLinkEditor
      isLink={isLink}
      editor={activeEditor}
      setIsLink={setIsLink}
      anchorElem={anchorElem}
      isLinkEditMode={isLinkEditMode}
      setIsLinkEditMode={setIsLinkEditMode}
    />,
    anchorElem,
  );
}

export default function FloatingLinkEditorPlugin({
  anchorElem = document.body,
  isLinkEditMode,
  setIsLinkEditMode,
}: {
  isLinkEditMode: boolean;
  anchorElem?: HTMLElement;
  setIsLinkEditMode: Dispatch<boolean>;
}): JSX.Element | null {
  const [editor] = useLexicalComposerContext();

  return useFloatingLinkEditorToolbar(
    editor,
    anchorElem,
    isLinkEditMode,
    setIsLinkEditMode,
  );
}

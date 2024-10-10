import * as React from 'react';
import { createPortal } from 'react-dom';
import { useRef, useState, useEffect, useCallback, KeyboardEvent } from 'react';

import { computePosition } from '@floating-ui/dom';
import { $findMatchingParent } from '@lexical/utils';
import { $isLinkNode, $toggleLink, TOGGLE_LINK_COMMAND } from '@lexical/link';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $getSelection,
  $setSelection,
  $isRangeSelection,
  KEY_ESCAPE_COMMAND,
  FORMAT_TEXT_COMMAND,
  KEY_MODIFIER_COMMAND,
  COMMAND_PRIORITY_HIGH,
  $createRangeSelection,
  COMMAND_PRIORITY_NORMAL,
  COMMAND_PRIORITY_NORMAL as NORMAL_PRIORITY,
  SELECTION_CHANGE_COMMAND as ON_SELECTION_CHANGE,
} from 'lexical';

import { Input } from '@ui/form/Input';
import { Bold01 } from '@ui/media/icons/Bold01';
import { Link01 } from '@ui/media/icons/Link01';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Divider } from '@ui/presentation/Divider';
import { Italic01 } from '@ui/media/icons/Italic01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { BlockQuote } from '@ui/media/icons/BlockQuote';
import { getExternalUrl } from '@utils/getExternalLink.ts';
import { sanitizeUrl } from '@ui/form/Editor/utils/url.ts';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { Strikethrough01 } from '@ui/media/icons/Strikethrough01';
import { FloatingToolbarButton } from '@ui/form/Editor/components';
import { $isExtendedQuoteNode } from '@ui/form/Editor/nodes/ExtendedQuoteNode';

import { usePointerInteractions } from './../utils/usePointerInteractions';
import {
  registerEnterQuoteCommand,
  TOGGLE_BLOCKQUOTE_COMMAND,
  registerToggleQuoteCommand,
} from './../commands';

const DEFAULT_DOM_ELEMENT = document.body;

type FloatingMenuCoords = { x: number; y: number } | undefined;

export type FloatingMenuComponentProps = {
  shouldShow: boolean;
  editor: ReturnType<typeof useLexicalComposerContext>[0];
};

export type FloatingMenuPluginProps = {
  element?: HTMLElement;
};

export function FloatingMenu({
  editor,
  shouldShow,
}: FloatingMenuComponentProps) {
  const [isStrikethrough, setIsStrikethrough] = useState(false);
  const [isLink, setIsLink] = useState(false);
  const [isBlockquote, setIsBlockquote] = useState(false);
  const [isBold, setIsBold] = useState(false);
  const [isItalic, setIsItalic] = useState(false);
  const [showLinkInput, setShowLinkInput] = useState(false);
  const [linkUrl, setLinkUrl] = useState('');
  const linkInputRef = useRef<HTMLInputElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const toggleQuoteCommand = registerToggleQuoteCommand(editor);
    const blockquoteEnterCommand = registerEnterQuoteCommand(editor);

    return () => {
      toggleQuoteCommand();
      blockquoteEnterCommand();
    };
  }, [editor]);

  const toggleLink = useCallback(() => {
    if (!isLink) {
      setShowLinkInput(true);
      setTimeout(() => linkInputRef?.current?.focus(), 0);
    } else {
      editor.update(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          $toggleLink(null);
        }
      });
    }
  }, [editor, isLink]);

  const handleLinkSubmit = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        $toggleLink(linkUrl);
      }
    });
    setShowLinkInput(false);
    setLinkUrl('');
  }, [editor, linkUrl]);

  const handleLinkCancel = useCallback(() => {
    setShowLinkInput(false);
    setLinkUrl('');
  }, []);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter') {
        e.preventDefault();
        handleLinkSubmit();
      } else if (e.key === 'Escape') {
        e.preventDefault();
        handleLinkCancel();
      }
    },
    [handleLinkSubmit, handleLinkCancel],
  );

  useEffect(() => {
    if (!shouldShow) {
      setShowLinkInput(false);
    }
  }, [shouldShow]);
  useEffect(() => {
    return editor.registerUpdateListener(({ editorState }) => {
      editorState.read(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          setIsStrikethrough(selection.hasFormat('strikethrough'));
          setIsBold(selection.hasFormat('bold'));
          setIsItalic(selection.hasFormat('italic'));
          setIsLink(
            selection
              .getNodes()
              .some(
                (node) =>
                  $isLinkNode(node) ||
                  ($isLinkNode(node.getParent()) &&
                    node?.getParent()?.isInline()),
              ),
          );
          setIsBlockquote(
            selection
              .getNodes()
              .some(
                (node) =>
                  $isExtendedQuoteNode(node) ||
                  $isExtendedQuoteNode(node.getParent()),
              ),
          );
        }
      });
    });
  }, [editor]);

  return (
    <div
      ref={menuRef}
      className='flex items-center justify-between bg-gray-700 text-gray-25 border-[1px] border-gray-200 rounded-md p-1 gap-1'
    >
      {showLinkInput ? (
        <>
          <Input
            type='url'
            value={linkUrl}
            ref={linkInputRef}
            variant='unstyled'
            placeholder='Enter URL'
            onKeyDown={handleKeyDown}
            onChange={(e) => setLinkUrl(e.target.value)}
            className='border-none rounded px-2 py-0 min-h-[auto] text-sm text-gray-25'
          />
          <Divider className='w-[1px] h-3 border-b-0 border-l-[1px] border-gray-500 mx-2' />

          <FloatingToolbarButton
            aria-label='Open link'
            icon={<LinkExternal02 className='text-inherit' />}
            onClick={() => {
              const link = getExternalUrl(sanitizeUrl(linkUrl));

              window.open(link, '_blank', 'noopener,noreferrer');
            }}
          />
          <FloatingToolbarButton
            aria-label='Cancel'
            onClick={handleLinkCancel}
            icon={<Trash01 className='text-inherit' />}
          />
        </>
      ) : (
        <>
          <Tooltip label='Bold: ⌘ + B'>
            <div>
              <FloatingToolbarButton
                active={isBold}
                aria-label='Format text to bold'
                icon={<Bold01 className='text-inherit' />}
                onClick={() => {
                  editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'bold');
                }}
              />
            </div>
          </Tooltip>
          <Tooltip label='Italic: ⌘ + I'>
            <div>
              <FloatingToolbarButton
                active={isItalic}
                aria-label='Format text with italic'
                icon={<Italic01 className='text-inherit' />}
                onClick={() => {
                  editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'italic');
                }}
              />
            </div>
          </Tooltip>
          <Tooltip label='Strikethrough: ⌘ + Shift + S'>
            <div>
              <FloatingToolbarButton
                active={isStrikethrough}
                aria-label='Format text with a strikethrough'
                icon={<Strikethrough01 className='text-inherit' />}
                onClick={() => {
                  editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'strikethrough');
                }}
              />
            </div>
          </Tooltip>
          <Tooltip label='Insert or remove link'>
            <FloatingToolbarButton
              active={isLink}
              onClick={toggleLink}
              aria-label='Insert or remove link'
              icon={<Link01 className='text-inherit' />}
            />
          </Tooltip>
          <Tooltip label='Blockquote: ⌘ + Shift + >'>
            <div>
              <FloatingToolbarButton
                active={isBlockquote}
                aria-label='Format text with block quote'
                icon={<BlockQuote className='text-inherit' />}
                onClick={() => {
                  editor.dispatchCommand(TOGGLE_BLOCKQUOTE_COMMAND, undefined);
                }}
              />
            </div>
          </Tooltip>
        </>
      )}
    </div>
  );
}

export function FloatingMenuPlugin({ element }: FloatingMenuPluginProps) {
  const ref = useRef<HTMLDivElement>(null);
  const [coords, setCoords] = useState<FloatingMenuCoords>(undefined);
  const show = coords !== undefined;

  const [editor] = useLexicalComposerContext();
  const { isPointerDown, isPointerReleased } = usePointerInteractions();

  const calculatePosition = useCallback(() => {
    const domSelection = getSelection();
    const domRange =
      domSelection?.rangeCount !== 0 && domSelection?.getRangeAt(0);

    if (!domRange || !ref.current || isPointerDown) return setCoords(undefined);

    computePosition(domRange, ref.current, { placement: 'top' })
      .then((pos) => {
        setCoords({ x: pos.x, y: pos.y - 10 });
      })
      .catch(() => {
        setCoords(undefined);
      });
  }, [isPointerDown]);

  const $handleSelectionChange = useCallback(() => {
    if (editor.isComposing()) return false;

    if (editor.getRootElement() !== document.activeElement) {
      setCoords(undefined);

      return true;
    }

    const selection = $getSelection();

    if ($isRangeSelection(selection) && !selection.anchor.is(selection.focus)) {
      const node = selection.getNodes()[0];
      const linkNode = $findMatchingParent(node, $isLinkNode);

      if ($isLinkNode(linkNode)) {
        return false;
      }

      calculatePosition();
    } else {
      setCoords(undefined);
    }

    return true;
  }, [editor, calculatePosition]);

  const handleClickOutside = useCallback(
    (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        editor.update(() => {
          setCoords(undefined);

          const currentSelection = $getSelection();

          if ($isRangeSelection(currentSelection)) {
            const newSelection = $createRangeSelection();

            newSelection.anchor.set(
              currentSelection.focus.key,
              currentSelection.focus.offset,
              currentSelection.focus.type,
            );
            newSelection.focus.set(
              currentSelection.focus.key,
              currentSelection.focus.offset,
              currentSelection.focus.type,
            );
            newSelection.dirty = true;
            $setSelection(newSelection);
          }
        });
      }
    },
    [editor],
  );

  useEffect(() => {
    const unregisterCommand = editor.registerCommand(
      ON_SELECTION_CHANGE,
      $handleSelectionChange,
      NORMAL_PRIORITY,
    );

    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      unregisterCommand();
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [editor, $handleSelectionChange, handleClickOutside]);

  useEffect(() => {
    const unregisterCommand = editor.registerCommand(
      KEY_ESCAPE_COMMAND,
      () => {
        setCoords(undefined);

        return false;
      },
      COMMAND_PRIORITY_HIGH,
    );

    return unregisterCommand;
  }, [editor]);

  useEffect(() => {
    const removeKeyboardHandler = editor.registerCommand(
      KEY_MODIFIER_COMMAND,
      (event: KeyboardEvent) => {
        if (event.key === 's' && (event.metaKey || event.ctrlKey)) {
          event.preventDefault();
          editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'strikethrough');

          return true;
        }

        if (
          event.key === '.' &&
          event.shiftKey &&
          (event.metaKey || event.ctrlKey)
        ) {
          event.preventDefault();
          editor.dispatchCommand(TOGGLE_BLOCKQUOTE_COMMAND, undefined);

          return true;
        }

        if (event.key === 'k' && (event.metaKey || event.ctrlKey)) {
          editor.dispatchCommand(TOGGLE_LINK_COMMAND, 'https://example.com');
          calculatePosition();

          return true;
        }

        return false;
      },
      COMMAND_PRIORITY_NORMAL,
    );

    return () => {
      removeKeyboardHandler();
    };
  }, [editor]);

  useEffect(() => {
    if (!show && isPointerReleased) {
      editor.getEditorState().read(() => {
        $handleSelectionChange();
      });
    }
  }, [isPointerReleased, $handleSelectionChange, editor]);

  return createPortal(
    <div
      ref={ref}
      aria-hidden={!show}
      style={{
        position: 'absolute',
        top: coords?.y,
        left: coords?.x,
        visibility: show ? 'visible' : 'hidden',
        opacity: show ? 1 : 0,
      }}
    >
      <FloatingMenu editor={editor} shouldShow={show} />
    </div>,
    element ?? DEFAULT_DOM_ELEMENT,
  );
}

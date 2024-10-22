import * as React from 'react';
import { createPortal } from 'react-dom';
import { useRef, useState, useEffect, useCallback, KeyboardEvent } from 'react';

import { computePosition } from '@floating-ui/dom';
import { $isLinkNode, $toggleLink } from '@lexical/link';
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

import { Bold01 } from '@ui/media/icons/Bold01';
import { Link01 } from '@ui/media/icons/Link01';
import { Italic01 } from '@ui/media/icons/Italic01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { BlockQuote } from '@ui/media/icons/BlockQuote';
import { BracketsPlus } from '@ui/media/icons/BracketsPlus';
import { Strikethrough01 } from '@ui/media/icons/Strikethrough01';
import { FloatingToolbarButton } from '@ui/form/Editor/components';
import { getSelectedNode } from '@ui/form/Editor/utils/getSelectedNode.ts';
import { $isExtendedQuoteNode } from '@ui/form/Editor/nodes/ExtendedQuoteNode';

import { usePointerInteractions } from './../utils/usePointerInteractions';
import {
  INSERT_VARIABLE_NODE,
  registerEnterQuoteCommand,
  TOGGLE_BLOCKQUOTE_COMMAND,
  registerToggleQuoteCommand,
  registerInsertVariableNodeCommand,
} from './../commands';

const DEFAULT_DOM_ELEMENT = document.body;

type FloatingMenuCoords = { x: number; y: number } | undefined;

export type FloatingMenuComponentProps = {
  variableOptions: string[];
  editor: ReturnType<typeof useLexicalComposerContext>[0];
};

export type FloatingMenuPluginProps = {
  element?: HTMLElement;
  variableOptions: string[];
};

export function FloatingMenu({
  editor,
  variableOptions,
}: FloatingMenuComponentProps) {
  const [isStrikethrough, setIsStrikethrough] = useState(false);
  const [isLink, setIsLink] = useState(false);
  const [isBlockquote, setIsBlockquote] = useState(false);
  const [isBold, setIsBold] = useState(false);
  const [isItalic, setIsItalic] = useState(false);
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
    editor.update(() => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        if (isLink) {
          $toggleLink(null);
        } else {
          $toggleLink('https://');
        }
      }
    });
  }, [editor, isLink]);

  useEffect(() => {
    return editor.registerUpdateListener(({ editorState }) => {
      editorState.read(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          setIsStrikethrough(selection.hasFormat('strikethrough'));
          setIsBold(selection.hasFormat('bold'));
          setIsItalic(selection.hasFormat('italic'));

          const node = getSelectedNode(selection);

          // Update links
          const parent = node.getParent();

          if ($isLinkNode(parent) || $isLinkNode(node)) {
            setIsLink(true);
          } else {
            setIsLink(false);
          }
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
        <Tooltip label='Insert or remove link: ⌘ + K'>
          <div>
            <FloatingToolbarButton
              active={isLink}
              onClick={toggleLink}
              aria-label='Insert or remove link'
              icon={<Link01 className='text-inherit' />}
            />
          </div>
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
        <Tooltip label='Add variable: {'>
          <div>
            <FloatingToolbarButton
              active={isBlockquote}
              aria-label='Add variable'
              icon={<BracketsPlus className='text-inherit' />}
              onClick={() => {
                editor.dispatchCommand(INSERT_VARIABLE_NODE, {
                  label: variableOptions?.[0].toLowerCase(),
                  value: variableOptions?.[0].toLowerCase(),
                });
              }}
            />
          </div>
        </Tooltip>
      </>
    </div>
  );
}

export function FloatingMenuPlugin({
  element,
  variableOptions,
}: FloatingMenuPluginProps) {
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
    if (editor.isComposing()) {
      return false;
    }

    if (editor.getRootElement() !== document.activeElement) {
      setCoords(undefined);

      return true;
    }

    const selection = $getSelection();

    if ($isRangeSelection(selection) && !selection.anchor.is(selection.focus)) {
      const node = getSelectedNode(selection);

      if ($isLinkNode(node)) {
        setCoords(undefined);

        return false;
      }
      calculatePosition();
    } else {
      setCoords(undefined);
    }

    return false;
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
    const unregisterInsertVariableCommand = registerInsertVariableNodeCommand(
      editor,
      () => {
        setCoords(undefined);
      },
    );
    const unregisterCommand = editor.registerCommand(
      ON_SELECTION_CHANGE,
      $handleSelectionChange,
      NORMAL_PRIORITY,
    );

    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      unregisterCommand();
      unregisterInsertVariableCommand();
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
          editor.update(() => {
            const selection = $getSelection();

            if ($isRangeSelection(selection)) {
              $toggleLink('');
            }
          });

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
      <FloatingMenu editor={editor} variableOptions={variableOptions} />
    </div>,
    element ?? DEFAULT_DOM_ELEMENT,
  );
}

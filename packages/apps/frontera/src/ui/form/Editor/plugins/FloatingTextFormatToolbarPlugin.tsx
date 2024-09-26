import { createPortal } from 'react-dom';
import { useRef, useState, useEffect, useCallback } from 'react';

import { computePosition } from '@floating-ui/dom';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $getSelection,
  $isRangeSelection,
  FORMAT_TEXT_COMMAND,
  COMMAND_PRIORITY_NORMAL as NORMAL_PRIORITY,
  SELECTION_CHANGE_COMMAND as ON_SELECTION_CHANGE,
} from 'lexical';

import { cn } from '@ui/utils/cn.ts';
import { IconButton } from '@ui/form/IconButton';
import { Code01 } from '@ui/media/icons/Code01.tsx';
import { Bold01 } from '@ui/media/icons/Bold01.tsx';
import { Italic01 } from '@ui/media/icons/Italic01.tsx';
import { Underline01 } from '@ui/media/icons/Underline01.tsx';
import { Strikethrough01 } from '@ui/media/icons/Strikethrough01.tsx';

import { usePointerInteractions } from './../utils/usePointerInteractions.tsx';

const DEFAULT_DOM_ELEMENT = document.body;

type FloatingMenuCoords = { x: number; y: number } | undefined;

export type FloatingMenuComponentProps = {
  shouldShow: boolean;
  editor: ReturnType<typeof useLexicalComposerContext>[0];
};

export type FloatingMenuPluginProps = {
  element?: HTMLElement;
};

export function FloatingMenu({ editor }: FloatingMenuComponentProps) {
  const [state, setState] = useState<FloatingMenuState>({
    isBold: false,
    isCode: false,
    isItalic: false,
    isStrikethrough: false,
    isUnderline: false,
  });

  useEffect(() => {
    const unregisterListener = editor.registerUpdateListener(
      ({ editorState }) => {
        editorState.read(() => {
          const selection = $getSelection();

          if (!$isRangeSelection(selection)) return;
          setState({
            isBold: selection.hasFormat('bold'),
            isCode: selection.hasFormat('code'),
            isItalic: selection.hasFormat('italic'),
            isStrikethrough: selection.hasFormat('strikethrough'),
            isUnderline: selection.hasFormat('underline'),
          });
        });
      },
    );

    return unregisterListener;
  }, [editor]);

  return (
    <div className='flex items-center justify-between bg-white border-[1px] border-gray-200 rounded-md '>
      <IconButton
        variant='ghost'
        icon={<Bold01 />}
        aria-label='Format text as bold'
        className={cn('rounded-r-none', {
          'bg-gray-100': state.isBold,
        })}
        onClick={() => {
          editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'bold');
        }}
      />
      <IconButton
        variant='ghost'
        icon={<Italic01 />}
        aria-label='Format text as italics'
        className={cn('rounded-none', {
          'bg-gray-100': state.isItalic,
        })}
        onClick={() => {
          editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'italic');
        }}
      />
      <IconButton
        variant='ghost'
        icon={<Underline01 />}
        aria-label='Format text to underlined'
        className={cn('rounded-none', {
          'bg-gray-100': state.isUnderline,
        })}
        onClick={() => {
          editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'underline');
        }}
      />
      <IconButton
        variant='ghost'
        icon={<Strikethrough01 />}
        aria-label='Format text with a strikethrough'
        className={cn('rounded-none', {
          'bg-gray-100': state.isStrikethrough,
        })}
        onClick={() => {
          editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'strikethrough');
        }}
      />
      <IconButton
        variant='ghost'
        icon={<Code01 />}
        aria-label='Format text with inline code'
        className={cn('rounded-l-none', {
          'bg-gray-100': state.isCode,
        })}
        onClick={() => {
          editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'code');
        }}
      />
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
      calculatePosition();
    } else {
      setCoords(undefined);
    }

    return true;
  }, [editor, calculatePosition]);

  useEffect(() => {
    const unregisterCommand = editor.registerCommand(
      ON_SELECTION_CHANGE,
      $handleSelectionChange,
      NORMAL_PRIORITY,
    );

    return unregisterCommand;
  }, [editor, $handleSelectionChange]);

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

export type FloatingMenuState = {
  isBold: boolean;
  isCode: boolean;
  isItalic: boolean;
  isUnderline: boolean;
  isStrikethrough: boolean;
};

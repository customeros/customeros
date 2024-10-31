import * as React from 'react';
import { useState, useEffect, useCallback } from 'react';

import { $setBlocksType } from '@lexical/selection';
import { $isQuoteNode, $createQuoteNode } from '@lexical/rich-text';
import { $isListNode, $createListNode, $isListItemNode } from '@lexical/list';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  LexicalNode,
  $getSelection,
  $isRangeSelection,
  FORMAT_TEXT_COMMAND,
  $createParagraphNode,
} from 'lexical';

import { cn } from '@ui/utils/cn.ts';
import { Bold01 } from '@ui/media/icons/Bold01';
import { IconButton } from '@ui/form/IconButton';
import { Italic01 } from '@ui/media/icons/Italic01';
import { BlockQuote } from '@ui/media/icons/BlockQuote';
import { ListBulleted } from '@ui/media/icons/ListBulleted';
import { ListNumbered } from '@ui/media/icons/ListNumbered';
import { Strikethrough01 } from '@ui/media/icons/Strikethrough01';

const activeStyle = 'bg-gray-100 text-gray-700 hover:bg-gray-100';

export default function ToolbarPlugin(): JSX.Element {
  const [editor] = useLexicalComposerContext();
  const [isStrikethrough, setIsStrikethrough] = useState(false);
  const [isBlockquote, setIsBlockquote] = useState(false);
  const [isBold, setIsBold] = useState(false);
  const [isItalic, setIsItalic] = useState(false);
  const [isOrderedList, setIsOrderedList] = useState(false);
  const [isUnorderedList, setIsUnorderedList] = useState(false);

  const toggleBlockquote = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if (!isBlockquote) {
        $setBlocksType(selection, $createQuoteNode);
      } else {
        $setBlocksType(selection, $createParagraphNode);
      }
    });
  }, [editor, isBlockquote]);

  const toggleOrderedList = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if (!isOrderedList) {
        $setBlocksType(selection, () => $createListNode('number'));
      } else {
        $setBlocksType(selection, $createParagraphNode);
      }
    });
  }, [editor, isOrderedList]);

  const toggleUnorderedList = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if (!isUnorderedList) {
        $setBlocksType(selection, () => $createListNode('bullet'));
      } else {
        $setBlocksType(selection, $createParagraphNode);
      }
    });
  }, [editor, isUnorderedList]);

  useEffect(() => {
    return editor.registerUpdateListener(({ editorState }) => {
      editorState.read(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          setIsStrikethrough(selection.hasFormat('strikethrough'));
          setIsBold(selection.hasFormat('bold'));
          setIsItalic(selection.hasFormat('italic'));
          setIsBlockquote(
            selection
              .getNodes()
              .some((n) => $isQuoteNode(n) || $isQuoteNode(n.getParent())),
          );

          let isUnordered = false;
          let isOrdered = false;
          // const isCheck = false;

          selection.getNodes().forEach((node) => {
            let currentNode: LexicalNode | null = node;

            while (currentNode != null) {
              if ($isListItemNode(currentNode)) {
                const parent = currentNode.getParent();

                if ($isListNode(parent)) {
                  if (parent.getListType() === 'bullet') {
                    isUnordered = true;
                  } else if (parent.getListType() === 'number') {
                    isOrdered = true;
                  }
                }
              }

              // Move to parent node to continue checking
              currentNode = currentNode.getParent() as LexicalNode | null;
            }
          });

          setIsUnorderedList(isUnordered);
          setIsOrderedList(isOrdered);
        }
      });
    });
  }, [editor]);

  return (
    <div className='flex items-center'>
      <>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Format text to bold'
          icon={<Bold01 className='text-inherit' />}
          className={cn('rounded-sm', {
            [activeStyle]: isBold,
          })}
          onClick={() => {
            editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'bold');
          }}
        />

        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Format text with italic'
          icon={<Italic01 className='text-inherit' />}
          className={cn('rounded-sm', {
            [activeStyle]: isItalic,
          })}
          onClick={() => {
            editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'italic');
          }}
        />
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Format text with a strikethrough'
          icon={<Strikethrough01 className='text-inherit' />}
          className={cn('rounded-sm', {
            [activeStyle]: isStrikethrough,
          })}
          onClick={() => {
            editor.dispatchCommand(FORMAT_TEXT_COMMAND, 'strikethrough');
          }}
        />
      </>
      <div className='h-5 w-[1px] bg-gray-400 mx-1' />
      <>
        <IconButton
          size='xs'
          variant='ghost'
          onClick={toggleUnorderedList}
          aria-label='Format text as an bullet list'
          icon={<ListBulleted className='text-inherit' />}
          className={cn('rounded-sm', {
            [activeStyle]: isUnorderedList,
          })}
        />
        <IconButton
          size='xs'
          variant='ghost'
          onClick={toggleOrderedList}
          aria-label='Format text as an ordered list'
          icon={<ListNumbered className='text-inherit' />}
          className={cn('rounded-sm', {
            [activeStyle]: isOrderedList,
          })}
        />
      </>
      <div className='h-5 w-[1px] bg-gray-400 mx-0.5' />

      <IconButton
        size='xs'
        variant='ghost'
        onClick={toggleBlockquote}
        aria-label='Format text with block quote'
        icon={<BlockQuote className='text-inherit' />}
        className={cn('rounded-sm', {
          [activeStyle]: isBlockquote,
        })}
      />
    </div>
  );
}

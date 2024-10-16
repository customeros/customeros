import * as React from 'react';
import { createPortal } from 'react-dom';
import { useRef, useState, useEffect, useCallback } from 'react';

import { $findMatchingParent } from '@lexical/utils';
import { offset, computePosition } from '@floating-ui/dom';
import { $isLinkNode, $toggleLink, $createLinkNode } from '@lexical/link';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $getSelection,
  CLICK_COMMAND,
  $createTextNode,
  $isRangeSelection,
  COMMAND_PRIORITY_HIGH,
  COMMAND_PRIORITY_NORMAL,
  SELECTION_CHANGE_COMMAND,
} from 'lexical';

import { Input } from '@ui/form/Input';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Divider } from '@ui/presentation/Divider';
import { getExternalUrl } from '@utils/getExternalLink';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { FloatingToolbarButton } from '@ui/form/Editor/components';

import { sanitizeUrl } from '../utils/url';
import { getSelectedNode } from '../utils/getSelectedNode';
import { usePointerInteractions } from '../utils/usePointerInteractions';

const DEFAULT_DOM_ELEMENT = document.body;

type FloatingLinkEditorComponentProps = {
  isLink: boolean;
  editor: ReturnType<typeof useLexicalComposerContext>[0];
  setIsLink: React.Dispatch<React.SetStateAction<boolean>>;
};

export function FloatingLinkEditor({
  editor,
  isLink,
  setIsLink,
}: FloatingLinkEditorComponentProps) {
  const [linkUrl, setLinkUrl] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    editor.getEditorState().read(() => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        const node = getSelectedNode(selection);
        const parent = node.getParent();

        if ($isLinkNode(parent)) {
          setLinkUrl(parent.getURL());
        } else if ($isLinkNode(node)) {
          setLinkUrl(node.getURL());
        }
      }
    });
  }, [editor, isLink]);

  useEffect(() => {
    if (isLink && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isLink]);

  const handleLinkSubmission = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        const node = getSelectedNode(selection);
        const parent = node.getParent();

        if (linkUrl.trim() === '') {
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
          let linkNode;

          if ($isLinkNode(parent)) {
            parent.setURL(sanitizeUrl(linkUrl));
            linkNode = parent;
          } else if ($isLinkNode(node)) {
            node.setURL(sanitizeUrl(linkUrl));
            linkNode = node;
          } else {
            linkNode = $createLinkNode(sanitizeUrl(linkUrl));
            selection.insertNodes([linkNode]);
          }

          const spaceNode = $createTextNode(' ');

          linkNode.insertAfter(spaceNode);
          spaceNode.select(0, 0);
        }

        setIsLink(false);
      }
    });
    setIsLink(false);
  }, [editor, linkUrl, setIsLink]);

  const handleDeleteLink = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        $toggleLink(null);
        setTimeout(() => {
          setIsLink(false);
        }, 0);
      }
    });
  }, [editor, setIsLink]);

  return (
    <div className='bg-gray-700 flex items-center min-w-[auto] max-w-[800px] p-1 pl-3 shadow-lg rounded-md'>
      <Input
        size='sm'
        ref={inputRef}
        value={linkUrl}
        variant='unstyled'
        placeholder='Enter a URL'
        onChange={(event) => setLinkUrl(event.target.value)}
        className='leading-none min-h-0 pointer-events-auto text-gray-25 overflow-ellipsis'
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
        icon={<LinkExternal02 className='text-inherit' />}
        onClick={() => {
          const link = getExternalUrl(sanitizeUrl(linkUrl));

          window.open(link, '_blank', 'noopener,noreferrer');
        }}
      />
      <FloatingToolbarButton
        aria-label='Delete link'
        onClick={handleDeleteLink}
        icon={<Trash01 className='text-inherit' />}
      />
    </div>
  );
}

export function FloatingLinkEditorPlugin({
  anchorElem = DEFAULT_DOM_ELEMENT,
}: {
  anchorElem?: HTMLElement;
}): JSX.Element | null {
  const [editor] = useLexicalComposerContext();
  const [isLink, setIsLink] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const [menuPosition, setMenuPosition] = useState<{
    top: number;
    left: number;
  } | null>(null);
  const anchorRef = useRef<HTMLElement | null>(null);
  const { isPointerDown, isPointerReleased } = usePointerInteractions();
  const closeTimeoutRef = useRef<number | null>(null);

  const updateMenuPosition = useCallback(() => {
    if (anchorRef.current && ref.current && !isPointerDown) {
      computePosition(anchorRef.current, ref.current, {
        placement: 'top-start',
        middleware: [
          offset({
            mainAxis: (anchorRef.current?.offsetHeight ?? 0) + 18,
            crossAxis: 0,
          }),
        ],
      }).then(({ x, y }) => {
        setMenuPosition({ top: y, left: x });
      });
    }
  }, [anchorRef, ref, isPointerDown]);

  const $handleSelectionChange = useCallback(() => {
    if (editor.isComposing()) return false;

    if (editor.getRootElement() !== document.activeElement) {
      setMenuPosition(null);

      return false;
    }

    const selection = $getSelection();

    if ($isRangeSelection(selection)) {
      const node = getSelectedNode(selection);
      const linkParent = $findMatchingParent(node, $isLinkNode);
      const linkNode = $isLinkNode(linkParent)
        ? linkParent
        : $isLinkNode(node)
        ? node
        : null;

      if (linkNode) {
        setIsLink(true);

        const element = editor.getElementByKey(
          linkNode.getKey(),
        ) as HTMLElement;

        if (element) {
          anchorRef.current = element;
          requestAnimationFrame(updateMenuPosition);
        }
      } else {
        setIsLink(false);
        anchorRef.current = null;
        setMenuPosition(null);
      }
    } else {
      setIsLink(false);
      anchorRef.current = null;
      setMenuPosition(null);
    }

    return false;
  }, [editor, updateMenuPosition]);

  useEffect(() => {
    return editor.registerCommand(
      SELECTION_CHANGE_COMMAND,
      $handleSelectionChange,
      COMMAND_PRIORITY_HIGH,
    );
  }, [editor, $handleSelectionChange]);

  useEffect(() => {
    if (!isLink && isPointerReleased) {
      editor.getEditorState().read(() => {
        $handleSelectionChange();
      });
    }
  }, [isPointerReleased, $handleSelectionChange, editor, isLink]);

  const closeMenu = useCallback(() => {
    if (closeTimeoutRef.current === null) {
      closeTimeoutRef.current = window.setTimeout(() => {
        setIsLink(false);
        setMenuPosition(null);
        closeTimeoutRef.current = null;
      }, 100);
    }
  }, []);

  useEffect(() => {
    return editor.registerCommand(
      CLICK_COMMAND,
      () => {
        if (isLink) {
          return true;
        }
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          const node = getSelectedNode(selection);
          const linkNode = $findMatchingParent(node, $isLinkNode);

          linkNode?.select(0, linkNode?.getTextContentSize());
        }

        return false;
      },
      COMMAND_PRIORITY_NORMAL,
    );
  }, [editor]);
  useEffect(() => {
    const rootElement = editor.getRootElement();

    if (!rootElement) {
      return;
    }

    const handleDoubleClick = (event: MouseEvent) => {
      event.preventDefault();

      editor.update(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          const node = getSelectedNode(selection);
          const linkNode = $findMatchingParent(node, $isLinkNode);

          if (linkNode) {
            linkNode.select(0, linkNode.getTextContentSize());
            setIsLink(true);
            updateMenuPosition();
          }
        }
      });
    };

    rootElement.addEventListener('dblclick', handleDoubleClick);

    return () => {
      rootElement.removeEventListener('dblclick', handleDoubleClick);
    };
  }, [editor, updateMenuPosition]);
  useEffect(() => {
    if (!isLink && isPointerReleased) {
      editor.getEditorState().read(() => {
        $handleSelectionChange();
      });
    }
  }, [isPointerReleased, $handleSelectionChange, editor, isLink]);
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        ref.current &&
        !ref.current.contains(event.target as Node) &&
        anchorRef.current &&
        !anchorRef.current.contains(event.target as Node)
      ) {
        closeMenu();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [closeMenu]);

  return (
    <>
      {isLink &&
        createPortal(
          <div
            ref={ref}
            tabIndex={-1}
            aria-hidden={!isLink}
            style={{
              position: 'absolute',
              top: menuPosition?.top ?? 0,
              left: menuPosition?.left ?? 0,
              visibility: isLink && menuPosition ? 'visible' : 'hidden',
              opacity: isLink && menuPosition ? 1 : 0,
              pointerEvents: 'all',
            }}
          >
            {isLink && menuPosition && (
              <FloatingLinkEditor
                editor={editor}
                isLink={isLink}
                setIsLink={setIsLink}
              />
            )}
          </div>,
          anchorElem,
        )}
    </>
  );
}
export default FloatingLinkEditorPlugin;

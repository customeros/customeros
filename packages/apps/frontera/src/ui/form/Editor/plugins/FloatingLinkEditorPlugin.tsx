import { createPortal } from 'react-dom';
import {
  useRef,
  useState,
  Dispatch,
  useEffect,
  useCallback,
  ReactPortal,
  KeyboardEvent,
  SetStateAction,
} from 'react';

import { offset, computePosition } from '@floating-ui/dom';
import { mergeRegister, $findMatchingParent } from '@lexical/utils';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  $isLinkNode,
  $toggleLink,
  $createLinkNode,
  $isAutoLinkNode,
} from '@lexical/link';
import {
  $getSelection,
  CLICK_COMMAND,
  $createTextNode,
  $isLineBreakNode,
  $isRangeSelection,
  COMMAND_PRIORITY_LOW,
  COMMAND_PRIORITY_HIGH,
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
  setIsLink: Dispatch<SetStateAction<boolean>>;
  editor: ReturnType<typeof useLexicalComposerContext>[0];
};

export function FloatingLinkEditor({
  editor,
  isLink,
  setIsLink,
}: FloatingLinkEditorComponentProps) {
  const [linkUrl, setLinkUrl] = useState('');
  const inputRef = useRef<HTMLInputElement | null>(null);
  const editorRef = useRef<HTMLDivElement | null>(null);

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

  const handleLinkSubmission = useCallback(() => {
    editor.update(() => {
      const selection = $getSelection();

      if ($isRangeSelection(selection)) {
        const node = getSelectedNode(selection);
        const parent = node.getParent();

        if (linkUrl.trim() === '') {
          if ($isLinkNode(parent)) {
            $toggleLink(null);
          } else if ($isLinkNode(node)) {
            $toggleLink(null);
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

          const spaceNode = $createTextNode('');

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
        setIsLink(false);
      }
    });
  }, [editor, setIsLink]);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        editorRef.current &&
        !editorRef.current?.contains(event.target as Node)
      ) {
        // todo - this condition could be more sofisticated
        if (linkUrl.trim().length || linkUrl !== 'https://') {
          handleLinkSubmission();
        } else {
          handleDeleteLink();
        }
      }
    }

    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isLink, handleLinkSubmission]);

  const monitorInputInteraction = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') {
      event.preventDefault();
      handleLinkSubmission();
      setIsLink(false);
    } else if (event.key === 'Escape') {
      event.preventDefault();
      setIsLink(false);
    }
  };

  return (
    <div
      ref={editorRef}
      className='bg-gray-700 flex items-center min-w-[240px] max-w-[240px] p-1 pl-3 shadow-lg rounded-md'
    >
      <Input
        size='sm'
        ref={inputRef}
        value={linkUrl}
        variant='unstyled'
        placeholder='Enter a URL'
        onMouseDown={(event) => event.stopPropagation()}
        onChange={(event) => setLinkUrl(event.target.value)}
        onKeyDown={(event) => {
          monitorInputInteraction(event);
        }}
        className='leading-none min-h-0 pointer-events-auto text-gray-25 overflow-ellipsis'
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          inputRef?.current?.focus();
        }}
      />

      <Divider className='w-[1px] h-3 border-b-0 border-l-[1px] border-gray-500 mx-2' />

      {!!linkUrl.trim().length && linkUrl.trim() !== 'https://' && (
        <FloatingToolbarButton
          aria-label='Open link'
          icon={<LinkExternal02 className='text-inherit' />}
          onMouseDown={(event) => {
            event.preventDefault();
            event.stopPropagation();
          }}
          onClick={() => {
            const link = getExternalUrl(sanitizeUrl(linkUrl));

            window.open(link, '_blank', 'noopener,noreferrer');
          }}
        />
      )}

      <FloatingToolbarButton
        aria-label='Delete link'
        onClick={handleDeleteLink}
        icon={<Trash01 className='text-inherit' />}
        onMouseDown={(event) => {
          event.preventDefault();
          event.stopPropagation();
        }}
      />
    </div>
  );
}

export function FloatingLinkEditorPlugin({
  anchorElem = DEFAULT_DOM_ELEMENT,
}: {
  anchorElem?: HTMLElement;
}): ReactPortal | null {
  const [editor] = useLexicalComposerContext();
  const [isLink, setIsLink] = useState(false);
  const ref = useRef<HTMLDivElement | null>(null);
  const [menuPosition, setMenuPosition] = useState<{
    top: number;
    left: number;
  } | null>(null);
  const anchorRef = useRef<HTMLElement | null>(null);
  const { isPointerDown, isPointerReleased } = usePointerInteractions();

  const updateMenuPosition = useCallback(() => {
    if (anchorRef.current && ref.current && !isPointerDown) {
      computePosition(anchorRef.current, ref.current, {
        placement: 'bottom-start',
        middleware: [offset(8)],
      }).then(({ x, y }) => {
        setMenuPosition({ top: y, left: x });
      });
    }
  }, [anchorRef, ref, isPointerDown]);

  const $handleSelectionChange = useCallback(() => {
    if (editor.isComposing()) return false;

    if (
      editor.getRootElement() !== document.activeElement ||
      !isPointerReleased
    ) {
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

    return true;
  }, [editor, updateMenuPosition]);

  useEffect(() => {
    return mergeRegister(
      editor.registerUpdateListener(({ editorState }) => {
        editorState.read(() => {
          $handleSelectionChange();
        });
      }),
      editor.registerCommand(
        SELECTION_CHANGE_COMMAND,
        $handleSelectionChange,
        COMMAND_PRIORITY_LOW,
      ),
    );
  }, [editor, $handleSelectionChange]);

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
              (focusAutoLinkNode && !focusAutoLinkNode.is(autoLinkNode))
            );
          });

        if (!badNode) {
          $handleSelectionChange();
        } else {
          setIsLink(false);
        }
      }
    }

    return mergeRegister(
      editor.registerUpdateListener(({ editorState }) => {
        editorState.read(() => {
          $updateToolbar();
        });
      }),

      editor.registerCommand(
        CLICK_COMMAND,
        (payload) => {
          const selection = $getSelection();

          if ($isRangeSelection(selection)) {
            const node = getSelectedNode(selection);
            const linkNode = $findMatchingParent(node, $isLinkNode);

            if ($isLinkNode(linkNode) && (payload.metaKey || payload.ctrlKey)) {
              window.open(linkNode.getURL(), '_blank');

              return true;
            }

            if ($isLinkNode(linkNode) && !menuPosition) {
              const element = editor.getElementByKey(
                linkNode.getKey(),
              ) as HTMLElement;

              if (element) {
                anchorRef.current = element;
                requestAnimationFrame(updateMenuPosition);
              }
              setIsLink(true);

              return true;
            }
          }

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),
    );
  }, [editor]);

  return createPortal(
    <div
      ref={ref}
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
      {isLink && (
        <FloatingLinkEditor
          isLink={isLink}
          editor={editor}
          setIsLink={setIsLink}
        />
      )}
    </div>,
    anchorElem,
  );
}
export default FloatingLinkEditorPlugin;

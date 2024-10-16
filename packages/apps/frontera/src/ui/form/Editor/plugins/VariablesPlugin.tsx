import { createPortal } from 'react-dom';
import { useRef, useMemo, useState, useEffect, useCallback } from 'react';

import { mergeRegister } from '@lexical/utils';
import { flip, shift, offset, computePosition } from '@floating-ui/dom';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  TextNode,
  $getSelection,
  $setSelection,
  $createTextNode,
  KEY_TAB_COMMAND,
  $isRangeSelection,
  KEY_ENTER_COMMAND,
  KEY_SPACE_COMMAND,
  KEY_ESCAPE_COMMAND,
  KEY_ARROW_UP_COMMAND,
  $createRangeSelection,
  COMMAND_PRIORITY_HIGH,
  KEY_BACKSPACE_COMMAND,
  KEY_ARROW_DOWN_COMMAND,
  SELECTION_CHANGE_COMMAND,
} from 'lexical';

import { cn } from '@ui/utils/cn';
import { SelectOption } from '@ui/utils/types';
import { User01 } from '@ui/media/icons/User01.tsx';
import { Building07 } from '@ui/media/icons/Building07.tsx';
import {
  VariableNode,
  $isVariableNode,
  $createVariableNode,
} from '@ui/form/Editor/nodes/VariableNode';

interface VariablesPluginProps {
  options: string[];
}

export default function VariablesPlugin({
  options = [],
}: VariablesPluginProps): JSX.Element | null {
  const [editor] = useLexicalComposerContext();
  const [menuItems, setMenuItems] = useState<SelectOption[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [showSuggestion, setShowSuggestion] = useState(false);
  const [menuPosition, setMenuPosition] = useState<{
    top: number;
    left: number;
  } | null>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const anchorRef = useRef<HTMLElement | null>(null);
  const closeTimeoutRef = useRef<number | null>(null);

  const selectableOptions = useMemo(
    () =>
      options.map((item) => ({
        label: item.toLowerCase(),
        value: item.toLowerCase(),
      })),
    [options],
  );

  const updateMenuPosition = useCallback(() => {
    if (anchorRef.current && menuRef.current) {
      computePosition(anchorRef.current, menuRef.current, {
        placement: 'bottom-start',
        middleware: [offset(8), flip(), shift()],
      }).then(({ x, y }) => {
        setMenuPosition({ top: y, left: x });
      });
    }
  }, []);

  const showMenuForVariableNode = useCallback(
    (_variableNode: VariableNode, targetElement: HTMLElement) => {
      anchorRef.current = targetElement;
      setShowSuggestion(true);
      setMenuItems(selectableOptions);
      setSelectedIndex(0);

      if (closeTimeoutRef.current !== null) {
        clearTimeout(closeTimeoutRef.current);
        closeTimeoutRef.current = null;
      }
      requestAnimationFrame(updateMenuPosition);
    },
    [selectableOptions, updateMenuPosition],
  );

  const closeMenu = useCallback(() => {
    if (closeTimeoutRef.current === null) {
      closeTimeoutRef.current = window.setTimeout(() => {
        setShowSuggestion(false);
        setMenuPosition(null);
        closeTimeoutRef.current = null;
      }, 100);
    }
  }, []);

  const insertVariableNode = useCallback(
    (variable: SelectOption, isNew: boolean) => {
      editor.update(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          const anchorNode = selection.anchor.getNode();

          if (anchorNode.getTextContent().endsWith('{')) {
            const textContent = anchorNode.getTextContent();

            (anchorNode as TextNode).setTextContent(textContent.slice(0, -1));
            selection.anchor.offset--;
            selection.focus.offset--;
          }

          const variableNode = $createVariableNode(variable, isNew);

          selection.insertNodes([variableNode]);

          const spaceNode = $createTextNode('');

          variableNode.insertAfter(spaceNode);

          const newSelection = $createRangeSelection();

          newSelection.anchor.set(spaceNode.getKey(), 1, 'text');
          newSelection.focus.set(spaceNode.getKey(), 1, 'text');
          $setSelection(newSelection);

          requestAnimationFrame(() => {
            const element = editor.getElementByKey(
              variableNode.getKey(),
            ) as HTMLElement;

            if (element) {
              showMenuForVariableNode(variableNode, element);
            }
          });
        }
      });
    },
    [editor, showMenuForVariableNode],
  );

  const onSelectOption = useCallback(
    (selectedOption: SelectOption) => {
      editor.update(() => {
        const selection = $getSelection();

        if (selection) {
          const nodes = selection.getNodes();
          const lastNode = nodes[nodes.length - 1];

          if (lastNode && $isVariableNode(lastNode)) {
            const newVariableNode = $createVariableNode(selectedOption);

            lastNode.replace(newVariableNode);

            const spaceNode = $createTextNode(' ');

            newVariableNode.insertAfter(spaceNode);
            newVariableNode.setSelected();

            const newSelection = $createRangeSelection();

            newSelection.anchor.set(spaceNode.getKey(), 1, 'text');
            newSelection.focus.set(spaceNode.getKey(), 1, 'text');
            $setSelection(newSelection);
          } else {
            insertVariableNode(selectedOption, false);
          }
        }
      });
      setShowSuggestion(false);
    },
    [editor, insertVariableNode],
  );

  useEffect(() => {
    return mergeRegister(
      editor.registerTextContentListener((textContent) => {
        const lastChar = textContent[textContent.length - 1];

        if (lastChar === '{') {
          setMenuItems(selectableOptions);
          setShowSuggestion(true);
          setSelectedIndex(0);
          insertVariableNode(selectableOptions[0], true);
        }
      }),

      editor.registerCommand(
        SELECTION_CHANGE_COMMAND,
        () => {
          const selection = $getSelection();

          if ($isRangeSelection(selection)) {
            const nodes = selection.getNodes();

            if (nodes.length === 1 && $isVariableNode(nodes[0])) {
              const element = editor.getElementByKey(
                nodes[0].getKey(),
              ) as HTMLElement;

              if (element) {
                showMenuForVariableNode(nodes[0], element);
              }

              return true;
            }
          }
          closeMenu();

          return false;
        },
        COMMAND_PRIORITY_HIGH,
      ),

      editor.registerCommand(
        KEY_BACKSPACE_COMMAND,
        (event: KeyboardEvent) => {
          const selection = $getSelection();

          if ($isRangeSelection(selection)) {
            const nodes = selection.getNodes();

            if (nodes.length === 1 && $isVariableNode(nodes[0])) {
              event.preventDefault();
              editor.update(() => {
                const variableNode = nodes[0] as VariableNode;
                const parentNode = variableNode.getParent();

                if (parentNode) {
                  const prevSibling = variableNode.getPreviousSibling();
                  const nextSibling = variableNode.getNextSibling();

                  variableNode.remove();

                  const newSelection = $createRangeSelection();

                  if (prevSibling) {
                    newSelection.anchor.set(
                      prevSibling.getKey(),
                      prevSibling.getTextContentSize(),
                      'text',
                    );
                    newSelection.focus.set(
                      prevSibling.getKey(),
                      prevSibling.getTextContentSize(),
                      'text',
                    );
                  } else if (nextSibling) {
                    newSelection.anchor.set(nextSibling.getKey(), 0, 'text');
                    newSelection.focus.set(nextSibling.getKey(), 0, 'text');
                  } else {
                    const newTextNode = $createTextNode('');

                    parentNode.append(newTextNode);
                    newSelection.anchor.set(newTextNode.getKey(), 0, 'text');
                    newSelection.focus.set(newTextNode.getKey(), 0, 'text');
                  }
                  $setSelection(newSelection);
                }
              });
              closeMenu();

              return true;
            }
          }

          return false;
        },
        COMMAND_PRIORITY_HIGH,
      ),
    );
  }, [
    editor,
    selectableOptions,
    insertVariableNode,
    showMenuForVariableNode,
    closeMenu,
  ]);

  useEffect(() => {
    const unregisterAll = mergeRegister(
      editor.registerCommand(
        KEY_ARROW_DOWN_COMMAND,
        (event: KeyboardEvent) => {
          if (!showSuggestion) return false;
          event.preventDefault();
          setSelectedIndex((prevIndex) =>
            prevIndex < menuItems.length - 1 ? prevIndex + 1 : 0,
          );

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),

      editor.registerCommand(
        KEY_ARROW_UP_COMMAND,
        (event: KeyboardEvent) => {
          if (!showSuggestion) return false;
          event.preventDefault();
          setSelectedIndex((prevIndex) =>
            prevIndex > 0 ? prevIndex - 1 : menuItems.length - 1,
          );

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),

      editor.registerCommand(
        KEY_ENTER_COMMAND,
        (event: KeyboardEvent | null) => {
          if (!showSuggestion) return false;
          event?.preventDefault();
          onSelectOption(menuItems[selectedIndex]);

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),

      editor.registerCommand(
        KEY_TAB_COMMAND,
        (event: KeyboardEvent) => {
          if (!showSuggestion) return false;
          event.preventDefault();
          onSelectOption(menuItems[selectedIndex]);
          closeMenu();

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),

      editor.registerCommand(
        KEY_SPACE_COMMAND,
        (event: KeyboardEvent) => {
          if (!showSuggestion) return false;
          event.preventDefault();
          onSelectOption(menuItems[selectedIndex]);
          closeMenu();

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),

      editor.registerCommand(
        KEY_ESCAPE_COMMAND,
        (event: KeyboardEvent) => {
          if (!showSuggestion) return false;
          event.preventDefault();
          closeMenu();

          return true;
        },
        COMMAND_PRIORITY_HIGH,
      ),
    );

    return () => unregisterAll();
  }, [
    editor,
    showSuggestion,
    menuItems,
    selectedIndex,
    onSelectOption,
    closeMenu,
  ]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        menuRef.current &&
        !menuRef.current.contains(event.target as Node) &&
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
      {showSuggestion &&
        createPortal(
          <div
            ref={menuRef}
            tabIndex={-1}
            className='bg-white min-w-[250px] py-1.5 px-[6px] shadow-lg border absolute rounded-md z-[99999]'
            style={{
              pointerEvents: 'all',
              visibility: menuPosition ? 'visible' : 'hidden',
              top: menuPosition ? `${menuPosition.top}px` : 0,
              left: menuPosition ? `${menuPosition.left}px` : 0,
            }}
          >
            <ul>
              {menuItems.map((option, i) => (
                <VariablesMenuItem
                  option={option}
                  key={option.value}
                  isSelected={selectedIndex === i}
                  onClick={() => onSelectOption(option)}
                  onMouseEnter={() => setSelectedIndex(i)}
                />
              ))}
            </ul>
          </div>,
          document.body,
        )}
    </>
  );
}

function VariablesMenuItem({
  option,
  isSelected,
  onClick,
  onMouseEnter,
}: {
  isSelected: boolean;
  onClick: () => void;
  option: SelectOption;
  onMouseEnter: () => void;
}) {
  const [type, ...rest] = option.label.split('_');

  return (
    <li
      tabIndex={1}
      role='option'
      onClick={onClick}
      aria-selected={isSelected}
      onMouseEnter={onMouseEnter}
      className={cn(
        'group flex gap-2 items-center text-start py-[6px] px-[10px] leading-[18px] text-gray-700 rounded-sm outline-none cursor-pointer hover:bg-gray-50 hover:rounded-md',
        isSelected && 'bg-gray-50 text-gray-700',
      )}
    >
      {type === 'contact' && (
        <User01
          className={cn('size-4 text-gray-500 group-hover:text-gray-700', {
            'text-gray-700': isSelected,
          })}
        />
      )}

      {type === 'organization' && (
        <Building07
          className={cn('size-4 text-gray-500 group-hover:text-gray-700', {
            'text-gray-700': isSelected,
          })}
        />
      )}
      <span className='text-sm first-letter:capitalize'>{rest?.join(' ')}</span>

      <span className='text-sm text-gray-500 '>{`{{${option.value}}}`}</span>
    </li>
  );
}

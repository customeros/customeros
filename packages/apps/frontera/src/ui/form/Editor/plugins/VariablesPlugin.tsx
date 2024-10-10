import * as ReactDOM from 'react-dom';
import { useMemo, useCallback } from 'react';

import { TextNode } from 'lexical';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  MenuOption,
  MenuTextMatch,
  LexicalTypeaheadMenuPlugin,
} from '@lexical/react/LexicalTypeaheadMenuPlugin';

import { cn } from '@ui/utils/cn';
import { SelectOption } from '@ui/utils/types';
import { $createVariableNode } from '@ui/form/Editor/nodes/VariableNode.tsx';

const TRIGGERS = '{{';

// Chars we expect to see in a variable name (non-space, non-punctuation).
const VALID_CHARS = '[^}}\\s]';

const LENGTH_LIMIT = 75;

const VariableMentionsRegex = new RegExp(
  '(^|\\s|\\()(' +
    TRIGGERS +
    '((?:' +
    VALID_CHARS +
    '){0,' +
    LENGTH_LIMIT +
    '})' +
    ')$',
);

// At most, 5 suggestions are shown in the popup.
const SUGGESTION_LIST_LENGTH_LIMIT = 5;

function checkForVariableMentions(
  text: string,
  minMatchLength: number,
): MenuTextMatch | null {
  const match = VariableMentionsRegex.exec(text);

  if (match !== null) {
    const maybeLeadingWhitespace = match[1];
    const matchingString = match[3];

    if (matchingString.length >= minMatchLength) {
      return {
        leadOffset: match.index + maybeLeadingWhitespace.length,
        matchingString,
        replaceableString: match[2],
      };
    }
  }

  return null;
}

function getPossibleQueryMatch(
  text: string,
  minMatch: number,
): MenuTextMatch | null {
  return checkForVariableMentions(text, minMatch);
}

class VariableTypeaheadOption extends MenuOption {
  label: string;
  value: string;

  constructor(item: SelectOption) {
    super(item.label);
    this.label = item.label;
    this.value = item.value;
  }
}

function VariablesTypeaheadMenuItem({
  index,
  option,
  selectOption,
  isSelected,
  onMouseEnter,
}: {
  index: number;
  isSelected: boolean;
  selectOption: () => void;
  onMouseEnter: () => void;
  option: VariableTypeaheadOption;
}) {
  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    selectOption();
  };

  return (
    <li
      tabIndex={-1}
      role='option'
      key={option.key}
      onMouseDown={handleClick}
      ref={option.setRefElement}
      onMouseEnter={onMouseEnter}
      style={{ pointerEvents: 'all' }}
      id={'typeahead-variable-item-' + index}
      className={cn(
        'flex gap-2 items-center text-start py-[6px] px-[10px] leading-[18px] text-gray-700 rounded-sm outline-none cursor-pointer hover:bg-gray-50 hover:rounded-md ',
        'data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed hover:data-[disabled]:bg-transparent',
        isSelected && 'bg-gray-50 text-gray-700',
      )}
    >
      <span className='text first-letter:capitalize'>{option.label}</span>
    </li>
  );
}

interface VariablesPluginProps {
  options: string[];
  onSearch?: (query: string | null) => void;
}

export default function VariablesPlugin({
  options = [],
  onSearch,
}: VariablesPluginProps): JSX.Element | null {
  const [editor] = useLexicalComposerContext();

  const selectableOptions = options.map((item) => ({
    label: item.split('_').join(' ').toLowerCase(),
    value: item.toLowerCase(),
  }));

  const _options = useMemo(
    () =>
      selectableOptions
        .map((item) => new VariableTypeaheadOption(item))
        .slice(0, SUGGESTION_LIST_LENGTH_LIMIT),
    [selectableOptions],
  );

  const onSelectOption = useCallback(
    (
      selectedOption: VariableTypeaheadOption,
      nodeToReplace: TextNode | null,
      closeMenu: () => void,
    ) => {
      editor.update(() => {
        const variableNode = $createVariableNode({
          label: selectedOption.value,
          value: selectedOption.value,
        });

        if (nodeToReplace) {
          nodeToReplace.replace(variableNode);
        }

        variableNode.select();
        closeMenu();
      });
    },
    [editor],
  );

  const checkForVariableMatch = useCallback(
    (text: string): MenuTextMatch | null => {
      return getPossibleQueryMatch(text, 1);
    },
    [],
  );

  return (
    <>
      <LexicalTypeaheadMenuPlugin<VariableTypeaheadOption>
        options={_options}
        onSelectOption={onSelectOption}
        triggerFn={checkForVariableMatch}
        onQueryChange={onSearch ?? (() => {})}
        menuRenderFn={(
          anchorElementRef,
          {
            selectedIndex,
            selectOptionAndCleanUp,
            setHighlightedIndex,
            options: opts,
          },
        ) =>
          anchorElementRef.current && _options.length
            ? ReactDOM.createPortal(
                <div
                  style={{}}
                  data-side='bottom'
                  className='relative bg-white min-w-[250px] py-1.5 px-[6px] shadow-lg border rounded-md data-[side=top]:animate-slideDownAndFade data-[side=right]:animate-slideLeftAndFade data-[side=bottom]:animate-slideUpAndFade data-[side=left]:animate-slideRightAndFade z-[500]'
                >
                  <ul>
                    {opts.map((option, i: number) => (
                      <VariablesTypeaheadMenuItem
                        index={i}
                        option={option}
                        key={option.key}
                        isSelected={selectedIndex === i}
                        onMouseEnter={() => {
                          setHighlightedIndex(i);
                        }}
                        selectOption={() => {
                          setHighlightedIndex(i);
                          selectOptionAndCleanUp(option);
                        }}
                      />
                    ))}
                  </ul>
                </div>,
                anchorElementRef.current,
              )
            : null
        }
      />
    </>
  );
}

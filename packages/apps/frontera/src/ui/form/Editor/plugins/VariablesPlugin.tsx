import * as ReactDOM from 'react-dom';
import { useMemo, useState, useCallback } from 'react';

import { TextNode } from 'lexical';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  MenuOption,
  MenuTextMatch,
  LexicalTypeaheadMenuPlugin,
  useBasicTypeaheadTriggerMatch,
} from '@lexical/react/LexicalTypeaheadMenuPlugin';

import { cn } from '@ui/utils/cn';

import { $createVariableNode } from '../nodes/VariableNode';

const PUNCTUATION =
  '\\.,\\+\\*\\?\\$\\@\\|#{}\\(\\)\\^\\-\\[\\]\\\\/!%\'"~=<>_:;';
const NAME = '\\b[A-Z][^\\s' + PUNCTUATION + ']';

const DocumentMentionsRegex = {
  NAME,
  PUNCTUATION,
};

const PUNC = DocumentMentionsRegex.PUNCTUATION;

const TRIGGERS = ['{{'].join('');

// Chars we expect to see in a mention (non-space, non-punctuation).
const VALID_CHARS = '[^' + TRIGGERS + PUNC + '\\s]';

// Non-standard series of chars. Each series must be preceded and followed by
// a valid char.
const VALID_JOINS =
  '(?:' +
  '\\.[ |$]|' + // E.g. "r. " in "Mr. Smith"
  ' |' + // E.g. " " in "Josh Duck"
  '[' +
  PUNC +
  ']|' + // E.g. "-' in "Salier-Hellendag"
  ')';

const LENGTH_LIMIT = 75;

const VariableSignRegex = new RegExp(
  '(^|\\s|\\()(' +
    TRIGGERS +
    '((?:' +
    VALID_CHARS +
    VALID_JOINS +
    '){0,' +
    LENGTH_LIMIT +
    '})' +
    ')$',
);

// 50 is the longest alias length limit.
const ALIAS_LENGTH_LIMIT = 50;

// Regex used to match alias.
const VariableSignRegexAliasRegex = new RegExp(
  '(^|\\s|\\()(' +
    TRIGGERS +
    '((?:' +
    VALID_CHARS +
    '){0,' +
    ALIAS_LENGTH_LIMIT +
    '})' +
    ')$',
);

function checkForVariableSign(
  text: string,
  minMatchLength: number,
): MenuTextMatch | null {
  let match = VariableSignRegex.exec(text);

  if (match === null) {
    match = VariableSignRegexAliasRegex.exec(text);
  }

  if (match !== null) {
    // The strategy ignores leading whitespace but we need to know it's
    // length to add it to the leadOffset
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

function getPossibleQueryMatch(text: string): MenuTextMatch | null {
  return checkForVariableSign(text, 0);
}

function sanitizeVariableName(variableName: string) {
  if (new RegExp(/^{{.*}}$/).exec(variableName)) return variableName;

  return '{{' + variableName.toLowerCase() + '}}';
}

class VariableTypeaheadOption extends MenuOption {
  name: string;

  constructor(name: string) {
    const _name = sanitizeVariableName(name);

    super(_name);
    this.name = _name;
  }
}

function VariablesTypeaheadMenuItem({
  index,
  isSelected,
  onClick,
  onMouseEnter,
  option,
}: {
  index: number;
  isSelected: boolean;
  onClick: () => void;
  onMouseEnter: () => void;
  option: VariableTypeaheadOption;
}) {
  return (
    <li
      tabIndex={-1}
      role='option'
      key={option.key}
      onMouseDown={onClick}
      ref={option.setRefElement}
      aria-selected={isSelected}
      onMouseEnter={onMouseEnter}
      id={'variable-item-' + index}
      className={cn(
        'flex gap-2 items-center text-start py-[6px] px-[10px] leading-[18px] text-gray-700  rounded-sm outline-none cursor-pointer hover:bg-gray-50 hover:rounded-md ',
        'data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed hover:data-[disabled]:bg-transparent',
        isSelected && 'bg-gray-50 text-gray-700',
      )}
    >
      {option.name}
    </li>
  );
}

interface MentionsPluginProps {
  options: string[];
}

export default function VariablesPlugin({
  options,
}: MentionsPluginProps): JSX.Element | null {
  const [editor] = useLexicalComposerContext();
  const [queryString, setQueryString] = useState<string | null>(null);

  const checkForSlashTriggerMatch = useBasicTypeaheadTriggerMatch('/', {
    minLength: 0,
  });

  const _options = useMemo(() => {
    const res = options.map((item) => new VariableTypeaheadOption(item));

    if (queryString) {
      return res.filter((item) =>
        item.name.toLowerCase()?.includes(queryString?.toLowerCase()),
      );
    }

    return res;
  }, [options, queryString]);

  const onSelectOption = useCallback(
    (
      selectedOption: VariableTypeaheadOption,
      nodeToReplace: TextNode | null,
      closeMenu: () => void,
    ) => {
      editor.update(() => {
        const mentionNode = $createVariableNode(selectedOption.name);

        if (nodeToReplace) {
          nodeToReplace.replace(mentionNode);
        }
        mentionNode.select();
        closeMenu();
      });
    },
    [editor],
  );

  const checkForMentionMatch = useCallback(
    (text: string) => {
      const slashMatch = checkForSlashTriggerMatch(text, editor);

      if (slashMatch !== null) {
        return null;
      }

      return getPossibleQueryMatch(text);
    },
    [checkForSlashTriggerMatch, editor],
  );

  return (
    <LexicalTypeaheadMenuPlugin<VariableTypeaheadOption>
      options={_options}
      onQueryChange={setQueryString}
      onSelectOption={onSelectOption}
      triggerFn={checkForMentionMatch}
      menuRenderFn={(
        anchorElementRef,
        { selectedIndex, selectOptionAndCleanUp, setHighlightedIndex },
      ) =>
        anchorElementRef.current && _options.length
          ? ReactDOM.createPortal(
              <div className='relative bg-white min-w-[250px] py-1.5 px-[6px] shadow-lg border rounded-md z-50'>
                <ul>
                  {_options.map((option, i: number) => (
                    <VariablesTypeaheadMenuItem
                      index={i}
                      option={option}
                      key={option.key}
                      isSelected={selectedIndex === i}
                      onMouseEnter={() => {
                        setHighlightedIndex(i);
                      }}
                      onClick={() => {
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
  );
}

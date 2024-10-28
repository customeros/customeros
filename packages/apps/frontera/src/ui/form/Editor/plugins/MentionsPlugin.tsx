import * as ReactDOM from 'react-dom';
import { useMemo, useCallback } from 'react';

import { TextNode } from 'lexical';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
import {
  MenuOption,
  MenuTextMatch,
  LexicalTypeaheadMenuPlugin,
  useBasicTypeaheadTriggerMatch,
} from '@lexical/react/LexicalTypeaheadMenuPlugin';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';

import { $createMentionNode } from '../nodes/MentionNode';

const PUNCTUATION =
  '\\.,\\+\\*\\?\\$\\@\\|#{}\\(\\)\\^\\-\\[\\]\\\\/!%\'"~=<>_:;';
const NAME = '\\b[A-Z][^\\s' + PUNCTUATION + ']';

const DocumentMentionsRegex = {
  NAME,
  PUNCTUATION,
};

const PUNC = DocumentMentionsRegex.PUNCTUATION;

const TRIGGERS = ['@'].join('');

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

const AtSignMentionsRegex = new RegExp(
  '(^|\\s|\\()(' +
    '[' +
    TRIGGERS +
    ']' +
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
const AtSignMentionsRegexAliasRegex = new RegExp(
  '(^|\\s|\\()(' +
    '[' +
    TRIGGERS +
    ']' +
    '((?:' +
    VALID_CHARS +
    '){0,' +
    ALIAS_LENGTH_LIMIT +
    '})' +
    ')$',
);

// At most, 5 suggestions are shown in the popup.
const SUGGESTION_LIST_LENGTH_LIMIT = 5;

function checkForAtSignMentions(
  text: string,
  minMatchLength: number,
): MenuTextMatch | null {
  let match = AtSignMentionsRegex.exec(text);

  if (match === null) {
    match = AtSignMentionsRegexAliasRegex.exec(text);
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
  return checkForAtSignMentions(text, 1);
}

class MentionTypeaheadOption extends MenuOption {
  name: string;
  picture: JSX.Element;

  constructor(name: string, picture: JSX.Element) {
    super(name);
    this.name = name;
    this.picture = picture;
  }
}

function MentionsTypeaheadMenuItem({
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
  option: MentionTypeaheadOption;
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
      id={'typeahead-item-' + index}
      className={cn(
        'flex gap-2 items-center text-start py-[6px] px-[10px] leading-[18px] text-gray-700  rounded-sm outline-none cursor-pointer hover:bg-gray-50 hover:rounded-md ',
        'data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed hover:data-[disabled]:bg-transparent',
        isSelected && 'bg-gray-50 text-gray-700',
      )}
    >
      {option.picture}
      <span className='text'>{option.name}</span>
    </li>
  );
}

interface MentionsPluginProps {
  options: string[];
  onSearch?: (query: string | null) => void;
}

export default function NewMentionsPlugin({
  options,
  onSearch,
}: MentionsPluginProps): JSX.Element | null {
  const [editor] = useLexicalComposerContext();

  const checkForSlashTriggerMatch = useBasicTypeaheadTriggerMatch('/', {
    minLength: 0,
  });

  const _options = useMemo(
    () =>
      options
        .map(
          (item) =>
            new MentionTypeaheadOption(
              item,
              <Avatar size='xs' name={item} textSize='xs' />,
            ),
        )
        .slice(0, SUGGESTION_LIST_LENGTH_LIMIT),
    [options],
  );

  const onSelectOption = useCallback(
    (
      selectedOption: MentionTypeaheadOption,
      nodeToReplace: TextNode | null,
      closeMenu: () => void,
    ) => {
      editor.update(() => {
        const mentionNode = $createMentionNode(selectedOption.name);

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
    <LexicalTypeaheadMenuPlugin<MentionTypeaheadOption>
      options={_options}
      onSelectOption={onSelectOption}
      triggerFn={checkForMentionMatch}
      onQueryChange={onSearch ?? (() => {})}
      menuRenderFn={(
        anchorElementRef,
        { selectedIndex, selectOptionAndCleanUp, setHighlightedIndex },
      ) =>
        anchorElementRef.current && _options.length
          ? ReactDOM.createPortal(
              <div className='relative bg-white min-w-[250px] py-1.5 px-[6px] shadow-lg border rounded-md z-50'>
                <ul>
                  {_options.map((option, i: number) => (
                    <MentionsTypeaheadMenuItem
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

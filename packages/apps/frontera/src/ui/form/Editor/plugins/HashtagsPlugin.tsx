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
import { SelectOption } from '@ui/utils/types';

import { $createHashtagNode } from '../nodes/HashtagNode';

const PUNCTUATION =
  '\\.,\\+\\*\\?\\$\\@\\|#{}\\(\\)\\^\\-\\[\\]\\\\/!%\'"~=<>_:;';
const NAME = '\\b[A-Z][^\\s' + PUNCTUATION + ']';

const DocumentMentionsRegex = {
  NAME,
  PUNCTUATION,
};

const PUNC = DocumentMentionsRegex.PUNCTUATION;

const TRIGGERS = ['#'].join('');

// Chars we expect to see in a mention (non-space, non-punctuation).
const VALID_CHARS = '[^' + TRIGGERS + PUNC + '\\s]';

const LENGTH_LIMIT = 75;

const HashSignMentionsRegex = new RegExp(
  '(^|\\s|\\()(' +
    '[' +
    TRIGGERS +
    ']' +
    '((?:' +
    VALID_CHARS +
    '){0,' +
    LENGTH_LIMIT +
    '})' +
    ')$',
);

// 50 is the longest alias length limit.
const ALIAS_LENGTH_LIMIT = 50;

// Regex used to match alias.
const HashSignMentionsRegexAliasRegex = new RegExp(
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

function checkForHashSignMentions(
  text: string,
  minMatchLength: number,
): MenuTextMatch | null {
  let match = HashSignMentionsRegex.exec(text);

  if (match === null) {
    match = HashSignMentionsRegexAliasRegex.exec(text);
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
  return checkForHashSignMentions(text, 1);
}

class HashtagTypeaheadOption extends MenuOption {
  label: string;
  value: string;

  constructor(item: SelectOption) {
    super(item.label);
    this.label = item.label;
    this.value = item.value;
  }
}

function HashtagsTypeaheadMenuItem({
  index,
  option,
  onClick,
  isSelected,
  onMouseEnter,
}: {
  index: number;
  isSelected: boolean;
  onClick: () => void;
  onMouseEnter: () => void;
  option: HashtagTypeaheadOption;
}) {
  return (
    <li
      key={option.key}
      tabIndex={-1}
      className={cn(
        'flex gap-2 items-center text-start py-[6px] px-[10px] leading-[18px] text-gray-700  rounded-sm outline-none cursor-pointer hover:bg-gray-50 hover:rounded-md ',
        'data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed hover:data-[disabled]:bg-transparent',
        isSelected && 'bg-gray-50 text-gray-700',
      )}
      ref={option.setRefElement}
      role='option'
      aria-selected={isSelected}
      id={'typeahead-hashtag-item-' + index}
      onMouseEnter={onMouseEnter}
      onClick={onClick}
    >
      <span className='text'>{option.label}</span>
    </li>
  );
}

interface HashtagsPluginProps {
  options: SelectOption[];
  onCreate?: (hashtag: string) => void;
  onSearch?: (query: string | null) => void;
}

export default function NewHashtagsPlugin({
  options,
  onSearch,
  onCreate,
}: HashtagsPluginProps): JSX.Element | null {
  const [editor] = useLexicalComposerContext();

  const checkForSlashTriggerMatch = useBasicTypeaheadTriggerMatch('/', {
    minLength: 0,
  });

  const _options = useMemo(
    () =>
      (options.length ? options : [{ label: 'Create tag', value: 'temp' }])
        .map((item) => new HashtagTypeaheadOption(item))
        .slice(0, SUGGESTION_LIST_LENGTH_LIMIT),
    [options],
  );

  const onSelectOption = useCallback(
    (
      selectedOption: HashtagTypeaheadOption,
      nodeToReplace: TextNode | null,
      closeMenu: () => void,
    ) => {
      editor.update(() => {
        const isCreating = selectedOption.value === 'temp';

        const hashtagNode = $createHashtagNode({
          value: selectedOption.value,
          label: isCreating
            ? nodeToReplace?.__text ?? ''
            : selectedOption.label,
        });

        if (isCreating) {
          onCreate?.(nodeToReplace?.__text ?? '');
        }

        if (nodeToReplace) {
          nodeToReplace.replace(hashtagNode);
        }

        hashtagNode.select();
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
    <LexicalTypeaheadMenuPlugin<HashtagTypeaheadOption>
      onQueryChange={onSearch ?? (() => {})}
      onSelectOption={onSelectOption}
      triggerFn={checkForMentionMatch}
      options={_options}
      menuRenderFn={(
        anchorElementRef,
        { selectedIndex, selectOptionAndCleanUp, setHighlightedIndex },
      ) =>
        anchorElementRef.current
          ? ReactDOM.createPortal(
              <div className='bg-white min-w-[250px] py-1.5 px-[6px] shadow-lg border rounded-md data-[side=top]:animate-slideDownAndFade data-[side=right]:animate-slideLeftAndFade data-[side=bottom]:animate-slideUpAndFade data-[side=left]:animate-slideRightAndFade z-10'>
                <ul>
                  {_options.map((option, i: number) => (
                    <HashtagsTypeaheadMenuItem
                      index={i}
                      isSelected={selectedIndex === i}
                      onClick={() => {
                        setHighlightedIndex(i);
                        selectOptionAndCleanUp(option);
                      }}
                      onMouseEnter={() => {
                        setHighlightedIndex(i);
                      }}
                      key={option.key}
                      option={option}
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

import { useField } from 'react-inverted-form';
import { FC, useRef, useState, useEffect, KeyboardEvent } from 'react';
import {
  GroupBase,
  components,
  OptionProps,
  MenuListProps,
  OptionsOrGroups,
} from 'react-select';

import { AnimatePresence } from 'framer-motion';

import { cn } from '@ui/utils/cn';
import { SelectOption } from '@ui/utils/types';
import { getMenuClassNames } from '@ui/form/Select';
import {
  MultiCreatableSelect,
  getMenuListClassNames,
} from '@ui/form/MultiCreatableSelect/MultiCreatableSelect';

import { TagButton } from './TagButton';
import { useTagButtonSlideAnimation } from './useTagButtonSlideAnimation';

interface EmailParticipantSelect {
  name: string;
  formId: string;
  tags?: Array<{ value: string; label: string }>;
}

interface Tag {
  label: string;
  value: string;
}
export const suggestedTags = [
  'meeting',
  'call',
  'voicemail',
  'email',
  'text-message',
];

export const TagsSelect: FC<EmailParticipantSelect> = ({
  formId,
  name,
  tags = [],
}) => {
  const { getInputProps } = useField(name, formId);
  const { onChange, value: selectedTags, onBlur } = getInputProps();
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [focusedOption, setFocusedOption] = useState<Tag | null>(null);
  const [inputVal, setInputVal] = useState('');
  const scope = useTagButtonSlideAnimation(!!selectedTags?.length);

  const getFilteredSuggestions = (
    filterString: string,
    callback: (options: OptionsOrGroups<unknown, GroupBase<unknown>>) => void,
  ) => {
    if (!filterString.slice(1).length) {
      callback(tags);

      return;
    }

    const options: OptionsOrGroups<unknown, GroupBase<unknown>> = tags.filter(
      (e) =>
        e.label.toLowerCase().includes(filterString.slice(1)?.toLowerCase()),
    );

    callback(options);
  };
  const handleInputChange = (d: string) => {
    setInputVal(d);
    if (d.length === 1 && d.startsWith('#')) {
      setMenuOpen(true);
    }
    if (!d.length || !d.startsWith('#')) {
      setMenuOpen(false);
    }
  };

  // this function is needed as tags are selected on 'Space' & 'Enter'
  const handleKeyDown = (event: KeyboardEvent) => {
    if (event.code === 'Backspace') {
      if (inputVal.length) {
        return;
      }
      event.preventDefault();
      const newSelected = [...selectedTags].slice(0, selectedTags.length - 1);
      onChange(newSelected);
    }
    if (event.code === 'Space' || event.code === 'Enter') {
      event.preventDefault();
      if (!isMenuOpen) return;

      if (focusedOption) {
        onChange([...selectedTags, focusedOption]);
        setMenuOpen(false);
        setFocusedOption(null);
        setInputVal('');
      }
    }
  };

  // FIXME - move this to outer scope
  const Option = (props: OptionProps<SelectOption>) => {
    const Or = useRef(null);

    useEffect(() => {
      if (props.isFocused) {
        setFocusedOption(props.data);
      }
    }, [props.isFocused, props.data.label]);

    return (
      <div ref={Or}>
        <components.Option {...props} key={props.data.label}>
          {props.data.label || props.data.value}
        </components.Option>
      </div>
    );
  };

  return (
    <>
      <AnimatePresence initial={false}>
        <div className='flex items-baseline' ref={scope}>
          {!selectedTags?.length && (
            <>
              <p className='text-gray-500 mr-2 whitespace-nowrap'>
                Suggested tags:
              </p>

              {suggestedTags?.map((tag) => (
                <TagButton
                  key={`tag-select-${tag}`}
                  onTagSet={() =>
                    onChange([
                      {
                        label: tag,
                        value:
                          tags?.find((e) => suggestedTags.includes(e.label))
                            ?.value || tag,
                      },
                    ])
                  }
                  tag={tag}
                />
              ))}
            </>
          )}
          {!!selectedTags?.length && (
            <MultiCreatableSelect
              formId={formId}
              name={name}
              Option={Option}
              classNames={{
                menu: ({ menuPlacement }) =>
                  getMenuClassNames(menuPlacement)('!z-[11]'),
                multiValueLabel: () =>
                  'p-0 gap-0 text-gray-700 m-0 mr-1 cursor-text font-base leading-4 before:content-["#"]',

                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                menuList: (props: MenuListProps<any, any, any>) =>
                  getMenuListClassNames(
                    cn({
                      'absolute top-[-300px] z-[999]':
                        props?.options?.length === 1 &&
                        props?.options?.[0]?.label ===
                          props.options?.[0]?.value &&
                        !suggestedTags.includes(props.options?.[0]?.label),
                    }),
                  ),
              }}
              placeholder='#Tag'
              backspaceRemovesValue
              onKeyDown={handleKeyDown}
              onChange={onChange}
              noOptionsMessage={() => null}
              loadOptions={(inputValue: string, callback) => {
                getFilteredSuggestions(inputValue, callback);
              }}
              formatCreateLabel={(input) => {
                if (input?.startsWith('#')) {
                  return `${input.slice(1)}`;
                }

                return input;
              }}
              onBlur={() => onBlur(selectedTags)}
              onMenuClose={() => setFocusedOption(null)}
              value={selectedTags}
              inputValue={inputVal}
              onInputChange={handleInputChange}
              menuIsOpen={isMenuOpen}
              menuPlacement='top'
              defaultOptions={tags}
              hideSelectedOptions
              isValidNewOption={(input) => input.startsWith('#')}
              getOptionLabel={(d) => {
                if (d.label?.startsWith('#')) {
                  return `${d.label.slice(1)}`;
                }

                return `${d.label}`;
              }}
              menuShouldBlockScroll
              onCreateOption={(input) => {
                if (input?.startsWith('#')) {
                  return {
                    value: input.slice(1),
                    label: input.slice(1),
                  };
                }

                return {
                  value: input,
                  label: input,
                };
              }}
              getNewOptionData={(input) => {
                if (input?.startsWith('#')) {
                  return {
                    value: input.slice(1),
                    label: input.slice(1),
                  };
                }

                return {
                  value: input,
                  label: input,
                };
              }}
            />
          )}
        </div>
      </AnimatePresence>
    </>
  );
};

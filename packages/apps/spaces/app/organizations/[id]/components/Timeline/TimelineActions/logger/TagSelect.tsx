'use client';
import React, { FC, KeyboardEvent, useEffect, useRef, useState } from 'react';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { OptionsOrGroups } from 'react-select';
import { AnimatePresence } from 'framer-motion';
import { chakraComponents } from '@ui/form/SyncSelect';
import { MultiCreatableSelect } from '@ui/form/MultiCreatableSelect';
import { tagsSelectStyles } from './tagSelectStyles';
import { TagButton } from './TagButton';
import { useTagButtonSlideAnimation } from './useTagButtonSlideAnimation';

interface EmailParticipantSelect {
  autofocus?: boolean;
}

interface Tag {
  label: string;
  value: string;
}
const commonTags = ['meeting', 'call', 'voicemail', 'email', 'text-message'];
const tags = commonTags.map((label) => ({ label, value: label }));

export const TagsSelect: FC<EmailParticipantSelect> = ({
  autofocus = false,
}) => {
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [focusedOption, setFocusedOption] = useState<Tag | null>(null);
  const [inputVal, setInputVal] = useState('');
  const [selectedTags, setTags] = useState<Array<Tag>>([]);
  const scope = useTagButtonSlideAnimation(!!selectedTags.length);
  const getFilteredSuggestions = (
    filterString: string,
    callback: (options: OptionsOrGroups<any, any>) => void,
  ) => {
    if (!filterString.slice(1).length) {
      callback(tags);
      return;
    }

    const options: OptionsOrGroups<string, any> = commonTags
      .map((label) => ({ label, value: label }))
      .filter((e) =>
        e.value.toLowerCase().includes(filterString.slice(1)?.toLowerCase()),
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

  const handleKeyDown = (event: KeyboardEvent) => {
    if (event.code === 'Enter') {
      event.preventDefault();
    }
    if (event.code === 'Space') {
      event.preventDefault();
      if (!isMenuOpen) return;

      if (focusedOption) {
        setTags((prevTags) => [...prevTags, focusedOption]);
        setMenuOpen(false);
        setFocusedOption(null);
        setInputVal('');
      }
    }
  };

  const Option = (props: any) => {
    const Or = useRef(null);

    useEffect(() => {
      if (props.isFocused) {
        setFocusedOption(props.data);
      }
    }, [props.isFocused, props.data.label]);

    return (
      <chakraComponents.Option {...props} key={props.data.label} ref={Or}>
        {props.data.label || props.data.value}
      </chakraComponents.Option>
    );
  };

  return (
    <>
      <AnimatePresence initial={false}>
        <Flex alignItems='center' ref={scope}>
          {!selectedTags.length && (
            <>
              <Text color='gray.500' mr={2} whiteSpace='nowrap'>
                Suggested tags:
              </Text>

              {tags.map((tag) => (
                <TagButton
                  key={`tag-select-${tag.value}`}
                  onTagSet={() => setTags([tag])}
                  tag={tag}
                />
              ))}
            </>
          )}
          {!!selectedTags.length && (
            <MultiCreatableSelect
              Option={Option}
              autoFocus={autofocus}
              name=''
              formId=''
              placeholder=''
              onKeyDown={handleKeyDown}
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
              onMenuClose={() => setFocusedOption(null)}
              value={selectedTags}
              onChange={(e) => {
                setTags(e);
              }}
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
              customStyles={tagsSelectStyles}
            />
          )}
        </Flex>
      </AnimatePresence>
    </>
  );
};

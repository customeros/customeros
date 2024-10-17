import { useMemo, forwardRef, useCallback } from 'react';
import {
  GroupBase,
  OptionProps,
  SelectInstance,
  OptionsOrGroups,
  MultiValueProps,
  components as reactSelectComponents,
} from 'react-select';

import merge from 'lodash/merge';

import { SelectOption } from '@ui/utils/types.ts';
import { Copy01 } from '@ui/media/icons/Copy01.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import {
  CreatableSelect,
  getDefaultClassNames,
} from '@ui/form/CreatableSelect';

import { MultiValueWithActionMenu } from './MultiValueWithActionMenu.tsx';

export const EmailMultiCreatableSelect = forwardRef<
  SelectInstance,
  {
    isMulti: boolean;
    placeholder?: string;
    noOptionsMessage: () => null;
    value: SelectOption<string>[];
    navigateAfterAddingToPeople: boolean;
    onChange: (value: SelectOption<string>[]) => void;
    onKeyDown: (e: React.KeyboardEvent<HTMLDivElement>) => void;
    options: Array<{ id: string; value: string; label: string }>;
  }
>(
  (
    {
      value,
      onChange,
      navigateAfterAddingToPeople,
      isMulti,
      onKeyDown,
      options,
      ...rest
    },
    ref,
  ) => {
    const [_, copyToClipboard] = useCopyToClipboard();

    const getFilteredSuggestions = (
      filterString: string,
      callback: (options: OptionsOrGroups<unknown, GroupBase<unknown>>) => void,
    ) => {
      if (!filterString.slice(1).length) {
        callback(options);

        return;
      }

      const opt: OptionsOrGroups<unknown, GroupBase<unknown>> = options.filter(
        (e) =>
          e.label
            .toLowerCase()
            .includes(filterString.slice(1)?.toLowerCase()) ||
          e.value.toLowerCase().includes(filterString.slice(1)?.toLowerCase()),
      );

      callback(opt);
    };
    const Option = useCallback((rest: OptionProps<SelectOption>) => {
      const fullLabel =
        rest?.data?.label.length > 1 &&
        rest?.data?.value.length > 1 &&
        `${rest.data.label}  ${rest.data.value}`;

      const emailOnly = rest?.data?.value.length > 1 && `${rest.data.value}`;

      const noEmail = rest?.data?.label && !rest?.data?.value && (
        <p>
          {rest.data.label} -
          <span className='text-gray-500 ml-1'>
            [No email for this contact]
          </span>
        </p>
      );

      return (
        <reactSelectComponents.Option {...rest}>
          {fullLabel || emailOnly || noEmail}
          {rest?.isFocused && (
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='Copy'
              className='h-5 p-0 self-end float-end'
              icon={<Copy01 className='size-3 text-gray-500' />}
              onClick={(e) => {
                e.stopPropagation();
                copyToClipboard(rest.data.value, 'Email copied');
              }}
            />
          )}
        </reactSelectComponents.Option>
      );
    }, []);

    const MultiValue = useCallback(
      (multiValueProps: MultiValueProps<SelectOption>) => {
        return (
          <MultiValueWithActionMenu
            {...multiValueProps}
            value={value}
            onChange={onChange}
            existingContacts={options}
            navigateAfterAddingToPeople={navigateAfterAddingToPeople}
          />
        );
      },
      [navigateAfterAddingToPeople],
    );

    const components = useMemo(
      () => ({
        MultiValueRemove: () => null,
        LoadingIndicator: () => null,
        DropdownIndicator: () => null,
        Option,
        MultiValue,
      }),
      [MultiValue, Option],
    );
    const defaultClassNames = useMemo(
      () => merge(getDefaultClassNames({ size: 'md' })),
      [],
    );

    return (
      <CreatableSelect
        unstyled
        menuIsOpen
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        ref={ref as any}
        defaultMenuIsOpen
        options={options}
        isClearable={false}
        onChange={onChange}
        tabSelectsValue={true}
        components={components}
        id={'email-multi-creatable-select'}
        classNames={{
          ...defaultClassNames,
        }}
        formatCreateLabel={(input: string) => {
          return input;
        }}
        loadOptions={(inputValue: string, callback) => {
          getFilteredSuggestions(inputValue, callback);
        }}
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        value={value?.map((e: { value: any; label: string | any[] }) => ({
          label: e.label.length > 1 ? e.label : e.value,
          value: e.value,
        }))}
        onKeyDown={(e) => {
          if (!isMulti && e.key !== 'Backspace' && value.length > 0) {
            e.stopPropagation();
            e.preventDefault();

            return;
          }

          if (onKeyDown) onKeyDown(e);
          e.stopPropagation();
        }}
        {...rest}
      />
    );
  },
);

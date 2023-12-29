import { useField } from 'react-inverted-form';
import React, { useMemo, forwardRef, useCallback } from 'react';

import { OptionProps, MultiValueProps } from 'chakra-react-select';

import { SelectOption } from '@ui/utils';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { chakraComponents } from '@ui/form/SyncSelect';
import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { multiCreatableSelectStyles } from '@ui/form/MultiCreatableSelect/styles';
import { emailRegex } from '@organization/src/components/Timeline/events/email/utils';
import {
  FormSelectProps,
  MultiCreatableSelect,
} from '@ui/form/MultiCreatableSelect';
import {
  Menu,
  MenuItem,
  MenuButton,
  MenuList as ChakraMenuList,
} from '@ui/overlay/Menu';

export const EmailFormMultiCreatableSelect = forwardRef<
  SelectInstance,
  FormSelectProps
>(({ name, formId, ...rest }, ref) => {
  const { getInputProps } = useField(name, formId);
  const { id, onChange, onBlur, value } = getInputProps();
  const [_, copyToClipboard] = useCopyToClipboard();

  const handleBlur = (stringVal: string) => {
    if (stringVal && emailRegex.test(stringVal)) {
      onBlur([...value, { label: stringVal, value: stringVal }]);

      return;
    }
    onBlur(value);
  };
  const Option = useCallback((rest: OptionProps<SelectOption>) => {
    return (
      <chakraComponents.Option {...rest}>
        {rest.data.label
          ? `${rest.data.label} - ${rest.data.value}`
          : rest.data.value}
        {rest?.isFocused && (
          <IconButton
            aria-label='Copy'
            size='xs'
            p={0}
            height={5}
            variant='ghost'
            icon={<Copy01 boxSize={3} color='gray.500' />}
            onClick={(e) => {
              e.stopPropagation();
              copyToClipboard(rest.data.value, 'Email copied!');
            }}
          />
        )}
      </chakraComponents.Option>
    );
  }, []);
  const MultiValue = useCallback(
    (rest: MultiValueProps<SelectOption>) => {
      return (
        <Menu isLazy closeOnSelect={false}>
          <MenuButton
            sx={{
              '&[aria-expanded="true"] > span > span': {
                bg: 'primary.50 !important',
                color: 'primary.700 !important',
                borderColor: 'primary.200 !important',
              },
            }}
          >
            <chakraComponents.MultiValue {...rest}>
              {rest.children}
            </chakraComponents.MultiValue>
          </MenuButton>
          <ChakraMenuList maxW={300}>
            {rest?.data?.value && (
              <MenuItem
                display='flex'
                justifyContent='space-between'
                onClick={(e) => {
                  e.stopPropagation();
                  copyToClipboard(rest?.data?.value, 'Email copied!');
                }}
              >
                {rest?.data?.value}
                <Copy01 boxSize={3} color='gray.500' ml={2} />
              </MenuItem>
            )}

            <MenuItem
              onClick={() => {
                const newValue = (
                  (rest?.selectProps?.value as Array<SelectOption>) ?? []
                )?.filter((e: SelectOption) => e.value !== rest?.data?.value);
                onChange(newValue);
              }}
            >
              Remove address
            </MenuItem>
          </ChakraMenuList>
        </Menu>
      );
    },
    [getInputProps],
  );

  const components = useMemo(
    () => ({
      MultiValueRemove: () => null,
    }),
    [],
  );

  return (
    <MultiCreatableSelect
      ref={ref}
      id={id}
      formId={formId}
      name={name}
      value={value}
      onBlur={(e) => handleBlur(e.target.value)}
      onChange={onChange}
      Option={Option}
      MultiValue={MultiValue}
      customStyles={multiCreatableSelectStyles}
      components={components}
      {...rest}
    />
  );
});

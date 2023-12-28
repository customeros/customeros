import { useField } from 'react-inverted-form';
import React, { forwardRef, useCallback } from 'react';

import { OptionProps, MultiValueProps } from 'chakra-react-select';
import {
  Menu,
  MenuItem,
  MenuButton,
  MenuList as ChakraMenuList,
} from '@chakra-ui/menu';

import { SelectOption } from '@ui/utils';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton';
import { chakraComponents } from '@ui/form/SyncSelect';
import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { emailRegex } from '@organization/src/components/Timeline/events/email/utils';
import {
  FormSelectProps,
  MultiCreatableSelect,
} from '@ui/form/MultiCreatableSelect';

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
    console.log('🏷️ ----- : HERE ');

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

  const MultiValue = useCallback((rest: MultiValueProps) => {
    return (
      <Menu isLazy closeOnSelect={false}>
        <MenuButton>
          <chakraComponents.MultiValue {...rest}>
            {rest.children}
          </chakraComponents.MultiValue>
        </MenuButton>
        <ChakraMenuList>
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
              <Copy01 boxSize={3} color='gray.500' />
            </MenuItem>
          )}

          <MenuItem
            onClick={(e) => {
              const newValue = value.filter(
                (e) => e.value !== rest?.data?.value,
              );
              onChange(newValue);
            }}
          >
            Remove address
          </MenuItem>
        </ChakraMenuList>
      </Menu>
    );
  }, []);

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
      {...rest}
    />
  );
});

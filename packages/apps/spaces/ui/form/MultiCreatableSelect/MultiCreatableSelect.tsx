import React, { useMemo, forwardRef, useCallback, ComponentType } from 'react';

import { Flex } from '@chakra-ui/react';
import {
  OptionProps,
  MenuListProps,
  ChakraStylesConfig,
} from 'chakra-react-select';

import { Tooltip } from '@ui/presentation/Tooltip';
import { SelectOption, chakraStyles } from '@ui/utils';
import { SelectInstance } from '@ui/form/SyncSelect/Select';
import { multiCreatableSelectStyles } from '@ui/form/MultiCreatableSelect/styles';
import {
  ControlProps,
  chakraComponents,
  AsyncCreatableProps,
  AsyncCreatableSelect,
  MultiValueGenericProps,
} from '@ui/form/SyncSelect';

// TODO: to be removed
export type CustomStylesFn = (
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  props: ChakraStylesConfig<any, any, any> | undefined,
) => chakraStyles;

// Exhaustively typing this Props interface does not offer any benefit at this moment
// TODO: Revisit this interface - naming is wrong and props need re-work
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface FormSelectProps extends AsyncCreatableProps<any, any, any> {
  name: string;
  formId: string;
  withTooltip?: boolean;
  // TODO: discard customStyles in favour of existing chakraStyles
  customStyles?: CustomStylesFn;
  optionAction?: (data: string) => JSX.Element;
  Option?: ComponentType<OptionProps<SelectOption>>;
}

export const MultiCreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ chakraStyles, ...props }, ref) => {
    const Control = useCallback(({ children, ...rest }: ControlProps) => {
      return (
        <chakraComponents.Control {...rest}>
          {children}
        </chakraComponents.Control>
      );
    }, []);
    const MultiValueLabel = useCallback(
      (rest: MultiValueGenericProps<SelectOption>) => {
        if (props?.withTooltip) {
          return (
            <chakraComponents.MultiValueLabel {...rest}>
              <Tooltip
                label={rest.data.label.length > 0 ? rest.data.value : ''}
                placement='top'
              >
                {rest.data.label || rest.data.value}
              </Tooltip>
            </chakraComponents.MultiValueLabel>
          );
        }

        return (
          <chakraComponents.MultiValueLabel {...rest}>
            {rest.data.label || rest.data.value}
          </chakraComponents.MultiValueLabel>
        );
      },
      [],
    );

    const Option = useCallback((rest: OptionProps<SelectOption>) => {
      return (
        <chakraComponents.Option {...rest}>
          <Flex
            flex={1}
            justifyContent='space-between'
            sx={{
              '& button': {
                display: 'none',
              },
            }}
            _hover={{
              '& button': {
                display: 'block',
              },
            }}
          >
            {rest.data.label || rest.data.value}
            {props?.optionAction && props.optionAction(rest.data.value)}
          </Flex>
        </chakraComponents.Option>
      );
    }, []);

    const MenuList = useCallback((rest: MenuListProps) => {
      return (
        <chakraComponents.MenuList {...rest}>
          {rest.children}
        </chakraComponents.MenuList>
      );
    }, []);

    const components = useMemo(
      () => ({
        Control,
        MultiValueLabel,
        MenuList,
        ClearIndicator: () => null,
        DropdownIndicator: () => null,
        Option: (props?.Option || Option) as ComponentType<OptionProps>,
      }),
      [Control, MultiValueLabel],
    );

    return (
      <AsyncCreatableSelect
        loadOptions={props?.loadOptions}
        variant='unstyled'
        focusBorderColor='transparent'
        ref={ref}
        components={components}
        tabSelectsValue={false}
        isMulti
        tagVariant='ghost'
        closeMenuOnSelect={false}
        chakraStyles={
          props?.customStyles?.(chakraStyles) ||
          multiCreatableSelectStyles(chakraStyles)
        }
        {...props}
      />
    );
  },
);

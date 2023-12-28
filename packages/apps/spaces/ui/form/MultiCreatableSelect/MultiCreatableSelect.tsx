import React, { useMemo, forwardRef, useCallback, ComponentType } from 'react';

import { Portal } from '@chakra-ui/react';
import {
  Menu,
  MenuItem,
  MenuButton,
  MenuList as ChakraMenuList,
} from '@chakra-ui/menu';
import {
  OptionProps,
  MenuListProps,
  MultiValueProps,
  ChakraStylesConfig,
} from 'chakra-react-select';

import { Button } from '@ui/form/Button';
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
  MultiValue?: ComponentType<MultiValueProps>;
  optionAction?: (data: string) => JSX.Element;
  Option?: ComponentType<OptionProps<SelectOption>>;
}

export const MultiCreatableSelect = forwardRef<SelectInstance, FormSelectProps>(
  ({ chakraStyles, ...props }, ref) => {
    const containerRef = React.useRef();

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
          {rest.data.label || rest.data.value}
          {props?.optionAction &&
            rest?.isFocused &&
            props.optionAction(rest.data.value)}
        </chakraComponents.Option>
      );
    }, []);

    const MenuList = useCallback((rest: MenuListProps) => {
      return (
        <chakraComponents.MenuList {...rest}>
          {rest.children}
          <Button>Remove option</Button>
        </chakraComponents.MenuList>
      );
    }, []);
    const MultiValue = useCallback((rest: MultiValueProps) => {
      return (
        <chakraComponents.MultiValue {...rest}>
          {rest.children}
        </chakraComponents.MultiValue>
      );
    }, []);

    const components = useMemo(
      () => ({
        Control,
        MultiValueLabel,
        MenuList,
        MultiValue: props?.MultiValue || MultiValue,
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

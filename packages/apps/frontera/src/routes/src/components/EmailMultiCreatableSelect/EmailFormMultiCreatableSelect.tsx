import { useMemo, useCallback } from 'react';
import {
  OptionProps,
  MultiValueProps,
  components as reactSelectComponents,
} from 'react-select';

import { observer } from 'mobx-react-lite';
import { EmailParticipantSelectUsecase } from '@domain/usecases/email-composer/email-participant-select.usecase.ts';

import { cn } from '@ui/utils/cn.ts';
import { SelectOption } from '@ui/utils/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Copy01 } from '@ui/media/icons/Copy01';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { getMultiValueLabelClassNames } from '@ui/form/CreatableSelect';
import { MultiValueWithActionMenu } from '@shared/components/EmailMultiCreatableSelect/MultiValueWithActionMenu';
import {
  Select,
  SelectProps,
  getMenuClassNames,
  getOptionClassNames,
  getMenuListClassNames,
  getContainerClassNames,
} from '@ui/form/Select';

interface EmailFormMultiCreatableSelectProps extends SelectProps {
  name: string;
  navigateAfterAddingToPeople: boolean;
  emailParticipantUseCase: EmailParticipantSelectUsecase;
}

export const EmailFormMultiCreatableSelect = observer(
  ({
    name,
    navigateAfterAddingToPeople,
    emailParticipantUseCase,
    ...props
  }: EmailFormMultiCreatableSelectProps) => {
    const [_, copyToClipboard] = useCopyToClipboard();

    const Option = useCallback((rest: OptionProps<SelectOption>) => {
      const noEmail = rest?.data?.label && !rest?.data?.value && (
        <p className='text-xs'>
          {rest.data.label} -
          <span className='text-gray-500 ml-1'>
            [No email for this contact]
          </span>
        </p>
      );

      return (
        <reactSelectComponents.Option {...rest}>
          {rest.data.label}

          {noEmail}
          {rest.data.value && (
            <>
              {rest.data.label.length > 0 && (
                <span className='text-gray-500 text-xs'> - </span>
              )}
              <span className='text-xs'> {rest.data.value}</span>
            </>
          )}
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
            name={name || ''}
            navigateAfterAddingToPeople={navigateAfterAddingToPeople}
            removeOption={emailParticipantUseCase.removeSelectedEmail}
          />
        );
      },
      [
        name,
        navigateAfterAddingToPeople,
        emailParticipantUseCase.removeSelectedEmail,
      ],
    );

    const components = useMemo(
      () => ({
        MultiValueRemove: () => null,
        LoadingIndicator: () => null,
        MultiValue,
        Option,
        ClearIndicator: () => null,
      }),
      [MultiValue, emailParticipantUseCase],
    );

    return (
      <Select
        isMulti
        autoFocus
        size='sm'
        name={name}
        backspaceRemovesValue
        menuPlacement={'auto'}
        components={components}
        filterOption={() => true}
        inputValue={emailParticipantUseCase.searchTerm}
        options={emailParticipantUseCase.emailOptionsList}
        menuIsOpen={emailParticipantUseCase.searchTerm.length > 0}
        styles={{ menuList: (base) => ({ ...base, maxHeight: '200px' }) }}
        onInputChange={(inputValue) => {
          emailParticipantUseCase.setSearchTerm(inputValue);
        }}
        value={emailParticipantUseCase.selectedEmails?.map(
          (email: SelectOption) => ({
            label: email.label || '',
            value: email.value || '',
          }),
        )}
        classNames={{
          input: () => 'pl-1',
          placeholder: () => 'pl-1 text-gray-400',
          container: ({ isFocused }) =>
            getContainerClassNames('flex flex-col min-h-[auto]', 'unstyled', {
              isFocused,
              size: 'sm',
            }),
          option: ({ isFocused }) =>
            getOptionClassNames('!cursor-pointer', { isFocused }),
          menuList: () => getMenuListClassNames(cn('w-full')),
          menu: ({ menuPlacement }) =>
            getMenuClassNames(menuPlacement)('bg-white', 'sm'),
          noOptionsMessage: () => 'text-gray-500',
          valueContainer: () => '!cursor-text mx-0 my-0.5',
          multiValueLabel: () => {
            return getMultiValueLabelClassNames(
              'bg-transparent hover:bg-transparent focus:bg-transparent active:bg-transparent rounded-sm px-1',
              'sm',
            );
          },
        }}
        {...props}
        closeMenuOnSelect={false}
        onChange={(e) => {
          emailParticipantUseCase.select(e);
        }}
        onKeyDown={(e) => {
          e.stopPropagation();

          if (e.key === 'Enter') {
            if (
              emailParticipantUseCase.searchTerm.length &&
              !emailParticipantUseCase.emailOptionsList.length
            ) {
              emailParticipantUseCase.addOption();
              e.preventDefault();
            }
          }
        }}
        noOptionsMessage={({ inputValue }) => (
          <div
            className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              emailParticipantUseCase.addOption();
            }}
          >
            <Plus />
            <span>{`Create "${inputValue}"`}</span>
          </div>
        )}
      />
    );
  },
);

import React, { ReactNode, ReactElement } from 'react';

import { observer } from 'mobx-react-lite';
import { FlowStepCommandMenuType } from '@store/UI/FlowStepCommandMenu.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

import { StepsHub } from './action';
import { TriggersHub, RecordAddedManually } from './trigger';

export const CommandsContainer = ({
  children,
  placeholder,
  dataTest,
}: {
  dataTest?: string;
  placeholder: string;
  children: ReactNode;
}) => {
  return (
    <Command
      data-test={dataTest}
      onClick={(e) => {
        e.stopPropagation();
      }}
      filter={(value, search, keywords) => {
        const extendValue = value.replace(/\s/g, '') + keywords;
        const searchWithoutSpaces = search.replace(/\s/g, '');

        if (
          extendValue.toLowerCase().includes(searchWithoutSpaces.toLowerCase())
        )
          return 1;

        return 0;
      }}
    >
      <CommandInput
        autoFocus
        className='p-1 text-sm'
        placeholder={placeholder}
        inputWrapperClassName='min-h-4'
        data-test={`${dataTest}-input`}
        wrapperClassName='py-2 px-4 mt-2'
      />
      <Command.List>
        <Command.Group>{children}</Command.Group>
      </Command.List>
    </Command>
  );
};

const Commands: Record<FlowStepCommandMenuType, ReactElement> = {
  StepsHub: <StepsHub />,
  Webhook: <div />, // todo
  TriggersHub: <TriggersHub />,
  RecordCreated: <div />, // todo
  RecordUpdated: <div />, // todo
  RecordAddedManually: <RecordAddedManually />,
  RecordMatchesCondition: <div />, // todo
};

const placeholderMap: Record<FlowStepCommandMenuType, string> = {
  StepsHub: 'Search a step',
  TriggersHub: 'Search a trigger',
  RecordAddedManually: 'Search a record',
  Webhook: '',
  RecordCreated: '',
  RecordUpdated: '',
  RecordMatchesCondition: '',
};

export const DropdownCommandMenu = observer(() => {
  const { ui } = useStore();

  return (
    <CommandsContainer
      dataTest={ui.flowCommandMenu.type}
      placeholder={placeholderMap[ui.flowCommandMenu.type]}
    >
      {Commands[ui.flowCommandMenu.type]}
    </CommandsContainer>
  );
});

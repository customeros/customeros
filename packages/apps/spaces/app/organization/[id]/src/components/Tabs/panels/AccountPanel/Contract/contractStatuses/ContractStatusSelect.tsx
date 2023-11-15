import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Select } from '@ui/form/SyncSelect';
import { ContractStatus } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { contractButtonSelect } from '@organization/src/components/Tabs/shared/contractSelectStyles';

import { contractOptionIcon } from './utils';

interface ContractStatusSelectProps {}

export const contractStatusOptions: SelectOption<ContractStatus>[] = [
  { label: 'Draft', value: ContractStatus.Draft },
  { label: 'Ended', value: ContractStatus.Ended },
  { label: 'Live', value: ContractStatus.Live },
];

export const ContractStatusSelect: React.FC<ContractStatusSelectProps> = () => {
  return (
    <Select
      isSearchable={false}
      isClearable={false}
      isMulti={false}
      placeholder='Status'
      options={contractStatusOptions}
      formatOptionLabel={(
        data: SelectOption<ContractStatus>,
        formatOptionLabelMeta,
      ) => {
        const icon = contractOptionIcon?.[data?.value];
        const isButton =
          formatOptionLabelMeta.selectValue?.[0]?.value === data.value &&
          formatOptionLabelMeta.context === 'value';

        return (
          <Flex alignItems='center' gap={isButton ? 1 : 2}>
            {icon && (
              <Flex alignItems='center' boxSize={isButton ? 3 : 4}>
                {icon}
              </Flex>
            )}
            <Text
              color={
                isButton && data.value === ContractStatus.Live
                  ? 'primary.800'
                  : 'gray.800'
              }
            >
              {data.label}
            </Text>
          </Flex>
        );
      }}
      chakraStyles={{
        ...contractButtonSelect,

        container: (props, state) => {
          const isLive = state.getValue()[0]?.value === ContractStatus.Live;

          return {
            ...props,
            px: 2,
            py: '1px',
            border: '1px solid',
            borderColor: isLive ? 'primary.200' : 'gray.300',
            backgroundColor: isLive ? 'primary.50' : 'transparent',
            color: isLive ? 'primary.700' : 'gray.500',

            borderRadius: 'md',
            fontSize: 'xs',
            maxHeight: '22px',

            '& > div': {
              p: 0,
              border: 'none',
              fontSize: 'xs',
              maxHeight: '22px',
              minH: 'auto',
            },
          };
        },
        valueContainer: (props, state) => {
          const isLive = state.getValue()[0]?.value === ContractStatus.Live;

          return {
            ...props,
            p: 0,
            border: 'none',
            fontSize: 'xs',
            maxHeight: '22px',
            minH: 'auto',
            color: isLive ? 'primary.700' : 'gray.500',
          };
        },

        menuList: (props) => {
          return {
            ...props,
            w: 'fit-content',
            minWidth: '125px',
            right: '60px',
          };
        },
      }}
    />
  );
};

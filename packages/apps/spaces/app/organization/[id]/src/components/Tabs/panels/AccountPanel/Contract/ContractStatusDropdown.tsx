'use client';

import React, { useState } from 'react';

import { Box } from '@ui/layout/Box';
import { Select } from '@ui/form/SyncSelect';
import { Check } from '@ui/media/icons/Check';
import { XClose } from '@ui/media/icons/XClose';
import { DotLive } from '@ui/media/icons/DotLive';

export const ContractStatusDropdown = () => {
  const [value, setValue] = useState({ label: 'Draft', value: 'draft' });

  return (
    <Box>
      <Select
        isSearchable={false}
        isClearable={false}
        isMulti={false}
        value={value}
        onChange={setValue}
        options={[
          {
            label: 'Draft',
            value: 'draft',
          },
          {
            label: 'Live',
            value: 'live',
          },
          {
            label: 'Ended',
            value: 'ended',
          },
        ]}
        formatOptionLabel={(value) => (
          <>
            {value.value === 'draft' ? (
              <Check />
            ) : value.value === 'live' ? (
              <DotLive />
            ) : (
              <XClose />
            )}
            {value.label}
          </>
        )}
        chakraStyles={{
          container: (props, state) => {
            const isCustomer = state.getValue()[0]?.value === 'live';

            return {
              ...props,
              px: 2,
              py: '1px',
              border: '1px solid',
              borderColor: isCustomer ? 'success.200' : 'gray.300',
              backgroundColor: isCustomer ? 'primary.50' : 'transparent',
              color: isCustomer ? 'primary.700' : 'gray.500',

              borderRadius: '2xl',
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
            const isCustomer = state.getValue()[0]?.value === 'live';

            return {
              ...props,
              p: 0,
              border: 'none',
              fontSize: 'xs',
              maxHeight: '22px',
              minH: 'auto',
              color: isCustomer ? 'primary.700' : 'gray.500',
            };
          },
          singleValue: (props) => {
            return {
              ...props,
              maxHeight: '22px',
              p: 0,
              minH: 'auto',
              color: 'inherit',
            };
          },
          input: (props) => {
            return {
              ...props,
              maxHeight: '22px',
              minH: 'auto',
              p: 0,
            };
          },
          inputContainer: (props) => {
            return {
              ...props,
              maxHeight: '22px',
              minH: 'auto',
              p: 0,
            };
          },

          control: (props) => {
            return {
              ...props,
              w: '100%',
              border: 'none',
            };
          },

          menuList: (props) => {
            return {
              ...props,
              w: 'fit-content',
              left: '-32px',
            };
          },
        }}
        leftElement={
          value.value === 'draft' ? (
            <Check />
          ) : value.value === 'live' ? (
            <DotLive />
          ) : (
            <XClose />
          )
        }
      />
    </Box>
  );
};

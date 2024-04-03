import React, { useState, useEffect } from 'react';

import { cn } from '@ui/utils/cn';
import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { InputProps } from '@ui/form/Input';
import { useOutsideClick } from '@ui/utils';
import { Button } from '@ui/form/Button/Button';
import { Divider } from '@ui/presentation/Divider/Divider';
import { EmailSelect } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/EmailsInputGroup/EmailSelect';

interface EmailsInputGroupProps extends InputProps {
  formId: string;
  modal?: boolean;
  cc: Array<{ label: string; value: string }>;

  to?: { label: string; value: string } | null;
  bcc: Array<{ label: string; value: string }>;
}

export const EmailsInputGroup = ({
  to,
  cc = [],
  bcc = [],
  formId,
}: EmailsInputGroupProps) => {
  const [showCC, setShowCC] = useState(false);
  const [showBCC, setShowBCC] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const [focusedItemIndex, setFocusedItemIndex] = useState<false | number>(
    false,
  );
  const ref = React.useRef(null);
  useOutsideClick({
    ref: ref,
    handler: () => {
      setIsFocused(false);
      setFocusedItemIndex(false);
      setShowCC(false);
      setShowBCC(false);
    },
  });

  const handleFocus = (index: number) => {
    setIsFocused(true);
    setFocusedItemIndex(index);
  };

  useEffect(() => {
    if (showCC && !isFocused) {
      handleFocus(1);
    }
  }, [showCC]);

  useEffect(() => {
    if (showBCC && !isFocused) {
      handleFocus(2);
    }
  }, [showBCC]);

  return (
    <div ref={ref}>
      <div className='flex relative items-center h-8'>
        <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
          Send invoice
        </p>
        <Divider />

        <div className='flex'>
          {!showCC && !cc.length && (
            <Button
              variant='ghost'
              color='gray.400'
              size='sm'
              className='text-sm px-1 mx-1'
              onClick={() => {
                setShowCC(true);
                setFocusedItemIndex(1);
              }}
            >
              CC
            </Button>
          )}

          {!showBCC && !bcc.length && (
            <Button
              variant='ghost'
              size='sm'
              className='text-sm px-1 '
              color='gray.400'
              onClick={() => {
                setShowBCC(true);
                setFocusedItemIndex(2);
              }}
            >
              BCC
            </Button>
          )}
        </div>
      </div>

      <Box width='100%'>
        {isFocused && (
          <>
            <EmailSelect
              formId={formId}
              fieldName='billingEmail'
              entryType='To'
              autofocus={focusedItemIndex === 0}
            />
            {(showCC || !!cc.length) && (
              <EmailSelect
                formId={formId}
                fieldName='billingEmailCC'
                entryType='CC'
                autofocus={focusedItemIndex === 1}
              />
            )}
            {(showBCC || !!bcc.length) && (
              <EmailSelect
                formId={formId}
                fieldName='billingEmailBCC'
                entryType='BCC'
                autofocus={focusedItemIndex === 2}
              />
            )}
          </>
        )}
      </Box>

      {!isFocused && (
        <Flex flexDir='column' flex={1}>
          <div
            onClick={() => handleFocus(0)}
            role='button'
            aria-label='Click to input participant data'
            className={cn('overflow-hidden', {
              'flex-1': !bcc.length,
            })}
          >
            <span className='text-sm font-semibold text-gray-700 mr-1'>To</span>
            <Text color='gray.500' noOfLines={1}>
              <>{to?.value ? to.value : `⚠️ [invalid email]`}</>
            </Text>
          </div>

          {!!cc.length && (
            <div
              onClick={() => handleFocus(1)}
              role='button'
              aria-label='Click to input participant data'
              className={cn('overflow-hidden', {
                'flex-1': !bcc.length,
              })}
            >
              <span className='text-sm font-semibold text-gray-700 mr-1'>
                CC
              </span>
              <p className='text-gray-500 whitespace-nowrap overflow-ellipsis overflow-hidden'>
                {[...cc].map((email) => email.value).join(', ')}
              </p>
            </div>
          )}
          {!!bcc.length && (
            <div
              onClick={() => handleFocus(2)}
              role='button'
              className={cn('overflow-hidden', {
                'flex-1': !bcc.length,
              })}
              aria-label='Click to input participant data'
            >
              <span className='text-sm font-semibold text-gray-700 mr-1'>
                BCC
              </span>
              <p className='text-gray-500 whitespace-nowrap overflow-ellipsis overflow-hidden'>
                {[...bcc].map((email) => email.value).join(', ')}
              </p>
            </div>
          )}
        </Flex>
      )}
    </div>
  );
};

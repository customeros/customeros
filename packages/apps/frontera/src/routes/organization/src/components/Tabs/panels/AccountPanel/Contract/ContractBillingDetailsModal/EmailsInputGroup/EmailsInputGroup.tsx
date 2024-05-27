import React, { useRef, useMemo, useState, useEffect } from 'react';

import { cn } from '@ui/utils/cn';
import { InputProps } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { FormInput } from '@ui/form/Input/FormInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { SelectOption } from '@shared/types/SelectOptions';
import { Divider } from '@ui/presentation/Divider/Divider';
import { validateEmail } from '@shared/util/emailValidation';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { EmailSelect } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/EmailsInputGroup/EmailSelect';

interface EmailsInputGroupProps extends InputProps {
  formId: string;
  modal?: boolean;
  to?: string | null;
  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;
}

const EmailList = ({
  emailList,
}: {
  emailList: Array<SelectOption<string>>;
}) => {
  return (
    <p className='text-gray-500 whitespace-nowrap overflow-ellipsis overflow-hidden h-8'>
      {[...emailList].map((email, i) => {
        const validationMessage = validateEmail(email.value);

        return (
          <React.Fragment key={email.value}>
            <Tooltip label={validationMessage || ''}>
              <span
                className={cn('mr-1 text-base', {
                  'text-warning-700': validateEmail(email.value),
                })}
              >
                {email.value}
                {i < emailList.length - 1 && ','}
              </span>
            </Tooltip>
          </React.Fragment>
        );
      })}
    </p>
  );
};
const ToEmailInput = ({
  formId,
  email,
}: {
  formId: string;
  email?: string | null;
}) => {
  const ref = useRef<HTMLInputElement>(null);

  const validationMessage = useMemo(() => {
    if (!email) return '';

    return validateEmail(email) ?? '';
  }, [email]);

  useEffect(() => {
    ref.current?.focus();
  }, [validationMessage]);

  return (
    <Tooltip label={validationMessage} align='start' side={'bottom'}>
      <FormInput
        ref={ref}
        variant='unstyled'
        className={cn('text-warning-700' && validationMessage)}
        formId={formId}
        autoComplete='off'
        label='To'
        labelProps={{
          className: 'text-sm mb-0  font-semibold inline-block pt-0',
        }}
        size='sm'
        name='billingEmail'
        placeholder='To email address'
      />
    </Tooltip>
  );
};

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
      <div className='flex relative items-center h-8  mt-2'>
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
      <ToEmailInput formId={formId} email={to} />

      <div className='flex-col flex-1 w-full gap-4'>
        {isFocused && (
          <>
            {(showCC || !!cc.length) && (
              <EmailSelect
                formId={formId}
                fieldName='billingEmailCC'
                entryType='CC'
                placeholder='CC email addresses'
                autofocus={focusedItemIndex === 1}
              />
            )}
            {(showBCC || !!bcc.length) && (
              <EmailSelect
                formId={formId}
                fieldName='billingEmailBCC'
                placeholder='BCC email addresses'
                entryType='BCC'
                autofocus={focusedItemIndex === 2}
              />
            )}
          </>
        )}
      </div>

      {!isFocused && (
        <div className='flex-col flex-1 gap-4'>
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
              <EmailList emailList={cc} />
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
              <EmailList emailList={bcc} />
            </div>
          )}
        </div>
      )}
    </div>
  );
};

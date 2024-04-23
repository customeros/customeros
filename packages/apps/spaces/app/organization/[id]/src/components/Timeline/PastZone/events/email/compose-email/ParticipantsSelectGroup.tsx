import Image from 'next/image';
import React, { useState, useEffect } from 'react';

import { cn } from '@ui/utils/cn';
import { InputProps } from '@ui/form/Input';
import { useOutsideClick } from '@ui/utils';
import { Button } from '@ui/form/Button/Button';
import { EmailSubjectInput } from '@organization/src/components/Timeline/PastZone/events/email/compose-email/EmailSubjectInput';
import { EmailParticipantSelect } from '@organization/src/components/Timeline/PastZone/events/email/compose-email/EmailParticipantSelect';

interface ParticipantSelectGroupGroupProps extends InputProps {
  formId: string;
  modal?: boolean;
  to: Array<{ label: string; value: string }>;

  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;
}

export const ParticipantsSelectGroup = ({
  to = [],
  cc = [],
  bcc = [],
  modal,
  formId,
}: ParticipantSelectGroupGroupProps) => {
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
    <div className='flex justify-between mt-3' ref={ref}>
      <div className='w-[100%]'>
        {isFocused && (
          <>
            <EmailParticipantSelect
              formId={formId}
              fieldName='to'
              entryType='To'
              autofocus={focusedItemIndex === 0}
            />
            {(showCC || !!cc.length) && (
              <EmailParticipantSelect
                formId={formId}
                fieldName='cc'
                entryType='CC'
                autofocus={focusedItemIndex === 1}
              />
            )}
            {(showBCC || !!bcc.length) && (
              <EmailParticipantSelect
                formId={formId}
                fieldName='bcc'
                entryType='BCC'
                autofocus={focusedItemIndex === 2}
              />
            )}
          </>
        )}

        {!isFocused && (
          <div
            className={cn(isFocused ? 'flex-1' : 'unset', 'flex mt-1 flex-col')}
          >
            <div
              className={cn(
                !cc.length && !bcc.length ? 'flex-1' : 'unset',
                'flex',
              )}
              onClick={() => handleFocus(0)}
              role='button'
              aria-label='Click to input participant data'
            >
              <span className='text-gray-700 font-semibold mr-1'>To:</span>
              <span className='text-gray-500 line-clamp-1'>
                {!!to?.length && (
                  <>
                    {to
                      ?.map((email) =>
                        email?.value
                          ? email.value
                          : `⚠️ ${email.label} [invalid email]`,
                      )
                      .join(', ')}
                  </>
                )}
              </span>
            </div>

            {!!cc.length && (
              <div
                className={cn(!bcc.length ? 'flex-1' : 'unset', 'flex')}
                onClick={() => handleFocus(1)}
                onFocusCapture={() => handleFocus(1)}
                role='button'
                aria-label='Click to input participant data'
              >
                <span className='text-gray-700 font-semibold mr-1'>CC:</span>
                <p className='text-gray-500 line-clamp-1'>
                  {[...cc].map((email) => email.value).join(', ')}
                </p>
              </div>
            )}
            {!!bcc.length && (
              <div
                className='flex'
                onClick={() => handleFocus(2)}
                onFocusCapture={() => handleFocus(2)}
                role='button'
                aria-label='Click to input participant data'
              >
                <span className='text-gray-700 font-semibold mr-1'>BCC:</span>
                <p className='text-gray-500 line-clamp-1'>
                  {[...bcc].map((email) => email.value).join(', ')}
                </p>
              </div>
            )}
          </div>
        )}
        <EmailSubjectInput formId={formId} fieldName='subject' />
      </div>
      <div className='flex max-w-[64px] mr-4 items-start'>
        {!showCC && (
          <Button
            className='text-gray-400 font-semibold px-1'
            variant='ghost'
            size='sm'
            onClick={() => {
              setShowCC(true);
              setFocusedItemIndex(1);
            }}
          >
            CC
          </Button>
        )}

        {!showBCC && (
          <Button
            className='text-gray-400 font-semibold px-1'
            variant='ghost'
            size='sm'
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

      {!modal && (
        <div>
          <Image
            src={'/backgrounds/organization/post-stamp.webp'}
            alt='Email'
            width={54}
            height={70}
          />
        </div>
      )}
    </div>
  );
};

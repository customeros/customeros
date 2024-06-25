import { useField } from 'react-inverted-form';
import { useSearchParams } from 'react-router-dom';
import React, { useState, useEffect } from 'react';

import postStamp from '@assets/backgrounds/organization/post-stamp.webp';

import { cn } from '@ui/utils/cn';
import { InputProps } from '@ui/form/Input';
import { Google } from '@ui/media/logos/Google';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Microsoft } from '@ui/media/icons/Microsoft';
import { FormSelect } from '@ui/form/Select/FormSelect';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { EmailSubjectInput } from '@organization/components/Timeline/PastZone/events/email/compose-email/EmailSubjectInput';
import { EmailParticipantSelect } from '@organization/components/Timeline/PastZone/events/email/compose-email/EmailParticipantSelect';

interface ParticipantSelectGroupGroupProps extends InputProps {
  formId: string;
  modal?: boolean;
  attendees: Array<string>;
  to: Array<{ label: string; value: string }>;

  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;
}

export const ParticipantsSelectGroup = ({
  attendees = [],
  to = [],
  cc = [],
  bcc = [],
  modal,
  formId,
}: ParticipantSelectGroupGroupProps) => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const params = new URLSearchParams(searchParams?.toString() ?? '');
  const id = params.get('events') ?? undefined;

  const client = getGraphQLClient();
  const { data: globalCacheData } = useGlobalCacheQuery(client);

  const { getInputProps: fromGetInputProps } = useField('from', formId);
  const { getInputProps: fromProviderGetInputProps } = useField(
    'fromProvider',
    formId,
  );
  const { onChange: fromOnChange } = fromGetInputProps();
  const { onChange: fromProviderOnChange } = fromProviderGetInputProps();

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

  const [fromOptions, setFromOptions] = useState(
    [] as Array<{
      label: string;
      value: string;
      active: boolean;
      provider: string;
    }>,
  );

  useEffect(() => {
    if (globalCacheData) {
      const options = [] as Array<{
        label: string;
        value: string;
        active: boolean;
        provider: string;
      }>;

      globalCacheData?.global_Cache?.activeEmailTokens
        .filter((a) => (id && attendees.indexOf(a.email) > -1) || !id)
        .forEach((v) => {
          options.push({
            label: v.email,
            value: v.email,
            provider: v.provider,
            active: true,
          });
        });

      globalCacheData?.global_Cache?.inactiveEmailTokens
        .filter((a) => (id && attendees.indexOf(a.email) > -1) || !id)
        .forEach((v) => {
          options.push({
            label: v.email,
            value: v.email,
            provider: v.provider,
            active: false,
          });
        });

      setFromOptions(options);
    }
  }, [
    globalCacheData?.global_Cache?.activeEmailTokens,
    globalCacheData?.global_Cache?.inactiveEmailTokens,
    id,
    attendees,
  ]);

  useEffect(() => {
    if (fromOptions && fromOptions.length > 0) {
      const activeOption = fromOptions.filter(
        (v) => v.value === store.session.value.profile.email && v.active,
      );

      if (activeOption && activeOption.length > 0) {
        fromOnChange(activeOption[0]);
        fromProviderOnChange(activeOption[0].provider);
      } else {
        const firstActive = fromOptions.filter((v) => v.active);
        if (firstActive && firstActive.length > 0) {
          fromOnChange(firstActive[0]);
          fromProviderOnChange(firstActive[0].provider);
        }
      }
    }
  }, [fromOptions]);

  return (
    <div className='flex justify-between mt-3' ref={ref}>
      <div className='w-[100%]'>
        <div className='flex items-baseline mb-[-1px] mt-0 flex-1 overflow-visible'>
          <span className='text-gray-700 font-semibold mr-1'>From:</span>
          <FormSelect
            formId={formId}
            name='from'
            options={fromOptions}
            getOptionLabel={(props) => {
              const { value } = props;

              const activeOption = (fromOptions ?? []).find(
                (v) => v.value === value,
              );

              return (
                <div
                  className={
                    'flex items-center gap-2 justify-between w-100 ' +
                    (activeOption?.active === false ? 'opacity-50' : '')
                  }
                >
                  <div className={'flex'}>
                    {activeOption?.provider === 'google' ? (
                      <Google className='size-5 mr-2' />
                    ) : (
                      <Microsoft className='size-5 mr-2' />
                    )}
                    <span>{activeOption?.label}</span>
                  </div>
                  <div className={'flex'}>
                    {activeOption?.active === false && (
                      <span className='text-red-500'>Expired</span>
                    )}
                  </div>
                </div>
              ) as unknown as string;
            }}
            isOptionDisabled={(option) => !option.active}
          />
        </div>
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
          <img src={postStamp} alt='Email' width={54} height={70} />
        </div>
      )}
    </div>
  );
};

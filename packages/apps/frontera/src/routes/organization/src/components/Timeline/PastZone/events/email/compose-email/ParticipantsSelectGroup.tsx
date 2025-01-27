import { useSearchParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { Button } from '@ui/form/Button/Button';
import { Input, InputProps } from '@ui/form/Input';
import { Select, getContainerClassNames } from '@ui/form/Select';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';

import { EmailParticipantSelect } from './EmailParticipantSelect';

import postStamp from '/backgrounds/organization/post-stamp.webp';

interface ParticipantSelectGroupGroupProps extends InputProps {
  modal?: boolean;
  emailUseCase: TimelineEmailUsecase;
}

export const ParticipantsSelectGroup = observer(
  ({ modal, emailUseCase }: ParticipantSelectGroupGroupProps) => {
    const [searchParams] = useSearchParams();
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    const id = params.get('events') ?? undefined;

    const from = emailUseCase.fromSelector;

    const [showCC, setShowCC] = useState(false);
    const [showBCC, setShowBCC] = useState(false);
    const [isFocused, setIsFocused] = useState(false);
    const [focusedItemIndex, setFocusedItemIndex] = useState<false | number>(
      false,
    );
    const ref = useRef(null);

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
      <div
        ref={ref}
        className='flex justify-between mt-3'
        onMouseDown={(e) => {
          e.stopPropagation();
        }}
      >
        <div className='w-[100%]'>
          <div className='flex items-baseline mb-[-1px] mt-0 flex-1 overflow-visible'>
            <span className='text-gray-700 font-semibold mr-1 text-sm'>
              From:
            </span>

            <div className='w-full'>
              <Select
                id={id}
                ref={ref}
                size='xs'
                name='from'
                menuPlacement={'auto'}
                value={from.selectedEmail}
                options={from.emailOptions}
                onChange={(e) => from.select(e)}
                placeholder={'Enter name or email...'}
                isOptionDisabled={(option) => !option.active}
                classNames={{
                  container: () =>
                    getContainerClassNames(
                      'focus-within:border-transparent focus-within:hover:border-transparent hover:border-transparent focus:border-transparent hover:focus:border-transparent focus-within:border-transparent',
                      'flushed',
                      {
                        size: 'xs',
                      },
                    ),
                  input: () =>
                    'hover:border-transparent focus:border-transparent hover:focus:border-transparent focus-within:border-transparent',
                }}
                getOptionLabel={(props) => {
                  const { value } = props;

                  const activeOption = (from.emailOptions ?? []).find(
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
              />
            </div>
          </div>

          <EmailParticipantSelect
            entryType='To'
            autofocus={focusedItemIndex === 0}
            emailParticipantUseCase={emailUseCase.toSelector}
          />
          <>
            {(showCC || !!emailUseCase.ccSelector.selectedEmails.length) && (
              <EmailParticipantSelect
                entryType='CC'
                autofocus={focusedItemIndex === 1}
                emailParticipantUseCase={emailUseCase.ccSelector}
              />
            )}
            {(showBCC || !!emailUseCase.bccSelector.selectedEmails.length) && (
              <EmailParticipantSelect
                entryType='BCC'
                autofocus={focusedItemIndex === 2}
                emailParticipantUseCase={emailUseCase.bccSelector}
              />
            )}
          </>

          <div className='flex items-center flex-1'>
            <span className='text-gray-700 font-semibold mr-1 text-sm'>
              Subject:
            </span>
            <Input
              autoFocus
              size='xs'
              variant='unstyled'
              value={emailUseCase.subject}
              placeholder={'Enter subject...'}
              className='text-gray-500 height-[5px] text-sm'
              onChange={(e) => emailUseCase.updateSubject(e.target.value)}
            />
          </div>
        </div>
        <div className='flex max-w-[64px] mr-4 items-start'>
          {!showCC && !emailUseCase.ccSelector.selectedEmails.length && (
            <Button
              size='sm'
              variant='ghost'
              className='text-gray-400 font-semibold px-1'
              onClick={() => {
                setShowCC(true);
                setFocusedItemIndex(1);
              }}
            >
              CC
            </Button>
          )}

          {!showBCC && !emailUseCase.bccSelector.selectedEmails.length && (
            <Button
              size='sm'
              variant='ghost'
              color='gray.400'
              className='text-gray-400 font-semibold px-1'
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
          <div className='min-w-[54px] max-w-[54px]'>
            <img width={54} alt='Email' height={70} src={postStamp} />
          </div>
        )}
      </div>
    );
  },
);

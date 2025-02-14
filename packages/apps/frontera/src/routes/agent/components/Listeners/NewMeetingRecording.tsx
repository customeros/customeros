import { useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { NewMeetingRecordingUsecase } from '@domain/usecases/agents/listeners/new-meeting-recording.usecase';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { Combobox } from '@ui/form/Combobox';
import { Image } from '@ui/media/Image/Image';
import { Spinner } from '@ui/feedback/Spinner';
import { Logo, LogoName } from '@ui/media/Logo';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { getOptionClassNames } from '@ui/form/Select';
import logoCustomerOs from '@shared/assets/customer-os-small.png';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

const usecase = new NewMeetingRecordingUsecase();

export const NewMeetingRecording = observer(() => {
  const store = useStore();
  const { id } = useParams<{ id: string }>();

  useEffect(() => {
    usecase.init(id!);

    return () => {
      usecase.reset();
    };
  }, [id]);

  return (
    <div>
      <h2 className='text-sm font-medium mb-4'>Detect new meeting recording</h2>

      {usecase.listenerErrors && (
        <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] mb-4'>
          <Icon stroke='none' className='mr-2' name='dot-single' />
          <span className='text-sm'>{usecase.listenerErrors}</span>
        </div>
      )}

      <div className='flex flex-col gap-1 mb-4'>
        <p className='text-sm font-medium'>Notetaker</p>
        <p className='text-sm mb-1'>Which app should take the meeting notes?</p>

        <Popover open={usecase.isOpen} onOpenChange={() => usecase.toggle()}>
          <PopoverTrigger className={'flex items-center w-full'}>
            <div className='flex w-[320px] items-center bg-white border border-grayModern-200 px-2 py-1 rounded-md min-h-8'>
              {!usecase.notetaker ? (
                <Image
                  width={16}
                  height={16}
                  alt='CustomerOS'
                  src={logoCustomerOs}
                  className='logo-image rounded mr-2'
                />
              ) : (
                <Logo
                  className='mr-2'
                  name={usecase.notetaker.value as LogoName}
                />
              )}
              {usecase.notetaker ? (
                <div className='text-sm'>{usecase.notetaker.label}</div>
              ) : (
                <span className='text-gray-400 text-sm'>Select option...</span>
              )}
            </div>
          </PopoverTrigger>
          <PopoverContent align='start' className='min-w-[264px] w-[320px]'>
            <Combobox
              isSearchable={false}
              value={usecase.notetaker}
              placeholder='Search notetaker apps'
              options={NewMeetingRecordingUsecase.notetakerOptions}
              onChange={(option) => {
                usecase.setNotetaker(option.value);
                usecase.toggle();
                usecase.execute();
              }}
              noOptionsMessage={({ inputValue }) => (
                <div className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'>
                  <span>{`No results matching "${inputValue}"`}</span>
                </div>
              )}
              components={{
                Option: ({ children, innerProps, ...props }) => (
                  <div
                    {...innerProps}
                    className={getOptionClassNames('', props)}
                  >
                    <Logo
                      className='mr-2'
                      name={
                        NewMeetingRecordingUsecase.notetakerOptions.find(
                          (option) => option.label === children,
                        )?.value as LogoName
                      }
                    />
                    {children}
                  </div>
                ),
              }}
            />
          </PopoverContent>
        </Popover>
      </div>

      {usecase.notetaker && (
        <div className='flex flex-col gap-1'>
          <p className='text-sm font-medium'>Webhook URL</p>
          <p className='text-sm mb-1'>
            {`Use this webhook in Zapier to start listening for new ${usecase.notetaker?.label} meetings`}
          </p>

          <div
            className={cn(
              'flex items-center justify-between gap-1 px-2 py-1 border border-grayModern-200 rounded-md max-w-[320px] min-h-8 bg-white transition-all duration-200',
              usecase.isWebhookUrlLoading && 'opacity-50 cursor-not-allowed',
            )}
          >
            <p
              className={cn(
                'text-sm truncate',
                !usecase.webhookUrl && 'text-gray-400',
                usecase.isWebhookUrlLoading && 'text-gray-400',
              )}
            >
              {usecase.isWebhookUrlLoading
                ? 'Loading...'
                : usecase.webhookUrl || 'Webhook URL'}
            </p>
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='Copy webhook URL'
              icon={<Icon name='copy-03' />}
              isLoading={usecase.isWebhookUrlLoading}
              spinner={
                <Spinner
                  size='xs'
                  label='Retrieving webhook URL'
                  className='text-grayModern-500 fill-grayModern-100'
                />
              }
              onClick={() => {
                navigator.clipboard.writeText(usecase.webhookUrl);
                store.ui.toastSuccess(
                  'Webhook URL copied to clipboard',
                  'webhook-url-copied',
                );
              }}
            />
          </div>
        </div>
      )}
    </div>
  );
});

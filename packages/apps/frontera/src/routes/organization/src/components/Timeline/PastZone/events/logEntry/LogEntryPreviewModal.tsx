import React from 'react';

import copy from 'copy-to-clipboard';
import noteImg from '@assets/images/note-img-preview.png';

import { User } from '@graphql/types';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetTagsQuery } from '@organization/graphql/getTags.generated';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';
import { useLogEntryUpdateContext } from '@organization/components/Timeline/PastZone/events/logEntry/context/LogEntryUpdateModalContext';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

import { PreviewEditor } from './preview/PreviewEditor';
import { PreviewTags } from './preview/tags/PreviewTags';
import { LogEntryDatePicker } from './preview/LogEntryDatePicker';
import { LogEntryExternalLink } from './preview/LogEntryExternalLink';

const getAuthor = (user: User) => {
  if (!user?.firstName && !user?.lastName) {
    return 'Unknown';
  }

  return `${user.firstName} ${user.lastName}`.trim();
};

export const LogEntryPreviewModal: React.FC = () => {
  const { modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();
  const store = useStore();

  const event = modalContent as LogEntryWithAliases;
  const author = getAuthor(event?.logEntryCreatedBy);
  const authorEmail = event?.logEntryCreatedBy?.emails?.[0]?.email;
  const client = getGraphQLClient();
  const { data } = useGetTagsQuery(client);
  const isAuthor =
    !!event.logEntryCreatedBy &&
    event.logEntryCreatedBy?.emails?.findIndex(
      (e) => store.session.value.profile.email === e.email,
    ) !== -1;
  const { formId } = useLogEntryUpdateContext();

  if (!event.content) return null;

  return (
    <>
      <div className='py-4 px-6 pb-1 sticky top-0 rounded-xl'>
        <div className='flex justify-between items-center'>
          <div className='flex items-center'>
            <h2 className='text-lg font-semibold'>Log entry</h2>
          </div>
          <div className='flex justify-end items-center'>
            <Tooltip label='Copy link' side='bottom' asChild={false}>
              <div>
                <IconButton
                  className='text-sm text-gray-500 mr-1'
                  variant='ghost'
                  aria-label='Copy link to this entry'
                  size='xs'
                  icon={<Link03 className='text-gray-500' />}
                  onClick={() => copy(window.location.href)}
                />
              </div>
            </Tooltip>
            <Tooltip
              label='Close'
              aria-label='close'
              side='bottom'
              asChild={false}
            >
              <div>
                <IconButton
                  className='text-sm text-gray-500'
                  variant='ghost'
                  aria-label='Close preview'
                  color='gray.500'
                  size='xs'
                  icon={<XClose className='text-gray-500 size-5' />}
                  onClick={closeModal}
                />
              </div>
            </Tooltip>
          </div>
        </div>
      </div>
      <div className='mt-0 p-6 pt-0 overflow-auto max-h-[calc(100vh-9rem)]'>
        <div className='relative'>
          <img
            className='absolute top-[-2px] right-[-23px] w-[174px] h-[123px]'
            src={noteImg}
            alt=''
          />
        </div>
        <div className='flex flex-col items-start gap-2'>
          <div className='flex flex-col'>
            <LogEntryDatePicker event={event} formId={formId} />
          </div>
          <div className='flex flex-col'>
            <p className='text-sm font-semibold'>Author</p>
            <Tooltip label={authorEmail as string} hasArrow>
              <p className='text-sm'>{author}</p>
            </Tooltip>
          </div>

          <div className='flex flex-col w-full'>
            <p className='text-sm font-semibold'>Entry</p>

            {!isAuthor && (
              <HtmlContentRenderer
                className='text-sm'
                htmlContent={`${event.content}`}
              />
            )}
            {isAuthor && (
              <PreviewEditor
                formId={formId}
                initialContent={`${event.content}`}
                tags={data?.tags}
                onClose={closeModal}
              />
            )}
          </div>

          <PreviewTags
            isAuthor={isAuthor}
            tags={event?.tags}
            tagOptions={data?.tags}
            id={event.id}
          />

          {event?.externalLinks?.[0]?.externalUrl && (
            <LogEntryExternalLink externalLink={event?.externalLinks?.[0]} />
          )}
        </div>
      </div>
    </>
  );
};

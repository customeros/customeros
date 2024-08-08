import { useMemo, useCallback } from 'react';

import noteIcon from '@assets/images/event-ill-log-stub.png';

import { cn } from '@ui/utils/cn';
import { Phone } from '@ui/media/icons/Phone';
import { User, Contact } from '@graphql/types';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Calendar } from '@ui/media/icons/Calendar';
import { Card, CardContent } from '@ui/presentation/Card/Card';
import { MessageTextSquare01 } from '@ui/media/icons/MessageTextSquare01';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { HtmlContentRenderer } from '@ui/presentation/HtmlContentRenderer/HtmlContentRenderer';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface LogEntryStubProps {
  data: LogEntryWithAliases;
}

function getAuthor(user: User | Contact) {
  if (!user) return 'Unknown';

  if (user.name) {
    return user.name;
  }

  if (user.firstName || user.lastName) {
    return `${user.firstName ?? ''} ${user.lastName ?? ''}`.trim();
  }

  return 'Unknown';
}

export const LogEntryStub = ({ data }: LogEntryStubProps) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const isTemporary = !data?.updatedAt;

  const fullName = getAuthor(data?.logEntryCreatedBy);
  const getLogEntryIcon = useCallback((type: string | null) => {
    switch (type) {
      case 'email':
        return <Mail01 className='text-gray-500 size-3' />;
      case 'meeting':
        return <Calendar className='text-gray-500 size-3' />;
      case 'voicemail':
      case 'call':
        return <Phone className='text-gray-500 size-3' />;
      case 'text-message':
        return <MessageTextSquare01 className='text-gray-500 size-3' />;

      default:
        return null;
    }
  }, []);

  const getInlineTags = useCallback(() => {
    if (data.tags?.[0]?.name) {
      return data.tags?.[0]?.name;
    }
    const parser = new DOMParser();
    const doc = parser.parseFromString(`<p>${data?.content}</p>`, 'text/html');
    const element = doc.querySelector('.customeros-tag');

    // Return the inner HTML of the found element
    return element?.innerHTML || null;
  }, [data.tags, data.content]);

  const logEntryIcon = useMemo(() => {
    const firstTag = getInlineTags();
    const icon = getLogEntryIcon(firstTag);

    if (!icon) return null;

    return (
      <div className='flex mr-[10px] relative bg-white border border-gray-200 rounded-md p-2 right-[-12px] top-[4px]'>
        {icon}
      </div>
    );
  }, [getInlineTags]);

  return (
    <Card
      onClick={() => !isTemporary && openModal(data.id)}
      className={cn(
        isTemporary
          ? 'opacity-50 cursor-progress'
          : 'opacity-100 cursor-pointer',
        'hover:shadow-md max-w-[549px] flex flex-col bg-white ml-6 shadow-xs border border-gray-200 rounded-lg transition-all duration-200 ease-in-out',
      )}
    >
      <CardContent
        data-test='timeline-log-entry'
        className='px-3 py-2 flex-1 flex'
      >
        <div className='flex w-full justify-between relative h-fit'>
          <div className='w-[460px] line-clamp-4 text-sm text-gray-700 h-fit'>
            <span>{fullName}</span>
            <span className='text-gray-500 mx-1'>wrote</span>
            <HtmlContentRenderer
              showAsInlineText
              data-test='timeline-log-entry-text'
              htmlContent={`${data?.content?.substring(0, 500) || ''}`}
              className='relative pointer-events-none text-sm z-10 *:line-clamp-4'
            />
          </div>

          <div className='h-[86px]'>
            <div className='absolute top-[-2px] right-[-12px]'>
              <img alt='' height={94} width={124} src={noteIcon} />
            </div>
            {logEntryIcon && logEntryIcon}
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

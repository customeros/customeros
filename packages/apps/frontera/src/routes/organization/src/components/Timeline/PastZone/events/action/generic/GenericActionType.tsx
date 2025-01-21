import { useMemo } from 'react';

import { Action } from '@graphql/types';
import { Globe06 } from '@ui/media/icons/Globe06';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface GenericActionTypeProps {
  data: Action;
}

export const GenericActionType = ({ data }: GenericActionTypeProps) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();

  // todo, remove when content comes as valid markdown
  const formattedContent = useMemo(() => {
    return data?.content?.split('\n')?.map((line, index) => {
      if (line.trimStart().startsWith('- ')) {
        // Handle list items
        const text = line.replace('-', '').trim();

        return (
          <div key={index} className='list-item list-disc ml-4'>
            {text.startsWith('/') ? <span>{text}</span> : text}
          </div>
        );
      }

      return <div key={index}>{line}</div>;
    });
  }, [data?.content]);

  if (!data?.content) return null;

  return (
    <div>
      <div
        onClick={() => openModal(data.id)}
        className='flex cursor-pointer min-h-[40px]'
      >
        <Globe06 className='text-gray-500 mt-0.5' />
        <p className=' max-w-[500px] ml-2 text-sm text-gray-700 whitespace-pre-line'>
          {formattedContent}
        </p>
      </div>
    </div>
  );
};

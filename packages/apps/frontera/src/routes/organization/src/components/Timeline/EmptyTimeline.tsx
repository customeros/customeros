import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { TimelineActions } from '@organization/components/Timeline/FutureZone/TimelineActions/TimelineActions';

import { FutureZone } from './FutureZone/FutureZone';
import EmptyTimelineIllustration from './assets/EmptyTimelineIllustration';

interface EmptyTimelineProps {
  invalidateQuery: () => void;
}

export const EmptyTimeline = observer(
  ({ invalidateQuery }: EmptyTimelineProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.value.get(id);

    return (
      <div className='flex flex-col h-[calc(100vh-5rem)] overflow-auto w-full'>
        <div className='flex flex-col items-center flex-1 max-h-[50%] bg-[url(/backgrounds/organization/dotted-bg-pattern.svg)] bg-no-repeat bg-contain bg-center'>
          <div className='flex flex-col items-center justify-center h-full max-w-[390px]'>
            <EmptyTimelineIllustration />
            <h1 className='text-gray-900 text-lg font-semibold mt-3 mb-2'>
              {organization?.value?.name || 'Unknown'} has no events yet
            </h1>
            <span className='text-gray-600 text-xs text-center'>
              This organization’s events will show up here once a data source
              has been linked
            </span>
          </div>
        </div>
        <div className='flex bg-[#F9F9FB] flex-col flex-1'>
          <div>
            <TimelineActions invalidateQuery={invalidateQuery} />
          </div>
          <div className='flex flex-1 h-full bg-[#F9F9FB]'>
            <FutureZone />
          </div>
        </div>
      </div>
    );
  },
);

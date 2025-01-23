import { Link, useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Radar } from '@ui/media/icons/Radar';
// import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
// import { DotsVertical } from '@ui/media/icons/DotsVertical';

export const Header = observer(() => {
  const { id } = useParams<{ id: string }>();
  const store = useStore();

  if (!id) {
    throw new Error('No id provided');
  }

  const agent = store.agents.getById(id);

  return (
    <div className='w-full border-b border-b-gray-200 px-3 py-[11px] flex justify-between items-center'>
      <div className='flex items-center gap-1'>
        <Link
          to='/agents'
          className='text-md font-medium text-grayModern-500 hover:text-grayModern-700 transition-colors'
        >
          Agents
        </Link>
        <ChevronRight className='w-4 h-4 text-gray-500' />
        <div className='flex items-center gap-1'>
          <div className='flex items-center rounded-sm p-0.5 bg-grayModern-100'>
            <Radar className='w-4 h-4 text-grayModern-500' />
          </div>
          <div className='flex items-center gap-1'>
            <p className='text-md font-medium'>{agent?.value.name}</p>
            {/* <IconButton
              size='xxs'
              variant='ghost'
              icon={<DotsVertical />}
              aria-label='agent options'
            /> */}
          </div>
        </div>
      </div>
    </div>
  );
});

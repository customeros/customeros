import { useNavigate } from 'react-router-dom';

import { Radar } from '@ui/media/icons/Radar';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface AgentCardProps {
  id: string;
  name: string;
  icon?: string;
  hasError?: boolean;
  status: 'ON' | 'OFF';
}

export const AgentCard = ({ id, name, status }: AgentCardProps) => {
  const navigate = useNavigate();
  const tagColor = status === 'ON' ? 'success' : 'error';

  return (
    <div
      onClick={() => navigate(`/agents/${id}`)}
      className='p-3 rounded-lg flex items-center gap-2 border border-grayModern-200 mb-4 cursor-pointer hover:bg-white hover:border-grayModern-500 hover:shadow-lg transition-colors'
    >
      <div className='p-2 flex items-center'>
        <Radar className='text-grayModern-500' />
      </div>

      <div className='flex justify-between w-full'>
        <p className='line-clamp-1 text-sm font-medium'>{name}</p>
        <Tag colorScheme={tagColor}>
          <TagLabel>{status}</TagLabel>
        </Tag>
      </div>
    </div>
  );
};

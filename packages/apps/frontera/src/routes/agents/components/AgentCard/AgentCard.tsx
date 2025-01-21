import { Radar } from '@ui/media/icons/Radar';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface AgentCardProps {
  name: string;
  icon?: string;
  hasError?: boolean;
  status: 'ON' | 'OFF';
}

export const AgentCard = ({ name, status }: AgentCardProps) => {
  const tagColor = status === 'ON' ? 'success' : 'error';

  return (
    <div className='p-3 rounded-lg flex items-center gap-2 border border-grayModern-200 mb-4 cursor-pointer hover:bg-white transition-colors'>
      <div className='p-2 flex items-center'>
        <Radar className='text-grayModern-500' />
      </div>

      <div className='flex justify-between w-full'>
        <p className='line-clamp-1 font-medium'>{name}</p>
        <Tag colorScheme={tagColor}>
          <TagLabel>{status}</TagLabel>
        </Tag>
      </div>
    </div>
  );
};

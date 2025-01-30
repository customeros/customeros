import { useNavigate } from 'react-router-dom';

import { cn } from '@ui/utils/cn';
import { Icon, IconName } from '@ui/media/Icon';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface AgentCardProps {
  id: string;
  name: string;
  icon?: string;
  color: string;
  hasError?: boolean;
  status: 'ON' | 'OFF';
}

export const AgentCard = ({
  id,
  name,
  icon,
  color = 'grayModern',
  status,
}: AgentCardProps) => {
  const navigate = useNavigate();
  const tagColor = status === 'ON' ? 'success' : 'error';

  const [border, bg, iconColor] = [
    `border-${color}-200 hover:border-${color}-500`,
    `bg-${color}-50`,
    `text-${color}-500`,
  ];

  return (
    <div
      onClick={() => navigate(`/agents/${id}`)}
      className={cn(
        'p-3 rounded-lg flex items-center gap-2 border mb-4 cursor-pointer hover:bg-white hover:shadow-lg transition-colors',
        border,
      )}
    >
      <div className={cn('p-2 flex items-center rounded-lg', bg)}>
        <Icon
          className={cn('size-6', iconColor)}
          name={(icon as IconName) ?? 'radar'}
        />
      </div>

      <div className='flex justify-between w-full'>
        <p className='line-clamp-1 text-sm font-medium'>{name}</p>
        <Tag colorScheme={tagColor}>
          <TagLabel className='flex items-center gap-1'>
            <Icon
              stroke='none'
              name={status === 'OFF' ? 'dot-single' : 'dot-live-success'}
            />
            {status}
          </TagLabel>
        </Tag>
      </div>
    </div>
  );
};

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
  status,
  hasError,
  color = 'grayModern',
}: AgentCardProps) => {
  const navigate = useNavigate();
  const tagColor = hasError ? 'error' : status === 'ON' ? 'success' : 'error';

  const [ring, bg, iconColor] = colorMap[color];

  return (
    <div
      onClick={() => navigate(`/agents/${id}`)}
      className={cn(
        'p-3 rounded-lg flex items-center gap-2 ring-0 border mb-4 cursor-pointer hover:bg-white hover:shadow-lg transition-all ring-grayModern-200 hover:ring-1',
        ring,
      )}
    >
      <div
        className={cn('p-2 flex items-center rounded-lg bg-grayModern-50', bg)}
      >
        <Icon
          name={(icon as IconName) ?? 'radar'}
          className={cn('size-6 text-grayModern-500', iconColor)}
        />
      </div>

      <div className='flex justify-between w-full'>
        <p className='line-clamp-1 text-sm font-medium'>{name}</p>
        <Tag colorScheme={tagColor}>
          <TagLabel className='flex items-center gap-1'>
            <Icon
              stroke='none'
              name={
                hasError
                  ? 'dot-single'
                  : status === 'OFF'
                  ? 'dot-single'
                  : 'dot-live-success'
              }
            />
            {status}
          </TagLabel>
        </Tag>
      </div>
    </div>
  );
};

const colorMap: Record<string, [ring: string, bg: string, iconColor: string]> =
  {
    grayModern: [
      'hover:ring-grayModern-400',
      'hover:bg-grayModern-50',
      'hover:text-grayModern-500',
    ],
    error: [
      'hover:ring-error-400',
      'hover:bg-error-50',
      'hover:text-error-500',
    ],
    warning: [
      'hover:ring-warning-400',
      'hover:bg-warning-50',
      'hover:text-warning-500',
    ],
    success: [
      'hover:ring-success-400',
      'hover:bg-success-50',
      'hover:text-success-500',
    ],
    grayWarm: [
      'hover:ring-grayWarm-400',
      'hover:bg-grayWarm-50',
      'hover:text-grayWarm-500',
    ],
    moss: ['hover:ring-moss-400', 'hover:bg-moss-50', 'hover:text-moss-500'],
    blueLight: [
      'hover:ring-blueLight-400',
      'hover:bg-blueLight-50',
      'hover:text-blueLight-500',
    ],
    indigo: [
      'hover:ring-indigo-400',
      'hover:bg-indigo-50',
      'hover:text-indigo-500',
    ],
    violet: [
      'hover:ring-violet-400',
      'hover:bg-violet-50',
      'hover:text-violet-500',
    ],
    pink: ['hover:ring-pink-400', 'hover:bg-pink-50', 'hover:text-pink-500'],
  };

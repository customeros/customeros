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
  defaultName: string;
  status: 'ON' | 'OFF';
}

export const AgentCard = ({
  id,
  name,
  icon,
  status,
  hasError,
  defaultName,
  color = 'grayModern',
}: AgentCardProps) => {
  const navigate = useNavigate();
  const tagColor = hasError ? 'error' : status === 'ON' ? 'success' : 'error';

  const [ring, bg, iconColor] = colorMap[color];

  return (
    <div
      onClick={() => navigate(`/agents/${id}`)}
      className={cn(
        'p-3 min-w-[372px] flex-1 rounded-lg flex items-center gap-2 ring-0 border cursor-pointer group hover:bg-white hover:shadow-lg transition-all ring-grayModern-200 hover:ring-1',
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
        <div className='flex flex-col justify-center'>
          {name !== defaultName && (
            <p className='line-clamp-1 text-xs text-grayModern-500'>{name}</p>
          )}
          <p className='line-clamp-2 text-sm font-medium'>{defaultName}</p>
        </div>
        <div className='flex items-center'>
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
    </div>
  );
};

const colorMap: Record<string, [ring: string, bg: string, iconColor: string]> =
  {
    grayModern: [
      'group-hover:ring-grayModern-400',
      'group-hover:bg-grayModern-50',
      'group-hover:text-grayModern-500',
    ],
    error: [
      'group-hover:ring-error-400',
      'group-hover:bg-error-50',
      'group-hover:text-error-500',
    ],
    warning: [
      'group-hover:ring-warning-400',
      'group-hover:bg-warning-50',
      'group-hover:text-warning-500',
    ],
    success: [
      'group-hover:ring-success-400',
      'group-hover:bg-success-50',
      'group-hover:text-success-500',
    ],
    grayWarm: [
      'group-hover:ring-grayWarm-400',
      'group-hover:bg-grayWarm-50',
      'group-hover:text-grayWarm-500',
    ],
    moss: [
      'group-hover:ring-moss-400',
      'group-hover:bg-moss-50',
      'group-hover:text-moss-500',
    ],
    blueLight: [
      'group-hover:ring-blueLight-400',
      'group-hover:bg-blueLight-50',
      'group-hover:text-blueLight-500',
    ],
    indigo: [
      'group-hover:ring-indigo-400',
      'group-hover:bg-indigo-50',
      'group-hover:text-indigo-500',
    ],
    violet: [
      'group-hover:ring-violet-400',
      'group-hover:bg-violet-50',
      'group-hover:text-violet-500',
    ],
    pink: [
      'group-hover:ring-pink-400',
      'group-hover:bg-pink-50',
      'group-hover:text-pink-500',
    ],
  };

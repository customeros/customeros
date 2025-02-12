import { useNavigate } from 'react-router-dom';

import { cn } from '@ui/utils/cn';
import { Icon, IconName } from '@ui/media/Icon';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface AgentCardProps {
  id: string;
  name: string;
  icon?: string;
  hasError?: boolean;
  defaultName: string;
  status: 'ON' | 'OFF';
  colorMap: [ring: string, bg: string, iconColor: string];
}

export const AgentCard = ({
  id,
  name,
  icon,
  status,
  hasError,
  defaultName,
  colorMap,
}: AgentCardProps) => {
  const navigate = useNavigate();
  const tagColor =
    hasError && status === 'ON'
      ? 'warning'
      : status === 'OFF'
      ? 'grayModern'
      : 'success';

  const [ring, bg, iconColor] = colorMap;

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
          {name.toLowerCase() !== defaultName.toLowerCase() && (
            <p className='line-clamp-1 text-xs text-grayModern-500'>
              {defaultName}
            </p>
          )}
          <p className='line-clamp-2 text-sm font-medium'>
            {name || defaultName}
          </p>
        </div>
        <div className='flex items-center'>
          <Tag colorScheme={tagColor}>
            <TagLabel className='flex items-center gap-1'>
              <Icon
                stroke={hasError && status === 'ON' ? 'currentColor' : 'none'}
                className={cn('size-3', {
                  'text-warning-500': hasError && status === 'ON',
                })}
                name={
                  hasError && status === 'ON'
                    ? 'alert-triangle'
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

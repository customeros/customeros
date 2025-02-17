import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { UpdateAgentIconUsecase } from '@domain/usecases/agents/update-agent-icon.usecase.ts';
import { UpdateAgentColorUsecase } from '@domain/usecases/agents/update-agent-color.usecase.ts';

import { cn } from '@ui/utils/cn.ts';
import { Input } from '@ui/form/Input';
import { Icon, IconName } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';

export const IconAndColorPicker = observer(() => {
  const { id } = useParams();
  const iconUsecase = useMemo(() => new UpdateAgentIconUsecase(id), [id]);
  const colorUsecase = useMemo(() => new UpdateAgentColorUsecase(id), [id]);

  return (
    <div className='flex flex-col gap-2 w-64'>
      <div className={'w-full flex justify-evenly gap-2  bg-grayModern-700'}>
        {Object.entries(colorMap).map(([color, bg]) => (
          <Button
            key={color}
            variant={'ghost'}
            onClick={() => colorUsecase.execute(color)}
            className='p-0 hover:bg-transparent hover:focus:bg-transparent cursor-pointer '
          >
            <div
              className={cn(`size-4 rounded-full ring-0 `, bg, {
                [`ring-1 ring-offset-2 ring-offset-grayModern-700`]:
                  colorUsecase.agent?.value.color === color,
              })}
            />
          </Button>
        ))}
      </div>

      <div className='relative border-t border-grayModern-600 pt-2 mt-[6px]'>
        <Input
          type='text'
          size={'sm'}
          variant={'unstyled'}
          placeholder='Search'
          value={iconUsecase.searchQuery}
          className={'text-grayModern-25'}
          onChange={(e) => iconUsecase.setSearchQuery(e.target.value)}
        />
      </div>

      <div className='grid grid-cols-10 gap-2 overflow-y-auto p-1'>
        {iconUsecase.iconOptions.map((iconName: IconName) => (
          <IconButton
            size={'xxs'}
            key={iconName}
            variant={'ghost'}
            colorScheme={'gray'}
            aria-label={iconName}
            icon={
              <Icon name={iconName} className={'text-inherit size-[14px]'} />
            }
            onClick={(e) => {
              e.stopPropagation();
              iconUsecase.execute(iconName);
            }}
            className={cn(
              `text-grayModern-50 hover:text-grayModern-50 hover:bg-grayModern-600 focus:text-grayModern-50 focus:bg-grayModern-600 transition-colors `,
              {
                '!bg-grayModern-600 !text-grayModern-25 hover:bg-grayModern-600 !hover:text-grayModern-50':
                  iconUsecase.agent?.value.icon === iconName,
              },
            )}
          />
        ))}
      </div>
    </div>
  );
});

const colorMap = {
  grayModern: 'bg-grayModern-400 ring-grayModern-400',
  error: 'bg-error-400 ring-error-400',
  warning: 'bg-warning-400 ring-warning-400',
  success: 'bg-success-400 ring-success-400',
  grayWarm: 'bg-grayWarm-400 ring-grayWarm-400',
  moss: 'bg-moss-400 ring-moss-400',
  blueLight: 'bg-blueLight-400 ring-blueLight-400',
  indigo: 'bg-indigo-400 ring-indigo-400',
  violet: 'bg-violet-400 ring-violet-400',
  pink: 'bg-pink-400 ring-pink-400',
};

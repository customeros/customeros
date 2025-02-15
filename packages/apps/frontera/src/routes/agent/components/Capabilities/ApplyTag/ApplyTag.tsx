import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { AddTagToCompanyUsecase } from '@domain/usecases/agents/capabilities/add-tag-to-company.usecase';

import { cn } from '@ui/utils/cn.ts';
import { Icon } from '@ui/media/Icon';
import { Combobox } from '@ui/form/Combobox';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Tag, TagLabel, TagRightButton } from '@ui/presentation/Tag';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

export const ApplyTag = observer(() => {
  const store = useStore();
  const { id } = useParams<{ id: string }>();
  const usecase = useMemo(() => new AddTagToCompanyUsecase(id!), [id]);
  const tag = usecase.selectedTag
    ? store.tags.getById(usecase.selectedTag.value)
    : null;

  return (
    <div className='flex flex-col gap-4'>
      <p className='font-semibold text-sm'>Apply a tag to company</p>

      {usecase.listenerErrors && (
        <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] mb-4'>
          <Icon stroke='none' className='mr-2' name='dot-single' />
          <span className='text-sm'>{usecase.listenerErrors}</span>
        </div>
      )}

      <div>
        <p className='font-semibold text-sm'>Tag name</p>
        <p className='text-sm mb-2'>
          Choose or create one tag to identify companies that need support
        </p>

        <Popover>
          <Tooltip align='start' label='Company tags'>
            <PopoverTrigger className={cn('flex items-center w-full gap-2')}>
              <Icon name={'tag-01'} className='text-grayModern-500' />
              <div
                data-test={'apply-tag'}
                className='flex flex-wrap gap-1 w-fit items-center'
              >
                {tag ? (
                  <Tag
                    size={'md'}
                    variant='subtle'
                    className={'gap-1'}
                    colorScheme={
                      // eslint-disable-next-line @typescript-eslint/no-explicit-any
                      (tag?.colorScheme as unknown as any) ?? 'grayModern'
                    }
                  >
                    <TagLabel>{tag.tagName}</TagLabel>
                    <TagRightButton
                      onClick={(e) => {
                        e.stopPropagation();
                        usecase.execute();
                      }}
                    >
                      <Icon name={'x-close'} />
                    </TagRightButton>
                  </Tag>
                ) : (
                  <span className='text-gray-400 text-sm'>Company tag</span>
                )}
              </div>
            </PopoverTrigger>
          </Tooltip>
          <PopoverContent align='start' className='min-w-[264px] w-[375px]'>
            <Combobox
              isMulti={false}
              options={usecase.tagList}
              placeholder='Company tag'
              value={usecase.selectedTag}
              inputValue={usecase.searchTerm}
              onInputChange={usecase.setSearchTerm}
              onChange={(newValue) => {
                usecase.execute(newValue);
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !usecase.tagList.length) {
                  usecase.create();
                }
              }}
              noOptionsMessage={({ inputValue }) => (
                <div
                  onClick={() => {
                    usecase.create();
                  }}
                  className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'
                >
                  <Plus />
                  <span>{`Create "${inputValue}"`}</span>
                </div>
              )}
            />
          </PopoverContent>
        </Popover>
      </div>
    </div>
  );
});

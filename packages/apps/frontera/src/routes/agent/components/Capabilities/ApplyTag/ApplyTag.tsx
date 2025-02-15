import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { AddTagToCompanyUsecase } from '@domain/usecases/agents/capabilities/add-tag-to-company.usecase';

import { Icon } from '@ui/media/Icon';
import { Tags } from '@shared/components/Tags';

export const ApplyTag = observer(() => {
  const { id } = useParams<{ id: string }>();
  const usecase = useMemo(() => new AddTagToCompanyUsecase(id!), [id]);

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
        <p className='text-sm'>
          Choose or create one tag to identify companies that need support
        </p>
        <Tags
          isMulti={false}
          className='mt-1'
          options={usecase.tagList}
          placeholder='Company tag'
          onCreate={usecase.create}
          value={usecase.selectedTags}
          inputValue={usecase.searchTerm}
          setInputValue={usecase.setSearchTerm}
          leftAccessory={<Icon name='tag-01' className='mr-3 text-gray-500' />}
          onChange={(selection) => {
            usecase.execute(selection?.value);
          }}
        />
      </div>
    </div>
  );
});

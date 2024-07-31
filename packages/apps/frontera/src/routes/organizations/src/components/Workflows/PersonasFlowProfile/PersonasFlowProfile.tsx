import { Tag, TagLabel } from '@ui/presentation/Tag';
import { getContainerClassNames } from '@ui/form/Select';
import { Certificate02 } from '@ui/media/icons/Certificate02';

import { MultiSelectFilter } from '../components';

export const PersonasFlowProfile = () => {
  return (
    <>
      <div className='flex'>
        <p className='font-semibold'> Tag contact with </p>
        <Tag variant='subtle'>
          <TagLabel>Solo RevOps</TagLabel>
        </Tag>
      </div>
      <p className='font-medium leading-5 text-gray-500 mt-4 mb-2'>WHEN</p>

      <div>
        <MultiSelectFilter
          label='Job Title'
          description='is any of'
          placeholder='Job titles'
          icon={<Certificate02 className='mr-2 text-gray-500' />}
          classNames={{
            container: () => getContainerClassNames(undefined, 'unstyled', {}),
          }}
          options={[
            { label: 'CEO', value: 'CEO' },
            { label: 'CFO', value: 'CFO' },
            { label: 'COO', value: 'COO' },
            { label: 'CMO', value: 'CMO' },
            { label: 'CTO', value: 'CTO' },
            { label: 'VP', value: 'VP' },
            { label: 'Director', value: 'Director' },
            { label: 'Manager', value: 'Manager' },
          ]}
        />
      </div>
    </>
  );
};

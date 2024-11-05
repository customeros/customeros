import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { RadioButton } from '@ui/media/icons/RadioButton';
import { ColumnViewType } from '@shared/types/__generated__/graphql.types';

import { Layout } from '../components';
import { getFieldTypes } from './filedTypes';
import { Header } from '../components/Header';
import { CustomFieldItem } from '../components/CustomFieldItem';

export const OrganizationFields = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();

  const orgPreset = store.tableViewDefs.organizationsPreset;
  const coreFields =
    store.tableViewDefs.getById(orgPreset || '')?.value.columns || [];
  const search = searchParams?.get('search') || '';

  const activeTab = (tab: string) => searchParams?.get('view') === tab;
  const customFieldTypes = getFieldTypes(store);

  const filteredFields = coreFields.filter((field) => {
    const fieldName = customFieldTypes[field.columnType]?.fieldName || '';

    return (
      field.columnType !== ColumnViewType.OrganizationsAvatar &&
      (!search || fieldName.toLowerCase().includes(search.toLowerCase()))
    );
  });

  return (
    <Layout>
      <Header
        numberOfCustomFields={0}
        title='Organization Fields'
        numberOfCoreFields={filteredFields.length}
        subTitle='Create and manage custom fields for organizations'
      />

      <div className='flex items-center justify-between px-2 mt-4'>
        <p className='flex flex-2 font-medium'>Field name</p>
        <p className='flex flex-1 font-medium'>Type</p>
      </div>
      {activeTab('core') ? (
        filteredFields.map((field) => (
          <div
            key={field.columnId}
            className='flex justify-between items-center'
          >
            <CustomFieldItem field={field} store={store} />
          </div>
        ))
      ) : (
        <div className='flex items-center'>
          mariana
          <RadioButton />
        </div>
      )}
    </Layout>
  );
});

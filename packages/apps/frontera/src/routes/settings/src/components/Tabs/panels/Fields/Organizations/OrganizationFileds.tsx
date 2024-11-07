import { useSearchParams } from 'react-router-dom';

import { toJS } from 'mobx';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ColumnViewType } from '@shared/types/__generated__/graphql.types';

import { Layout } from '../components';
import { Header } from '../components/Header';
import { getDefaultFieldTypes } from './filedTypes';
import { CustomFieldItem } from '../components/CustomFieldItem';

export const OrganizationFields = observer(() => {
  const store = useStore();
  const [searchParams] = useSearchParams();

  const customFieldStore = store.customFields.toArray();

  const orgPreset = store.tableViewDefs.organizationsPreset;

  const coreFields =
    store.tableViewDefs.getById(orgPreset || '')?.value.columns || [];
  const search = searchParams?.get('search') || '';

  const activeTab = (tab: string) => searchParams?.get('view') === tab;
  const customFieldTypes = getDefaultFieldTypes(store);

  const filteredFields = coreFields.filter((field) => {
    const fieldName = customFieldTypes[field.columnType]?.fieldName || '';

    return (
      field.columnType !== ColumnViewType.OrganizationsAvatar &&
      (!search || fieldName.toLowerCase().includes(search.toLowerCase()))
    );
  });

  const customFields = customFieldStore
    .filter((f) => f.value.name.toLowerCase().includes(search.toLowerCase()))
    .map((f) => toJS(f.value));

  const filteredCustomFields = activeTab('custom') ? customFields : [];

  return (
    <Layout>
      <Header
        title='Organization fields'
        numberOfCoreFields={coreFields.length}
        numberOfCustomFields={customFieldStore.length}
        subTitle='Create and manage custom fields for organizations'
      />

      <div className='flex items-center justify-between mt-4'>
        {activeTab('core') && coreFields.length > 0 && (
          <>
            <p className='flex flex-2 font-medium text-sm'>Field name</p>
            <p className='flex flex-1 font-medium text-sm'>Type</p>
          </>
        )}
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
        <>
          {filteredCustomFields.length === 0 ? (
            search ? (
              <p>
                Nothing to search for yet. Go ahead, add you first custom field
              </p>
            ) : (
              <p className='text-center text-sm text-gray-500'>
                Nothing to search for yet. Go ahead, add you first custom
                field...
              </p>
            )
          ) : (
            filteredCustomFields.map((field) => (
              <div key={field.id} className='flex justify-between items-center'>
                <CustomFieldItem
                  field={field}
                  store={store}
                  isEditable={true}
                />
              </div>
            ))
          )}
        </>
      )}
    </Layout>
  );
});

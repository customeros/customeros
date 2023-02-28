import React, { useState } from 'react';
import { ChevronDown, User } from '../../ui-kit/atoms';
import { Controller } from 'react-hook-form';
import { SearchWithOverlay } from '../../ui-kit/atoms/search';
import { useUsers } from '../../../hooks/useUser';
import styles from './contact-details.module.scss';
import { Filter } from '../../../graphQL/__generated__/generated';

interface ContactOwnerProps {
  control: any;
  setValue: any;
}

export const ContactOwner: React.FC<ContactOwnerProps> = ({
  control,
  setValue,
}) => {
  const { data, loading, error, onLoadUsers } = useUsers();

  return (
    <div className={styles.searchWrapper}>
      <Controller
        name='ownerFullName'
        control={control}
        rules={{ required: true }}
        render={({ field }) => (
          <SearchWithOverlay
            options={data?.content}
            loadingOptions={loading}
            resourceLabel='users'
            value={field.value}
            searchBy={[
              { label: 'First name', field: 'FIRST_NAME' },
              { label: 'Last name', field: 'LAST_NAME' },
            ]}
            searchData={(where: Filter) => {
              return onLoadUsers({
                variables: {
                  pagination: { limit: 999, page: 0 },
                  where: where,
                },
              });
            }}
            itemTemplate={(e: any) => {
              return (
                <>
                  <span className='mr-3'>
                    <User style={{ transform: 'scale(0.8)' }} />
                  </span>
                  {e.firstName} {e.lastName}
                </>
              );
            }}
            onItemSelected={(e: any) => {
              setValue('ownerId', !e ? undefined : e.id);
              setValue(
                'ownerFullName',
                !e ? '-' : e.firstName + ' ' + e.lastName,
              );
            }}
            maxResults={5}
          />
        )}
      />
    </div>
  );
};

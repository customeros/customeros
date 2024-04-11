'use client';
import React, { FC } from 'react';

import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';

import { FormSelect } from '@ui/form/SyncSelect';
import { FormInput } from '@ui/form/Input/FormInput2';
import { SelectOption } from '@shared/types/SelectOptions';
import { countryOptions } from '@shared/util/countryOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { EmailsInputGroup } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/EmailsInputGroup/EmailsInputGroup';
import { BillingAddressDetailsFormDto } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm.dto';

interface BillingAddressDetailsForm {
  formId: string;
  values: BillingAddressDetailsFormDto;
  country?: SelectOption<string> | null;
}

export const BillingDetailsForm: FC<BillingAddressDetailsForm> = ({
  country,
  formId,
  values,
}) => {
  const client = getGraphQLClient();
  const { data: tenantSettingsData } = useTenantSettingsQuery(client);

  return (
    <div className='flex flex-col mt-2'>
      <FormInput
        label='Organization legal name'
        labelProps={{
          className: 'text-sm mb-0 font-semibold',
        }}
        formId={formId}
        name='organizationLegalName'
        placeholder='Organization legal name'
        autoComplete='off'
        className='overflow-hidden overflow-ellipsis mb-1'
      />

      <div className='flex flex-col'>
        <p className='text-sm font-semibold'>Billing address</p>
        <FormSelect
          label='Country'
          placeholder='Country'
          name='country'
          formId={formId}
          options={countryOptions}
        />
        <FormInput
          label='Address line 1'
          formId={formId}
          name='addressLine1'
          placeholder='Address line 1'
          isLabelHidden
          autoComplete='off'
          className='overflow-hidden overflow-ellipsis'
        />
        <FormInput
          label='Address line 2'
          formId={formId}
          name='addressLine2'
          isLabelHidden
          placeholder='Address line 2'
          autoComplete='off'
          className='overflow-hidden overflow-ellipsis'
        />
        {country?.value === 'US' && (
          <FormInput
            label='City'
            formId={formId}
            name='locality'
            placeholder='City'
            isLabelHidden
            autoComplete='off'
            className='overflow-hidden overflow-ellipsis'
          />
        )}
        <div className='flex'>
          {country?.value === 'US' ? (
            <FormInput
              label='State'
              name='region'
              placeholder='State'
              formId={formId}
              isLabelHidden
            />
          ) : (
            <FormInput
              label='City'
              formId={formId}
              name='locality'
              placeholder='City'
              autoComplete='off'
              className='overflow-hidden overflow-ellipsis'
              isLabelHidden
            />
          )}
          <FormInput
            label='ZIP/Postal code'
            formId={formId}
            name='postalCode'
            isLabelHidden
            placeholder='ZIP/Postal code'
            autoComplete='off'
            className='overflow-hidden overflow-ellipsis'
          />
        </div>
      </div>

      {tenantSettingsData?.tenantSettings?.billingEnabled && (
        <EmailsInputGroup
          formId={formId}
          to={values?.billingEmail}
          cc={values?.billingEmailCC ?? []}
          bcc={values?.billingEmailBCC ?? []}
        />
      )}
    </div>
  );
};

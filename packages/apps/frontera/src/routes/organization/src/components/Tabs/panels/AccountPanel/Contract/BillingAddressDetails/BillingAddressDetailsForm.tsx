import { FC } from 'react';

import { Input } from '@ui/form/Input';
import { Select } from '@ui/form/Select';
import { useStore } from '@shared/hooks/useStore';
import { countryOptions } from '@shared/util/countryOptions';
import { EmailsInputGroup } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/EmailsInputGroup/EmailsInputGroup';

interface BillingAddressDetailsForm {
  contractId: string;
}

export const BillingAddressDetailsForm: FC<BillingAddressDetailsForm> = ({
  contractId,
}) => {
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId);

  const tenantSettings = store.settings.tenant.value;

  const handleUpdateBillingDetails = (key: string, value: string) => {
    contractStore?.update((contract) => ({
      ...contract,
      billingDetails: {
        ...contract.billingDetails,
        [key]: value,
      },
    }));
  };

  return (
    <div className='flex flex-col mt-2 min-w-[342px]'>
      <Input
        name='organizationLegalName'
        placeholder='Organization legal name'
        autoComplete='off'
        className='overflow-hidden overflow-ellipsis mb-4'
        variant='unstyled'
        value={
          contractStore?.value?.billingDetails?.organizationLegalName || ''
        }
        onChange={(e) => {
          handleUpdateBillingDetails('organizationLegalName', e.target.value);
        }}
        size='xs'
      />

      <div className='flex flex-col'>
        <p className='text-sm font-semibold'>Billing address</p>
        <Select
          placeholder='Country'
          name='country'
          options={countryOptions}
          value={contractStore?.value?.billingDetails?.country}
        />
        <Input
          name='addressLine1'
          placeholder='Address line 1'
          autoComplete='off'
          className='overflow-hidden overflow-ellipsis'
          value={contractStore?.value?.billingDetails?.addressLine1 ?? ''}
          onChange={(e) => {
            handleUpdateBillingDetails('addressLine1', e.target.value);
          }}
        />
        <Input
          name='addressLine2'
          placeholder='Address line 2'
          autoComplete='off'
          className='overflow-hidden overflow-ellipsis'
          value={contractStore?.value?.billingDetails?.addressLine2 ?? ''}
          onChange={(e) => {
            handleUpdateBillingDetails('addressLine2', e.target.value);
          }}
        />
        {contractStore?.value?.billingDetails?.country === 'US' && (
          <Input
            name='locality'
            placeholder='City'
            autoComplete='off'
            className='overflow-hidden overflow-ellipsis'
            value={contractStore?.value?.billingDetails?.locality ?? ''}
            onChange={(e) => {
              handleUpdateBillingDetails('locality', e.target.value);
            }}
          />
        )}
        <div className='flex'>
          {contractStore?.value?.billingDetails?.country === 'US' ? (
            <Input
              name='region'
              placeholder='State'
              value={contractStore?.value?.billingDetails?.region ?? ''}
              onChange={(e) => {
                handleUpdateBillingDetails('region', e.target.value);
              }}
            />
          ) : (
            <Input
              placeholder='City'
              autoComplete='off'
              className='overflow-hidden overflow-ellipsis'
              value={contractStore?.value?.billingDetails?.locality ?? ''}
              onChange={(e) => {
                handleUpdateBillingDetails('locality', e.target.value);
              }}
            />
          )}
          <Input
            name='postalCode'
            placeholder='ZIP/Postal code'
            autoComplete='off'
            className='overflow-hidden overflow-ellipsis'
            value={contractStore?.value?.billingDetails?.postalCode ?? ''}
            onChange={(e) => {
              handleUpdateBillingDetails('postalCode', e.target.value);
            }}
          />
        </div>
      </div>

      {tenantSettings?.billingEnabled && (
        <EmailsInputGroup contractId={contractId} />
      )}
    </div>
  );
};

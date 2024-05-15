import { LogoUploader } from '@settings/components/LogoUploadComponent/LogoUploader';
import { PaymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/PaymentMethods';

import { FormInput } from '@ui/form/Input/FormInput';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { FormSelect } from '@ui/form/Select/FormSelect';
import { FormSwitch } from '@ui/form/Switch/FormSwitch';
import { SelectOption } from '@shared/types/SelectOptions';
import { Divider } from '@ui/presentation/Divider/Divider';
import { countryOptions } from '@shared/util/countryOptions';
import { FormMaskInput } from '@ui/form/Input/FormMaskInput';
import { currencyOptions } from '@shared/util/currencyOptions';

const opts = {
  mask: 'AA 000 000 000',
  definitions: {
    A: /[A-Za-z]/,
    '0': /[0-9]/,
  },
  prepare: function (value: string, mask: { _value: string }) {
    if (mask._value.length < 2) {
      return value.toUpperCase();
    }

    return value;
  },
  format: function (value: string) {
    return value.toUpperCase(); // Ensure the country code is uppercase
  },
  parse: function (value: string) {
    return value.toUpperCase(); // Ensure the input is treated as uppercase
  },
};

export const TenantBillingPanelDetailsForm = ({
  setIsInvoiceProviderDetailsHovered,
  setIsInvoiceProviderFocused,
  formId,
  country,
  legalName,
}: {
  formId: string;
  legalName?: string | null;
  country?: SelectOption<string> | null;
  setIsInvoiceProviderFocused: (newState: boolean) => void;
  setIsInvoiceProviderDetailsHovered: (newState: boolean) => void;
}) => {
  return (
    <div className='flex flex-col px-6 py-5 w-full gap-4'>
      <LogoUploader />
      <FormInput
        autoComplete='off'
        label='Organization legal name'
        placeholder='Legal name'
        labelProps={{
          className: 'text-sm mb-0 font-semibold',
        }}
        name='legalName'
        formId={formId}
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
        onFocus={() => setIsInvoiceProviderFocused(true)}
        onBlur={() => setIsInvoiceProviderFocused(false)}
      />

      <FormSelect
        label='Base currency'
        placeholder='Invoice currency'
        isLabelVisible
        name='baseCurrency'
        formId={formId}
        options={currencyOptions}
      />

      <div
        className='flex flex-col'
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
      >
        <div className='text-sm font-semibold'>Billing address</div>
        <FormSelect
          name='country'
          placeholder='Country'
          label='Country'
          formId={formId}
          options={countryOptions}
        />
        <FormInput
          autoComplete='off'
          label='Address line 1'
          placeholder='Address line 1'
          labelProps={{
            className: 'hidden',
          }}
          name='addressLine1'
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
        <FormInput
          autoComplete='off'
          label='Billing address line 2'
          name='addressLine2'
          placeholder='Address line 2'
          labelProps={{
            className: 'hidden',
          }}
          formId={formId}
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />

        {country?.value === 'US' && (
          <FormInput
            className='overflow-ellipsis'
            label='City'
            formId={formId}
            name='locality'
            placeholder='City'
            labelProps={{
              className: 'hidden',
            }}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
            autoComplete='off'
          />
        )}

        <div className='flex gap-2'>
          {country?.value === 'US' ? (
            <FormInput
              label='State'
              name='region'
              placeholder='State'
              labelProps={{
                className: 'hidden',
              }}
              formId={formId}
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
            />
          ) : (
            <FormInput
              className='overflow-ellipsis'
              label='City'
              formId={formId}
              name='locality'
              placeholder='City'
              labelProps={{
                className: 'hidden',
              }}
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
              autoComplete='off'
            />
          )}
          <FormInput
            autoComplete='off'
            label='Billing address zip/Postal code'
            name='zip'
            placeholder='ZIP/Postal code'
            labelProps={{
              className: 'hidden',
            }}
            formId={formId}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
          />
        </div>

        <FormMaskInput
          formId={formId}
          name='vatNumber'
          autoComplete='off'
          label='VAT number'
          options={{ opts }}
          labelProps={{
            className: 'text-sm mb-0 mt-4 font-semibold inline-block ',
          }}
          placeholder='VAT number'
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
      </div>
      <div className='flex flex-col'>
        <div className='flex relative items-center'>
          <span className='text-sm whitespace-nowrap mr-2 text-gray-500'>
            Email invoice
          </span>
          <Divider />
        </div>

        <Tooltip label='This email is configured by CustomerOS'>
          <FormInput
            formId={formId}
            autoComplete='off'
            label='From'
            labelProps={{
              className: 'text-sm mb-0 font-semibold inline-block pt-4',
            }}
            name='sendInvoicesFrom'
            placeholder=''
            onFocus={() => setIsInvoiceProviderFocused(true)}
          />
        </Tooltip>

        <FormInput
          className='overflow-ellipsis'
          autoComplete='off'
          label='BCC'
          labelProps={{
            className: 'text-sm mb-0 font-semibold inline-block pt-4',
          }}
          formId={formId}
          name='sendInvoicesBcc'
          placeholder='BCC'
          type='email'
        />
      </div>
      <PaymentMethods formId={formId} legalName={legalName} />
      <FormSwitch
        size='sm'
        name='check'
        formId={formId}
        label='Checks'
        labelProps={{
          className: 'text-sm font-semibold m-0',
        }}
      />
    </div>
  );
};

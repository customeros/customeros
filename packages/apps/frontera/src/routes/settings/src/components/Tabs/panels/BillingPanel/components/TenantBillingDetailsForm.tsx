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
        formId={formId}
        name='legalName'
        autoComplete='off'
        placeholder='Legal name'
        label='Organization legal name'
        onFocus={() => setIsInvoiceProviderFocused(true)}
        onBlur={() => setIsInvoiceProviderFocused(false)}
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
        labelProps={{
          className: 'text-sm mb-0 font-semibold',
        }}
      />

      <FormSelect
        isLabelVisible
        formId={formId}
        name='baseCurrency'
        label='Base currency'
        options={currencyOptions}
        placeholder='Invoice currency'
      />

      <div
        className='flex flex-col'
        onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
        onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
      >
        <div className='text-sm font-semibold'>Billing address</div>
        <FormSelect
          name='country'
          label='Country'
          formId={formId}
          placeholder='Country'
          options={countryOptions}
        />
        <FormInput
          formId={formId}
          autoComplete='off'
          name='addressLine1'
          label='Address line 1'
          placeholder='Address line 1'
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
          labelProps={{
            className: 'hidden',
          }}
        />
        <FormInput
          formId={formId}
          autoComplete='off'
          name='addressLine2'
          placeholder='Address line 2'
          label='Billing address line 2'
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
          labelProps={{
            className: 'hidden',
          }}
        />

        {country?.value === 'US' && (
          <FormInput
            label='City'
            formId={formId}
            name='locality'
            placeholder='City'
            autoComplete='off'
            className='overflow-ellipsis'
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
            labelProps={{
              className: 'hidden',
            }}
          />
        )}

        <div className='flex gap-2'>
          {country?.value === 'US' ? (
            <FormInput
              label='State'
              name='region'
              formId={formId}
              placeholder='State'
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
              labelProps={{
                className: 'hidden',
              }}
            />
          ) : (
            <FormInput
              label='City'
              formId={formId}
              name='locality'
              placeholder='City'
              autoComplete='off'
              className='overflow-ellipsis'
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
              labelProps={{
                className: 'hidden',
              }}
            />
          )}
          <FormInput
            name='zip'
            formId={formId}
            autoComplete='off'
            placeholder='ZIP/Postal code'
            label='Billing address zip/Postal code'
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
            labelProps={{
              className: 'hidden',
            }}
          />
        </div>

        <FormMaskInput
          formId={formId}
          name='vatNumber'
          autoComplete='off'
          label='VAT number'
          options={{ opts }}
          placeholder='VAT number'
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
          labelProps={{
            className: 'text-sm mb-0 mt-4 font-semibold inline-block ',
          }}
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
            label='From'
            placeholder=''
            formId={formId}
            autoComplete='off'
            name='sendInvoicesFrom'
            onFocus={() => setIsInvoiceProviderFocused(true)}
            labelProps={{
              className: 'text-sm mb-0 font-semibold inline-block pt-4',
            }}
          />
        </Tooltip>

        <FormInput
          label='BCC'
          type='email'
          formId={formId}
          placeholder='BCC'
          autoComplete='off'
          name='sendInvoicesBcc'
          className='overflow-ellipsis'
          labelProps={{
            className: 'text-sm mb-0 font-semibold inline-block pt-4',
          }}
        />
      </div>
      <PaymentMethods formId={formId} legalName={legalName} />
      <FormSwitch
        size='sm'
        name='check'
        label='Checks'
        formId={formId}
        labelProps={{
          className: 'text-sm font-semibold m-0',
        }}
      />
    </div>
  );
};

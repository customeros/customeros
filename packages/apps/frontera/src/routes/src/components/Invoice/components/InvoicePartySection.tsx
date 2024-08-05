import { FC } from 'react';

import countries from '@assets/countries/countries.json';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
type InvoiceHeaderProps = {
  zip?: string;
  title: string;
  name?: string;
  email?: string;
  region?: string;
  country?: string;
  locality?: string;
  vatNumber?: string;
  isBlurred?: boolean;
  isFocused?: boolean;
  onClick?: () => void;
  addressLine1?: string;
  addressLine2?: string;
};

export const InvoicePartySection: FC<InvoiceHeaderProps> = ({
  isBlurred,
  isFocused,
  zip = '',
  name = '',
  email = '',
  country = '',
  locality = '',
  addressLine1 = '',
  addressLine2 = '',
  region = '',
  title,
  onClick,
}) => {
  const isUSA = country === 'United States of America';
  const borderRightPosition = title === 'From' ? 'border-r-0' : 'border-r';
  const filterDynamicClass = isBlurred ? 'blur-[2px]' : 'filter-none';
  const oppacity = isFocused
    ? 'data-[focus=true]:opacity-100'
    : 'data-[focus=true]:opacity-0';

  const showOnlyButton =
    !zip && !email && !locality && !addressLine1 && !addressLine2 && onClick;

  return (
    <Tooltip label={onClick ? 'Edit billing details' : ''}>
      <div
        onClick={onClick}
        data-focus={isFocused}
        tabIndex={onClick ? 0 : -1}
        role={onClick ? 'button' : 'none'}
        className={cn(
          'data-[focus=true]:transition-opacity data-[focus=true]:ring-2 data-[focus=true]:ring-gray-700 data-[focus=true]:delay-250 data-[focus=true]:ease-in-out data-[focus=true]:filter',
          'flex flex-col flex-1 w-[170px] py-2 px-3 border-t border-b border-gray-300 relative transition duration-250 ease-in-out filter',
          borderRightPosition,
          filterDynamicClass,
          oppacity,
          {
            'hover:ring-gray-700 hover:ring-2 hover:transition-opacity hover:delay-250 hover:ease-in-out hover:filter':
              onClick,
          },
        )}
      >
        <span className='font-semibold mb-1 text-sm'>{title}</span>
        {showOnlyButton && (
          <div>
            <Button
              size='xs'
              variant='link'
              onClick={onClick}
              colorScheme='primary'
              className='p-0 font-medium text-primary-600 shadow-none'
            >
              Add billing details
            </Button>
          </div>
        )}

        {!showOnlyButton && (
          <>
            <span className='text-sm leading-5 mb-1 font-medium'>{name}</span>
            {/* this is a left over from the original code */}
            {/*{vatNumber && (*/}
            {/*  <Text fontSize='xs' mb={1} lineHeight={1.2}>*/}
            {/*    VAT number: {vatNumber}*/}
            {/*  </Text>*/}
            {/*)}*/}

            <span className='text-sm text-gray-500 leading-5 break-words'>
              {addressLine1}
              <span className='block leading-4'>{addressLine2}</span>
            </span>

            {isUSA && (
              <span className='leading-4 text-gray-500 text-sm break-words'>
                {locality && `${locality}, `} {region} {zip}
              </span>
            )}
            {!isUSA && (
              <span className='text-sm leading-4 text-gray-500 break-words'>
                {locality}
                {locality && zip && ', '} {zip}
              </span>
            )}

            <span className='text-sm leading-4 text-gray-500'>
              {countries
                .find((c) => c.name === country)
                ?.alpha3.toLocaleUpperCase()}
            </span>
            {email && (
              <span className='text-sm leading-4 text-gray-500 break-words'>
                {email}
              </span>
            )}
          </>
        )}
      </div>
    </Tooltip>
  );
};

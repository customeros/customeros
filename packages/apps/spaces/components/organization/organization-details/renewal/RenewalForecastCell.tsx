import { useState } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { CurrencyInput } from '@ui/form/CurrencyInput/CurrencyInput';

import { useUpdateRenewalForecastMutation } from '@spaces/graphql';

const formatCurrency = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  minimumFractionDigits: 0,
}).format;

interface RenewalForecastCellProps {
  organizationId: string;
  currentForecast?: number | null;
  previousForecast?: number | null;
}

export const RenewalForecastCell = ({
  organizationId,
  currentForecast = null,
  previousForecast = null,
}: RenewalForecastCellProps) => {
  const [_value, setValue] = useState<string>(() =>
    String(currentForecast ?? ''),
  );
  const [updateRenewalForecast] = useUpdateRenewalForecastMutation();

  const handleChange = (value: string) => {
    setValue(value);
  };

  const handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
    updateRenewalForecast({
      variables: {
        input: {
          id: organizationId,
          amount: Number(_value),
        },
      },
    });
  };

  return (
    <Flex flexDir='column' justify='center'>
      {/* <Text fontSize='sm' color='gray.700'>
        {currentForecast ? formatCurrency(currentForecast) : '-'}
      </Text> */}
      <CurrencyInput
        value={_value}
        size='sm'
        onBlur={handleBlur}
        onChange={handleChange}
        isLabelVisible={false}
      />
      {previousForecast && (
        <Text fontSize='sm' color='gray.500' textDecoration='line-through'>
          {formatCurrency(previousForecast)}
        </Text>
      )}
    </Flex>
  );
};

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

interface RenewalForecastCellProps {
  amount?: number | null;
  potentialAmount?: number | null;
}

export const RenewalForecastCell = ({
  amount = null,
  potentialAmount = null,
}: RenewalForecastCellProps) => {
  const formattedAmount =
    amount !== null && amount >= 0 ? formatCurrency(amount, 0) : 'Unknown';
  const formattedPotentialAmount = formatCurrency(potentialAmount ?? 0, 0);

  const showPotentialAmount =
    amount !== null &&
    potentialAmount !== null &&
    formattedAmount !== formattedPotentialAmount;

  if (formattedAmount === 'Unknown')
    return <Text color='gray.400'>Unknown</Text>;

  return (
    <Flex flexDir='column' justify='center'>
      <Text fontSize='sm' color={amount ? 'gray.700' : 'gray.500'}>
        {formattedAmount}
      </Text>
      {showPotentialAmount && (
        <Text fontSize='sm' color='gray.500' textDecoration='line-through'>
          {formattedPotentialAmount}
        </Text>
      )}
    </Flex>
  );
};

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

interface RenewalForecastCellProps {
  amount?: number | null;
  isUpdatedByUser?: boolean;
  potentialAmount?: number | null;
}

export const RenewalForecastCell = ({
  amount = null,
  potentialAmount = null,
}: RenewalForecastCellProps) => {
  return (
    <Flex flexDir='column' justify='center'>
      <Text fontSize='sm' color={amount ? 'gray.700' : 'gray.500'}>
        {amount !== null && amount >= 0 ? formatCurrency(amount) : 'Unknown'}
      </Text>
      {potentialAmount && potentialAmount !== amount && (
        <Text fontSize='sm' color='gray.500' textDecoration='line-through'>
          {formatCurrency(potentialAmount)}
        </Text>
      )}
    </Flex>
  );
};

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

const formatCurrency = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  minimumFractionDigits: 0,
}).format;

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
        {amount ? formatCurrency(amount) : 'Unknown'}
      </Text>
      {potentialAmount && potentialAmount !== amount && (
        <Text fontSize='sm' color='gray.500' textDecoration='line-through'>
          {formatCurrency(potentialAmount)}
        </Text>
      )}
    </Flex>
  );
};

export function formatCurrency(
  amount: number,
  maximumFractionDigits: number = 2,
  currency: string = 'USD',
): string {
  const hasFraction = amount % 1 !== 0;
  const minimumFractionDigits = hasFraction ? 2 : 0;

  return Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits,
    maximumFractionDigits: maximumFractionDigits,
  }).format(amount);
}

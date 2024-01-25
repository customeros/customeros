export function formatCurrency(
  amount: number,
  maximumFractionDigits: number = 2,
  currency: string = 'USD',
): string {
  return Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: maximumFractionDigits,
  }).format(amount);
}

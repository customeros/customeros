export function formatCurrency(
  amount: number,
  maximumFractionDigits: number = 2,
): string {
  return Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 0,
    maximumFractionDigits: maximumFractionDigits,
  }).format(amount);
}

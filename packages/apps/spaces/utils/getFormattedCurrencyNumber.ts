function abbreviateNumber(num: number): string {
  const SI_SYMBOL = ['', 'K', 'M', 'B', 'T', 'Qa', 'Qi'];
  const tier = (Math.log10(Math.abs(num)) / 3) | 0;
  if (tier == 0) return num.toString();
  const suffix = SI_SYMBOL[tier];
  const scale = Math.pow(10, tier * 3);
  const scaled = num / scale;
  if (Math.round(scaled) === scaled) {
    return scaled + suffix;
  } else {
    return scaled.toFixed(1) + suffix;
  }
}
export function formatCurrency(amount: number): string {
  if (`${amount}`.length >= 6) {
    const abbreviatedNumber = abbreviateNumber(amount);
    return `$${abbreviatedNumber}`;
  } else {
    return Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
    }).format(amount);
  }
}

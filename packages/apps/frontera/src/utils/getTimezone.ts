export function getTimezone(timeZone: string): string {
  const options: Intl.DateTimeFormatOptions = {
    timeZone: timeZone,
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
  };

  const currentDate = new Date();

  const formatter = new Intl.DateTimeFormat('en-US', options);

  const localTime = formatter.format(currentDate);

  return localTime;
}

// Example usage
// console.log('Local time in New York:', getTimezone('America/New_York'));

export const handleFormatAmountCurrency = (val: string) => `$` + val ?? '';
export const handleParseAmountCurrency = (val: string) => val.replace(/^\$/, '');

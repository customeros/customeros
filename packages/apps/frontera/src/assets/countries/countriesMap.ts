import countries from '@assets/countries/countries.json';

export const countryMap = new Map(
  countries.map((country) => [country.alpha2.toLowerCase(), country.name]),
);

import { DateTimeUtils } from '@utils/date';

export const IssueDateCell = ({ value }: { value: string }) => {
  return (
    <span>
      {DateTimeUtils.format(value, DateTimeUtils.defaultFormatShortString)}
    </span>
  );
};

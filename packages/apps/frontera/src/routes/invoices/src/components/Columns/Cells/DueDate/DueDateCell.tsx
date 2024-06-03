import { DateTimeUtils } from '@utils/date';

export const DueDateCell = ({ value }: { value: string }) => {
  return (
    <span>
      {DateTimeUtils.format(value, DateTimeUtils.defaultFormatShortString)}
    </span>
  );
};

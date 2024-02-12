import ReactDatePicker, { ReactDatePickerProps } from 'react-datepicker';

import 'react-datepicker/dist/react-datepicker.css';

export const InlineDatePicker = (props: ReactDatePickerProps) => {
  return <ReactDatePicker inline {...props} />;
};

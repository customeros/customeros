import { extendTheme } from '@chakra-ui/react';

import { Input } from '@ui/form/Input/Input.theme';
import { Button } from '@ui/form/Button/Button.theme';
import { Avatar } from '@ui/media/Avatar/Avatar.theme';
import { radioTheme } from '@ui/form/Radio/Radio.theme';
import { Card } from '@ui/presentation/Card/Card.theme';
import { TagTheme } from '@ui/presentation/Tag/Tag.theme';
import { switchTheme } from '@ui/form/Switch/Switch.theme';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.theme';
import { Textarea } from '@ui/form/Textarea/Textarea.theme';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.theme';
import { NumberInput } from '@ui/form/NumberInput/NumberInput.theme';

import { colors } from './colors';
import { shadows } from './shadows';

export const theme = extendTheme({
  fonts: {
    heading: 'var(--font-barlow)',
    body: 'var(--font-barlow)',
    sticky: 'var(--font-merriweather)',
  },
  colors,
  shadows,
  components: {
    Avatar,
    Button,
    Tag: TagTheme,
    Card,
    Input,
    NumberInput,
    Textarea,
    Checkbox,
    Tooltip,
    Switch: switchTheme,
    Radio: radioTheme,
  },
  // styles: {
  //   global: () => ({
  //     // Optionally set global CSS styles
  //     body: {
  //       // '--chakra-colors-chakra-body-text': colors.gray['700'], // no idea how to change this variable
  //       // '.react-datepicker': {
  //       //   border: '0px',
  //       //   fontFamily: 'var(--font-barlow)',
  //       //   backgroundColor: 'white',
  //       // },
  //       // '.react-datepicker__navigation-icon::before': {
  //       //   w: '7px',
  //       //   h: '7px',
  //       //   top: '9px',
  //       //   borderColor: 'gray.500',
  //       //   borderWidth: '2px 2px 0 0',
  //       // },
  //       // '.react-datepicker__navigation-icon--previous::before': {
  //       //   right: '-5px',
  //       // },
  //       // '.react-datepicker__navigation-icon--next::before': {
  //       //   left: '-5px',
  //       // },
  //       // '.react-datepicker__header': {
  //       //   p: 0,
  //       //   borderBottom: 'none',
  //       //   backgroundColor: 'white',
  //       // },
  //       // '.react-datepicker__month-container': {
  //       //   mt: '7px',
  //       //   width: '100%',
  //       //   height: '100%',
  //       // },
  //       // '.react-datepicker__current-month': {
  //       //   fontSize: '16px',
  //       //   fontWeight: '600',
  //       //   color: 'gray.700',
  //       // },
  //       // '.react-datepicker__month': {
  //       //   margin: '0',
  //       //   backgroundColor: 'white',
  //       // },
  //       // '.react-datepicker__day-names': {
  //       //   width: '100%',
  //       //   display: 'flex',
  //       //   mt: '3',
  //       //   mb: '1',
  //       //   borderBottom: 'none',
  //       //   justifyContent: 'space-between',
  //       //   backgroundColor: 'white',
  //       //   padding: '0',
  //       // },
  //       // '.react-datepicker__day-name': {
  //       //   w: '40px',
  //       //   mx: '8px',
  //       //   my: '10px',
  //       //   margin: '0',
  //       //   fontSize: '14px',
  //       //   fontWeight: '600',
  //       //   color: 'gray.700',
  //       // },
  //       // '.react-datepicker__week': {
  //       //   width: '100%',
  //       //   display: 'flex',
  //       //   marginBottom: '1',
  //       //   justifyContent: 'space-between',
  //       // },
  //       // '.react-datepicker__week:last-child': {
  //       //   marginBottom: '0',
  //       // },
  //       // '.react-datepicker__day': {
  //       //   width: '40px',
  //       //   margin: '0px',
  //       //   height: '40px',
  //       //   display: 'flex',
  //       //   fontSize: '14px',
  //       //   fontWeight: '400',
  //       //   alignItems: 'center',
  //       //   justifyContent: 'center',
  //       //   color: 'gray.700',
  //       // },
  //       // '.react-datepicker__day:hover': {
  //       //   bg: 'primary.50',
  //       //   border: '1px solid',
  //       //   borderColor: 'primary.200',
  //       //   borderRadius: 'full',
  //       // },
  //       // '.react-datepicker__day:focus-visible': {
  //       //   boxShadow: 'outline',
  //       //   outline: 'none',
  //       //   borderRadius: 'full',
  //       // },
  //       // '.react-datepicker__day--disabled': {
  //       //   color: 'gray.400',
  //       //   border: 'unset',
  //       //   backgroundColor: 'white',
  //       // },
  //       // '.react-datepicker__day--disabled:hover': {
  //       //   color: 'gray.400',
  //       //   border: 'unset',
  //       //   backgroundColor: 'white',
  //       // },
  //       // '.react-datepicker__day--outside-month': {
  //       //   color: 'gray.400',
  //       // },
  //       // '.react-datepicker__day--today': {
  //       //   border: 'unset',
  //       //   color: 'primary.500',
  //       // },
  //       // '.react-datepicker__day--selected': {
  //       //   bg: 'primary.50',
  //       //   border: '1px solid',
  //       //   borderColor: 'primary.200',
  //       //   borderRadius: 'full',
  //       // },
  //       // '.react-datepicker__day--selected:hover': {
  //       //   backgroundColor: 'primary.200',
  //       // },
  //       // '.react-datepicker__day--selected:focusVisible': {
  //       //   boxShadow: 'outline',
  //       //   outline: 'none',
  //       // },
  //       // '.react-datepicker__day--disabled.react-datepicker__day--today': {
  //       //   bg: 'gray.50',
  //       // },
  //       // '.react-datepicker__day--keyboard-selected': {
  //       //   bg: 'primary.50',
  //       //   border: '1px solid',
  //       //   borderColor: 'primary.200',
  //       //   borderRadius: 'full',
  //       // },
  //       // '.react-datepicker__day--keyboard-selected, .react-datepicker__day--keyboard-selected:focusVisible':
  //       //   {
  //       //     outline: 'none',
  //       //     bg: 'primary.50',
  //       //     border: '1px solid',
  //       //     borderColor: 'red.200',
  //       //     borderRadius: 'full',
  //       //   },
  //       // '.react-datepicker__day--keyboard-selected.react-datepicker__day--today':
  //       //   {
  //       //     color: 'primary.500',
  //       //     bg: 'primary.50',
  //       //     outline: 'none',
  //       //     boxShadow: 'outline',
  //       //   },
  //       // '.react-datepicker__day--keyboard-selected:hover': {
  //       //   backgroundColor: 'primary.50',
  //       //   boxShadow: 'outline',
  //       // },
  //       // '.react-datepicker-popper': {
  //       //   paddingTop: '4px',
  //       // },
  //       // '.hidden': {
  //       //   display: 'none',
  //       // },
  //       // '.selected-month-year-button p': {
  //       //   color: 'white',
  //       // },
  //     },
  //   }),
  // },
});

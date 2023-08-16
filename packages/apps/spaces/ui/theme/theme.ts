import { extendTheme } from '@chakra-ui/react';

import { Avatar } from '@ui/media/Avatar/Avatar.theme';
import { Button } from '@ui/form/Button/Button.theme';
import { Input } from '@ui/form/Input/Input.theme';
import { NumberInput } from '@ui/form/NumberInput/NumberInput.theme';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.theme';
import { Textarea } from '@ui/form/Textarea/Textarea.theme';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.theme';

import { colors } from './colors';
import { shadows } from './shadows';
import {Simulate} from "react-dom/test-utils";
import select = Simulate.select;
import {fontWeight} from "@mui/system";

export const theme = extendTheme({
  fonts: {
    heading: 'var(--font-barlow)',
    body: 'var(--font-barlow)',
  },
  colors,
  shadows,
  components: {
    Avatar,
    Button,
    Input,
    NumberInput,
    Textarea,
    Checkbox,
    Tooltip,
  },
});

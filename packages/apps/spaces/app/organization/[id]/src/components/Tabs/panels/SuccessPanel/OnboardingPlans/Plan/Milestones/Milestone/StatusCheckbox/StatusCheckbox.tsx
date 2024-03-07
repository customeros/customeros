import { Checkbox, CheckboxProps } from '@ui/form/Checkbox';
import { ArrowNarrowDownRight } from '@ui/media/icons/ArrowNarrowDownRight';

interface StatusCheckboxProps extends CheckboxProps {
  readOnly?: boolean;
  showCustomIcon?: boolean;
}

export const StatusCheckbox = ({
  showCustomIcon,
  ...props
}: StatusCheckboxProps) => {
  const customeIconColor = (() => {
    switch (props.colorScheme) {
      case 'gray':
        return 'gray.400';
      case 'warning':
        return 'warning.600';
      case 'success':
        return 'success.500';
      default:
        return 'gray.400';
    }
  })();

  return (
    <Checkbox
      mr='2'
      size='md'
      icon={
        showCustomIcon ? (
          <ArrowNarrowDownRight boxSize='3.5' color={customeIconColor} />
        ) : undefined
      }
      {...props}
      sx={{
        '& > span': {
          bg: `${props.colorScheme}.50`,
          borderColor: customeIconColor,
          _hover: {
            bg: `${props.colorScheme}.100`,
          },
          _focus: {
            borderColor: customeIconColor,
            bg: `${props.colorScheme}.100`,
          },
          '& > svg': {
            color: customeIconColor,
          },
        },
      }}
    />
  );
};

import { cn } from '@ui/utils/cn';
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
        return 'text-gray-400';
      case 'warning':
        return 'text-warning-600';
      case 'success':
        return 'text-success-500';
      default:
        return 'text-gray-400';
    }
  })();

  return (
    <Checkbox
      mr='2'
      size='md'
      icon={
        showCustomIcon ? (
          <ArrowNarrowDownRight
            className={cn(customeIconColor, 'size-[14px]')}
          />
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

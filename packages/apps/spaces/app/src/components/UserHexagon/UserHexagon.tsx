import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { genRandomNum } from '@shared/util/math';

interface UserHexagonProps {
  name: string;
}

export const UserHexagon = ({ name }: UserHexagonProps) => {
  const colors = useToken('colors', [
    'error.500',
    'rose.500',
    'orange.700',
    'orange.400',
    'warning.400',
    'yellow.400',
    'moss.400',
    'green.500',
    'greenLight.700',
    'cyan.400',
    'teal.400',
    'blue.500',
    'blueDark.700',
    'primary.500',
    'violet.600',
    'fuchsia.500',
    'fuchsia.700',
    'pink.500',
    'gray.400',
    'gray.500',
    'gray.950',
  ]);

  const randomColor = colors[genRandomNum(colors.length - 1)];

  return (
    <Flex w='26px' h='28px' align='center' justify='center'>
      <svg
        xmlns='http://www.w3.org/2000/svg'
        width='26'
        height='28'
        viewBox='0 0 26 28'
        fill='none'
        color={randomColor}
        style={{
          position: 'absolute',
        }}
      >
        <path
          d='M11.25 1.58771C12.3329 0.962498 13.6671 0.962498 14.75 1.58771L22.8744 6.27831C23.9573 6.90353 24.6244 8.05897 24.6244 9.3094V18.6906C24.6244 19.941 23.9573 21.0965 22.8744 21.7217L14.75 26.4123C13.6671 27.0375 12.3329 27.0375 11.25 26.4123L3.12564 21.7217C2.04274 21.0965 1.37564 19.941 1.37564 18.6906V9.3094C1.37564 8.05897 2.04274 6.90353 3.12564 6.27831L11.25 1.58771Z'
          fill='#FCFCFD'
          stroke='currentColor'
        />
      </svg>

      <Text fontSize='sm' color={randomColor} zIndex={2}>
        {getInitials(name)}
      </Text>
    </Flex>
  );
};

function getInitials(name: string) {
  const temp = name.toUpperCase().split(' ').splice(0, 2);

  return temp
    .map((s) => s[0])
    .join('')
    .trim();
}

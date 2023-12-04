'use client';
import { ReactNode, PropsWithChildren } from 'react';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { Card, CardBody, CardProps, CardHeader } from '@ui/presentation/Card';

import { HelpButton } from '../HelpButton';

interface ChartCardProps extends CardProps {
  stat?: string;
  title: string;
  hasData?: boolean;
  renderSubStat?: () => ReactNode;
  renderHelpContent?: () => ReactNode;
}

export const ChartCard = ({
  stat,
  title,
  hasData,
  children,
  renderSubStat,
  renderHelpContent,
  ...props
}: PropsWithChildren<ChartCardProps>) => {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <>
      <Card
        borderRadius='lg'
        boxShadow='xs'
        border='1px solid'
        borderColor='gray.200'
        _hover={{
          '& #help-button': {
            visibility: 'visible',
          },
        }}
        {...props}
      >
        <CardHeader pb='0' pt='4' px='6'>
          <Flex gap='2' align='center'>
            <Text fontSize='lg' fontWeight='normal'>
              {title}
            </Text>
            {renderHelpContent && (
              <HelpButton isOpen={isOpen} onOpen={onOpen} />
            )}
          </Flex>
          {stat && (
            <Heading
              fontSize={hasData ? undefined : '18px'}
              color={hasData ? 'gray.700' : 'gray.400'}
            >
              {hasData ? stat : 'No data yet'}
            </Heading>
          )}
          {hasData && renderSubStat && renderSubStat?.()}
        </CardHeader>
        <CardBody px='6' pb='6'>
          {children}
        </CardBody>
      </Card>

      <InfoDialog
        label={title}
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={onClose}
        confirmButtonLabel='Got it'
      >
        {renderHelpContent?.()}
      </InfoDialog>
    </>
  );
};

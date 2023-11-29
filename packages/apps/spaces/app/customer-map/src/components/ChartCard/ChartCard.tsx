'use client';
import { ReactNode, PropsWithChildren } from 'react';

import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { Card, CardBody, CardProps, CardHeader } from '@ui/presentation/Card';

interface ChartCardProps extends CardProps {
  stat?: string;
  title: string;
  renderSubStat?: () => ReactNode;
}

export const ChartCard = ({
  stat,
  title,
  children,
  renderSubStat,
  ...props
}: PropsWithChildren<ChartCardProps>) => {
  return (
    <Card
      borderRadius='lg'
      boxShadow='xs'
      border='1px solid'
      borderColor='gray.200'
      {...props}
    >
      <CardHeader pb='0' pt='4' px='6'>
        <Text fontSize='lg' fontWeight='normal'>
          {title}
        </Text>
        {stat && <Heading>{stat}</Heading>}
        {renderSubStat && renderSubStat?.()}
      </CardHeader>
      <CardBody px='6' pb='6'>
        {children}
      </CardBody>
    </Card>
  );
};

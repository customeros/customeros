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
    <Card {...props}>
      <CardHeader minH='8.24rem'>
        <Text fontSize='lg' fontWeight='normal'>
          {title}
        </Text>
        {stat && <Heading>{stat}</Heading>}
        {renderSubStat && renderSubStat?.()}
      </CardHeader>
      <CardBody>{children}</CardBody>
    </Card>
  );
};

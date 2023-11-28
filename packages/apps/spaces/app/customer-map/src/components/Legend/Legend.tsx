'use client';

import { scaleOrdinal } from '@visx/scale';
import { LegendItem, LegendLabel, LegendOrdinal } from '@visx/legend';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

interface LegendProps {
  data: {
    label: string;
    color: string;
    borderColor?: string;
  }[];
}

export const Legend = ({ data }: LegendProps) => {
  const scale = scaleOrdinal({
    domain: data.map((d) => d.label),
    range: data.map((d) => d.color),
  });

  return (
    <LegendOrdinal scale={scale}>
      {(labels) => (
        <Flex justify='flex-end'>
          {labels.map((label, i) => (
            <LegendItem key={`legend-quantile-${i}`} margin='0 0.5rem'>
              <svg
                width={data[i].borderColor ? 9 : 8}
                height={data[i].borderColor ? 9 : 8}
                style={{ marginRight: '0.25rem' }}
              >
                <circle
                  fill={label.value}
                  cx={data[i].borderColor ? 4.5 : 4}
                  cy={data[i].borderColor ? 4.5 : 4}
                  r={data[i].borderColor ? 4 : 4}
                  strokeWidth={data[i].borderColor ? 1 : 0}
                  stroke={data[i].borderColor ?? 'transparent'}
                />
              </svg>
              <LegendLabel align='left' margin='0 0 0 4px'>
                <Text fontSize='sm'>{label.text}</Text>
              </LegendLabel>
            </LegendItem>
          ))}
        </Flex>
      )}
    </LegendOrdinal>
  );
};

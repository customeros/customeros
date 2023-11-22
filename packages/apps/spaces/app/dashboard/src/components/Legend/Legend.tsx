'use client';

import { scaleOrdinal } from '@visx/scale';
import { LegendItem, LegendLabel, LegendOrdinal } from '@visx/legend';

import { Flex } from '@ui/layout/Flex';

interface LegendProps {
  data: {
    label: string;
    color: string;
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
            <LegendItem key={`legend-quantile-${i}`} margin='0 1rem'>
              <svg width={12} height={12} style={{ marginRight: '0.5rem' }}>
                <circle fill={label.value} cx={6} cy={6} r={6} />
              </svg>
              <LegendLabel align='left' margin='0 0 0 4px'>
                {label.text}
              </LegendLabel>
            </LegendItem>
          ))}
        </Flex>
      )}
    </LegendOrdinal>
  );
};

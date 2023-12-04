'use client';

import { Pie } from '@visx/shape';
import { Group } from '@visx/group';

import { useToken } from '@ui/utils';

import { AnimatedPie } from './AnimatedPie';

export type RevenueAtRiskDatum = {
  atRisk: number;
  highConfidence: number;
};

const mockData: RevenueAtRiskDatum = {
  atRisk: 1200,
  highConfidence: 820,
};

interface RevenueAtRiskProps {
  width: number;
  height: number;
  hasContracts?: boolean;
  data: RevenueAtRiskDatum;
}

const margin = { top: 0, right: 0, bottom: 0, left: 0 };

const mapData = (data: RevenueAtRiskDatum) => {
  return [
    {
      label: 'At Risk',
      value: data.atRisk,
    },
    {
      label: 'High Confidence',
      value: data.highConfidence,
    },
  ];
};

const RevenueAtRisk = ({
  width,
  height,
  data: _data,
  hasContracts,
}: RevenueAtRiskProps) => {
  const data = hasContracts ? _data : mockData;
  const [warning300, greenLight500] = useToken('colors', [
    hasContracts ? 'warning.300' : 'gray.100',
    hasContracts ? 'greenLight.500' : 'gray.300',
  ]);
  const mappedData = mapData(data);

  const colorScale = {
    Confidence: greenLight500,
    Risk: warning300,
  };

  const innerWidth = width - margin.left - margin.right;
  const innerHeight = height - margin.top - margin.bottom;
  const radius = Math.min(innerWidth, innerHeight) / 2;
  const centerY = innerHeight / 2;
  const centerX = innerWidth / 2;
  const donutThickness = 70;

  if (width < 10) return null;

  return (
    <svg width={width} height={height}>
      <Group top={centerY + margin.top} left={centerX + margin.left}>
        <Pie
          data={mappedData}
          pieValue={(d) => d.value}
          outerRadius={radius}
          innerRadius={radius - donutThickness}
          cornerRadius={8}
          padAngle={0.01}
        >
          {(pie) => {
            return (
              <AnimatedPie
                {...pie}
                animate
                getKey={(arc) => arc.data.label}
                getColor={(arc) =>
                  arc.data.label === 'At Risk'
                    ? colorScale.Risk
                    : colorScale.Confidence
                }
              />
            );
          }}
        </Pie>
      </Group>
    </svg>
  );
};

export default RevenueAtRisk;

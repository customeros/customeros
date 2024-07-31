import { Pie } from '@visx/shape';
import { Group } from '@visx/group';

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

const RevenueAtRiskChart = ({
  width,
  height,
  data: _data,
  hasContracts,
}: RevenueAtRiskProps) => {
  const data = hasContracts ? _data : mockData;

  const colors = {
    warning300: hasContracts ? '#FEC84B' : '#F2F4F7',
    greenLight500: hasContracts ? '#66C61C' : '#D0D5DD',
  };
  const mappedData = mapData(data);

  const colorScale = {
    Confidence: colors.greenLight500,
    Risk: colors.warning300,
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
          padAngle={0.01}
          cornerRadius={8}
          data={mappedData}
          outerRadius={radius}
          pieValue={(d) => d.value}
          innerRadius={radius - donutThickness}
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

export default RevenueAtRiskChart;

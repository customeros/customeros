'use client';

import { Pie } from '@visx/shape';
import { Group } from '@visx/group';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { useToken } from '@ui/utils';

import { AnimatedPie } from './AnimatedPie';

const mockData = [
  {
    label: 'At Risk',
    value: 355300,
  },
  {
    label: 'High Confidence',
    value: 1504990,
  },
];

interface RevenueAtRiskProps {}

const margin = { top: 0, right: 0, bottom: 0, left: 0 };

const RevenueAtRisk = (_props: RevenueAtRiskProps) => {
  const [yellow400] = useToken('colors', ['yellow.400']);

  const colorScale = {
    Confidence: '#66C61C',
    Risk: yellow400,
  };

  return (
    <ParentSize>
      {({ width, height }) => {
        if (width < 10) return null;

        const innerWidth = width - margin.left - margin.right;
        const innerHeight = height - margin.top - margin.bottom;
        const radius = Math.min(innerWidth, innerHeight) / 2;
        const centerY = innerHeight / 2;
        const centerX = innerWidth / 2;
        const donutThickness = 70;

        return (
          <svg width={width} height={height}>
            <Group top={centerY + margin.top} left={centerX + margin.left}>
              <Pie
                data={mockData}
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
      }}
    </ParentSize>
  );
};

export default RevenueAtRisk;

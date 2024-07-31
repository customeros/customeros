import { PatternLines } from '@visx/pattern';
import {
  XYChart,
  Tooltip,
  AnimatedAxis,
  AnimatedGrid,
  AnimatedBarStack,
  AnimatedBarSeries,
} from '@visx/xychart';

import { cn } from '@ui/utils/cn';

import { mockData } from './mock';
import { Legend } from '../../Legend';
import { getMonthLabel } from '../util';

export type RetentionRateDatum = {
  month: number;
  values: {
    renewed: number;
    churned: number;
  };
};
interface RetentionRateProps {
  width?: number;
  height?: number;
  hasContracts?: boolean;
  data: RetentionRateDatum[];
}

const getX = (d: RetentionRateDatum) => getMonthLabel(d.month);

const RetentionRateChart = ({
  width,
  data: _data,
  hasContracts,
}: RetentionRateProps) => {
  const data = hasContracts ? _data : mockData;

  const colors = {
    gray700: '#344054',
    warning950: hasContracts ? '#4E1D09' : '#D0D5DD',
    greenLight500: hasContracts ? '#66C61C' : '#EAECF0',
  };

  const colorScale = {
    Renewed: colors.greenLight500,
    Churned: colors.warning950,
  };

  const isMissingData = (dataPoint: 'renewed' | 'churned') =>
    data.every((d) => d.values[dataPoint] === 0);

  const legendData = [
    {
      label: 'Renewed',
      color: colorScale.Renewed,
      isMissingData: isMissingData('renewed'),
    },
    {
      label: 'Churned',
      color: colorScale.Churned,
      isMissingData: isMissingData('churned'),
    },
  ];

  return (
    <>
      <Legend data={legendData} />
      <XYChart
        height={200}
        width={width || 500}
        yScale={{ type: 'linear' }}
        margin={{ top: 12, right: 0, bottom: 20, left: 0 }}
        xScale={{
          type: 'band',
          paddingInner: 0.4,
          paddingOuter: 0.4,
        }}
      >
        <PatternLines
          width={8}
          height={8}
          strokeWidth={2}
          id='stripes-renewed'
          orientation={['diagonal']}
          stroke={colorScale.Renewed}
        />
        <PatternLines
          width={8}
          height={8}
          strokeWidth={2}
          id='stripes-churned'
          orientation={['diagonal']}
          stroke={colorScale.Churned}
        />
        <AnimatedBarStack>
          <AnimatedBarSeries
            radiusTop
            radius={4}
            data={data}
            dataKey='Renewed'
            yAccessor={(d) => d.values.renewed}
            xAccessor={(d) => getMonthLabel(d.month)}
            colorAccessor={(_, i) =>
              i === data.length - 1
                ? 'url(#stripes-renewed)'
                : colorScale.Renewed
            }
          />
          <AnimatedBarSeries
            radius={4}
            data={data}
            radiusBottom
            dataKey='Churned'
            yAccessor={(d) => -d.values.churned}
            xAccessor={(d) => getMonthLabel(d.month)}
            colorAccessor={(_, i) =>
              i === data.length - 1
                ? 'url(#stripes-churned)'
                : colorScale.Churned
            }
          />
        </AnimatedBarStack>

        <AnimatedAxis
          hideTicks
          hideAxisLine
          orientation='bottom'
          tickLabelProps={{
            fontSize: 12,
            fontWeight: 'medium',
            fontFamily: `var(--font-ibm-plex-sans)`,
          }}
        />
        <AnimatedGrid
          numTicks={1}
          columns={false}
          lineStyle={{ stroke: 'white', strokeWidth: 2 }}
        />
        <Tooltip
          snapTooltipToDatumY
          snapTooltipToDatumX
          style={{
            position: 'absolute',
            padding: '8px 12px',
            background: colors.gray700,
            borderRadius: '8px',
          }}
          renderTooltip={({ tooltipData }) => {
            const xLabel = getX(
              tooltipData?.nearestDatum?.datum as RetentionRateDatum,
            );
            const values = (
              tooltipData?.nearestDatum?.datum as RetentionRateDatum
            ).values;

            return (
              <div className='flex flex-col'>
                {hasContracts ? (
                  <>
                    <p className='text-white font-semibold text-sm'>{xLabel}</p>

                    <div className='flex flex-col'>
                      <TooltipEntry
                        label='Renewed'
                        value={values.renewed}
                        color={colorScale.Renewed}
                        isMissingData={isMissingData('renewed')}
                      />
                      <TooltipEntry
                        label='Churned'
                        value={values.churned}
                        color={colorScale.Churned}
                        isMissingData={isMissingData('churned')}
                      />
                    </div>
                  </>
                ) : (
                  <p className='text-white font-semibold text-sm'>
                    No data yet
                  </p>
                )}
              </div>
            );
          }}
        />
      </XYChart>
      <p className='text-gray-500 text-xs mt-2'>
        <i>*Key data missing.</i>
      </p>
    </>
  );
};

const TooltipEntry = ({
  color,
  label,
  value,
  isMissingData,
}: {
  color: string;
  label: string;
  value: number;
  isMissingData?: boolean;
}) => {
  return (
    <div className='flex items-center gap-4'>
      <div className='flex items-center flex-1 gap-2'>
        <div
          style={{ backgroundColor: color }}
          className='flex w-2 h-2  rounded-full border border-white'
        />
        <p className='text-white text-sm'>{label}</p>
      </div>
      <div className='flex justify-start'>
        <p
          className={cn(
            isMissingData ? 'text-gray-400' : 'text-white',
            'text-sm',
          )}
        >
          {isMissingData ? '*' : value}
        </p>
      </div>
    </div>
  );
};

export default RetentionRateChart;

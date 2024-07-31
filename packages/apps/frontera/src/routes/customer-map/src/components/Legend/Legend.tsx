import { scaleOrdinal } from '@visx/scale';
import { LegendItem, LegendLabel, LegendOrdinal } from '@visx/legend';

interface LegendProps {
  leftElement?: React.ReactNode;
  data: {
    label: string;
    color: string;
    borderColor?: string;
    isMissingData?: boolean;
  }[];
}

export const Legend = ({ data, leftElement }: LegendProps) => {
  const scale = scaleOrdinal({
    domain: data.map((d) => d.label),
    range: data.map((d) => d.color),
  });

  return (
    <LegendOrdinal scale={scale}>
      {(labels) => (
        <div className='flex justify-between'>
          <div>{leftElement}</div>
          <div className='flex'>
            {labels.map((label, i) => (
              <LegendItem margin='0 0.5rem' key={`legend-quantile-${i}`}>
                <svg
                  style={{ marginRight: '0.25rem' }}
                  width={data[i].borderColor ? 9 : 8}
                  height={data[i].borderColor ? 9 : 8}
                >
                  <circle
                    fill={label.value}
                    r={data[i].borderColor ? 4 : 4}
                    cx={data[i].borderColor ? 4.5 : 4}
                    cy={data[i].borderColor ? 4.5 : 4}
                    strokeWidth={data[i].borderColor ? 1 : 0}
                    stroke={data[i].borderColor ?? 'transparent'}
                  />
                </svg>
                <LegendLabel align='left' margin='0 0 0 4px'>
                  <p className='text-sm'>
                    {label.text}
                    {data[i].isMissingData && <span>*</span>}
                  </p>
                </LegendLabel>
              </LegendItem>
            ))}
          </div>
        </div>
      )}
    </LegendOrdinal>
  );
};

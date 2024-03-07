import * as d3 from 'd3';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { dodge } from './dodge';

type CircleData<Datum> = {
  x: Date;
  y: number;
  r: number;
  values: Datum;
};
export type DodgedCircleData<Datum> = {
  x: number;
  y: number;
  r: number;
  data: CircleData<Datum>;
};

const LOW_R = 4; // min radius value
const TOP_R = 60; // max radius value
const PADDING = 6; // padding between circles

interface UseBubblesOptions<Datum> {
  width: number;
  height: number;
  marginLeft: number;
  marginRight: number;
  data: CircleData<Datum>[];
}

export const useDodge = <Datum>({
  width,
  data,
  marginLeft,
  marginRight,
}: UseBubblesOptions<Datum>) => {
  const minMaxX = d3.extent(data, (d) => d.x) as [Date, Date];
  const minMaxR = d3.extent(data, (d) => d.r) as [number, number];

  const isTopRadiusDecreased = useFeatureIsOn('decrease-circle-top-radius');

  const xScale = d3
    .scaleTime()
    .domain(minMaxX)
    .range([marginLeft, width + marginRight]);

  const rScale = (() => {
    if (isTopRadiusDecreased) {
      return d3.scaleLinear().domain(minMaxR).range([LOW_R, 20]).nice();
    }

    return d3.scaleSqrt().domain(minMaxR).range([LOW_R, TOP_R]).nice();
  })();

  function transformData(input: CircleData<Datum>[]) {
    return dodge(input, {
      x: (d: CircleData<Datum>) => xScale(d.x),
      r: (d: CircleData<Datum>) => rScale(d.r),
      padding: PADDING,
    }) as DodgedCircleData<Datum>[];
  }

  return { transformData, minMaxX };
};

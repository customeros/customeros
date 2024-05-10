import sample from 'lodash/sample';
import genRandomNormalPoints, {
  PointsRange,
} from '@visx/mock-data/lib/generators/genRandomNormalPoints';

const genTimeValue = () =>
  sample([
    new Date('2023-01-01'),
    new Date('2023-01-15'),
    new Date('2023-01-25'),
    new Date('2023-02-01'),
    new Date('2023-02-15'),
    new Date('2023-02-25'),
    new Date('2023-03-01'),
    new Date('2023-03-15'),
    new Date('2023-03-25'),
    new Date('2023-04-01'),
    new Date('2023-04-15'),
    new Date('2023-04-25'),
    new Date('2023-05-01'),
    new Date('2023-05-15'),
    new Date('2023-05-25'),
    new Date('2023-06-01'),
    new Date('2023-06-15'),
    new Date('2023-06-25'),
    new Date('2023-07-02'),
    new Date('2023-07-15'),
    new Date('2023-07-25'),
    new Date('2023-08-01'),
    new Date('2023-08-15'),
    new Date('2023-08-25'),
    new Date('2023-09-01'),
    new Date('2023-09-15'),
    new Date('2023-09-25'),
    new Date('2023-10-01'),
    new Date('2023-10-15'),
    new Date('2023-10-25'),
    new Date('2023-11-01'),
    new Date('2023-11-15'),
    new Date('2023-11-25'),
    new Date('2023-12-01'),
    new Date('2023-12-15'),
    new Date('2023-12-31'),
  ]);

const points: PointsRange[] = genRandomNormalPoints(70, 0.5).filter(
  (_, i) => i < 70,
);

export const mockData = points.map(() => ({
  x: genTimeValue(),
  y: 0,
  r: sample([2, 15, 6, 12, 22, 5, 10, 15, 25, 7, 120]),
  values: {
    id: `${Math.random()}`,
    name: '',
    status: sample(['OK', 'AT_RISK', 'CHURNED']),
  },
}));

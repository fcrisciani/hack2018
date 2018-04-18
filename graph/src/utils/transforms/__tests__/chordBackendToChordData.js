import chordBackendToChordData from '../chordBackendToChordData';
import sample from './backendChordSample.json';

it('data transform', () => {
  expect(chordBackendToChordData(sample.graph)).toMatchSnapshot();
})
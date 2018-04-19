import React, { Component } from 'react';
import { ResponsiveChord } from 'nivo';
import request from 'superagent';
import chordBackendToChordData from './utils/transforms/chordBackendToChordData';
import HistorySelect from './appComponents/HistorySelect';
import PlayPause from './appComponents/PlayPause';
import sample from './mock/sample.json';
import './App.css';
import debounce from 'lodash/debounce'

const ENDPOINT_LOCALSTORAGE_KEY = 'hkdayendpoint';
let tickCount = 1;
const sum = (a, b) => {
  return a + b;
}
const getTotalCount = (matrix) => {
  return matrix.map((vector) => {
    return vector.reduce(sum, 0)
  }).reduce(sum, 0);
}
class App extends Component {
  constructor(props) {
    super(props);
    this.state = ({
      endpoint:localStorage.getItem(ENDPOINT_LOCALSTORAGE_KEY) || 'http://52.42.55.249:10000/chord',
      hasPauseCursor: false,
      playing: true,
      history: [],
      initialFetching: true,
      data: [[]],
      labels:[],
    })
  }
  tick = debounce(() => {

    this.doRequest();
    clearTimeout(this.timeout);
    const TICK_PERIOD = 1000;
    const timeoutT =  Math.min(TICK_PERIOD * tickCount++,10000)
    this.timeout = setTimeout(this.tick, timeoutT);
  }, 500)
  doRequest() {
    const setNewState = (newState) => {
      const time = new Date();
      const totalCount = getTotalCount(newState.matrix);
      this.state.history.push({
        time,
        item: newState,
        totalCount,
      });
      if (this.state.history.length > 100) {
        // cap history
        this.state.history.splice(0, 1);
      }
      if (this.state.playing) {
        this.setState({
          history: this.state.history.slice(),
          time,
          initialFetching: false,
          data: newState.matrix,
          labels: newState.labels,
          totalCount,
        })
      } else {
        this.setState({
          history: this.state.history.slice(),
        })
      }
    }
    const useSample = false;
    if (useSample) {
      const newState = chordBackendToChordData(sample.graph);
      setNewState(newState);
    } else {
      request
        .get(this.state.endpoint)
        .end((err, res) => {
          if(err) {
            console.log('error', err);
            return;
          }
          const graph = JSON.parse(res.text);
          const newState = chordBackendToChordData(graph.graph)
          setNewState(newState);
        })
    }
  }
  onHistoryChange = (historyItem) => {
    this.setState({
      playing: false,
      hasPauseCursor: true,
      time: historyItem.time,
      data: historyItem.item.matrix,
      labels: historyItem.item.labels,
      totalCount: historyItem.totalCount,
    })
  }
  componentDidMount() {
    clearTimeout(this.timeout);
    this.tick();
  }
  onPlayPauseChange = (newVal) => {
    this.setState({
      playing: newVal,
      hasPauseCursor: !newVal
    })
    if (newVal === true) {
      this.tick();
    }
  }
  storeEndpoint = debounce((endpoint) => {
    localStorage.setItem(ENDPOINT_LOCALSTORAGE_KEY, endpoint)
  }, 200)
  onInputChange = (e) => {
    const endpoint = e.target.value;
    this.setState({
      endpoint,
    })
    this.storeEndpoint(endpoint)
    this.tick();
  }
  render() {
    if (this.state.initialFetching) {
      return(
        <div className="loading">
          loading...
        </div>
      );
    }
    return (
      <div className="App">
        <label>
          endpoint
          <input
            style={{
              border: 0,
              width: 400,
              margin: '0 5px',
            }}
           onChange={this.onInputChange} value={this.state.endpoint} type="text"/>
        </label>
        <h3 className="title">
          <PlayPause onChange={this.onPlayPauseChange} playing={this.state.playing}/>
          {this.state.totalCount} connections @{this.state.time.toLocaleTimeString('en-US')}
        </h3>
        <ResponsiveChord
          matrix={this.state.data}
          keys={this.state.labels}
          margin={{
              "top": 60,
              "right": 60,
              "bottom": 90,
              "left": 60
          }}
          padAngle={0.02}
          innerRadiusRatio={0.96}
          innerRadiusOffset={0.02}
          arcOpacity={1}
          arcBorderWidth={1}
          arcBorderColor="inherit:darker(0.4)"
          ribbonOpacity={0.5}
          ribbonBorderWidth={1}
          ribbonBorderColor="inherit:darker(0.4)"
          enableLabel={true}
          label="id"
          labelOffset={12}
          labelRotation={-90}
          labelTextColor="inherit:darker(1)"
          colors="set1"
          isInteractive={true}
          arcHoverOpacity={1}
          arcHoverOthersOpacity={0.25}
          ribbonHoverOpacity={0.75}
          ribbonHoverOthersOpacity={0.25}
          animate={true}
          motionStiffness={90}
          motionDamping={17}
          legends={[{
              "anchor": "bottom",
              "direction": "row",
              "translateY": 70,
              "itemWidth": 80,
              "itemHeight": 14,
              "symbolSize": 14,
              "symbolShape": "circle"
            }]}
          />
            <div className="controls">
          <HistorySelect
              hasPauseCursor={this.state.hasPauseCursor}
              history={this.state.history}
              onChange={this.onHistoryChange}/>
        </div>
      </div>
    );
  }
}

export default App;

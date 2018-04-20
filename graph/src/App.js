import React, { Component } from 'react';
import { Chord } from 'nivo';
import request from 'superagent';
import chordBackendToChordData from './utils/transforms/chordBackendToChordData';
import HistorySelect from './appComponents/HistorySelect';
import PlayPause from './appComponents/PlayPause';
import List from './appComponents/List';
import sample from './mock/sample.json';
import samplePods from './mock/sample_pods.json';
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
      windowWidth: window.innerWidth,
      windowHeight: window.innerHeight,
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
    const setNewState = ({connected, disconnected}) => {
      const time = new Date();
      const totalCount = getTotalCount(connected.matrix);
      this.state.history.push({
        time,
        item: connected,
        totalCount,
        disconnectedLabels: disconnected.labels,
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
          data: connected.matrix,
          labels: connected.labels,
          totalCount,
          disconnectedLabels: disconnected.labels,
        })
      } else {
        this.setState({
          history: this.state.history.slice(),
        })
      }
    }
    const useSample = false;
    const useSamplePods = true;
    if (useSample) {
      const {connected, disconnected } = chordBackendToChordData(sample.graph);
      setNewState({connected, disconnected});
    } else if (useSamplePods) {
      const {connected, disconnected } = chordBackendToChordData(samplePods.graph);
      setNewState({connected, disconnected});
    } else {
      request
        .get(this.state.endpoint)
        .end((err, res) => {
          if(err) {
            console.log('error', err);
            return;
          }
          const graph = JSON.parse(res.text);
          const {connected, disconnected } = chordBackendToChordData(graph.graph)
          setNewState({connected, disconnected});
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
      disconnectedLabels: historyItem.disconnectedLabels,
    })
  }
  componentDidMount() {
    clearTimeout(this.timeout);
    this.tick();
    window.addEventListener('resize', () => {
      this.setState({
        windowWidth: window.innerWidth,
        windowHeight: window.innerHeight,
      })
    })
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
  getChordChartSize() {
    return Math.max(Math.min(this.state.windowWidth, this.state.windowHeight) * 0.5, 300)
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
        <div className="content">
          <div className="graphs">
            <div className="chord-container">
              <Chord
                width={this.getChordChartSize()}
                height={this.getChordChartSize()}
                matrix={this.state.data}
                keys={this.state.labels}
                margin={{
                    "top": 100,
                    "right": 100,
                    "bottom": 100,
                    "left": 100
                }}
                padAngle={0.04}
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
                    "itemWidth": 180,
                    "itemHeight": 14,
                    "symbolSize": 14,
                    "symbolShape": "circle"
                  }]}
              />
            </div>
            <div className="controls">
              <HistorySelect
                  hasPauseCursor={this.state.hasPauseCursor}
                  history={this.state.history}
                  onChange={this.onHistoryChange}/>
            </div>
          </div>
          <div className="disconnected-list">
            <h3>
              disconnected items
            </h3>
            <List items={this.state.disconnectedLabels}/>
          </div>
        </div>
      </div>
    );
  }
}

export default App;

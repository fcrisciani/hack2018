import React, { Component } from 'react';
import { Chord } from 'nivo';
// import {Chord} from './appComponents/nivo/lib/Chord';
import request from 'superagent';
import chordBackendToChordData from './utils/transforms/chordBackendToChordData';
import HistorySelect from './appComponents/HistorySelect';
import PlayPause from './appComponents/PlayPause';
import List from './appComponents/List';
import sample from './mock/sample.json';
import samplePods from './mock/sample_pods.json';
import './App.css';
import debounce from 'lodash/debounce'
import Dropdown from 'react-dropdown'
import 'react-dropdown/style.css'

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
const getEndpointPortion = (endpoint) => {
  const splitted = endpoint.split('/')
  const endpointPortion = splitted[splitted.length-1]
  return endpointPortion;
}
class App extends Component {
  constructor(props) {
    super(props);
    const endpoint = localStorage.getItem(ENDPOINT_LOCALSTORAGE_KEY) || 'http://52.42.55.249:10000/chord';
    const endpointPortion = getEndpointPortion(endpoint);
    this.state = ({
      endpointPortion,
      endpoint,
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
    const useSamplePods = false;
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
    const endpointPortion = getEndpointPortion(endpoint);
    this.setState({
      endpoint,
      endpointPortion
    })
    this.storeEndpoint(endpoint)
    this.tick();
  }
  getChordChartSize() {
    return Math.max(Math.min(this.state.windowWidth, this.state.windowHeight) * 1, 300)
  }
  onDropdownChange = (newVal) => {
    const endpointPortion = newVal.value;
    const endpoint = `http://52.42.55.249:10000/${endpointPortion}`
    this.setState({
      endpoint,
      endpointPortion,
    })
    this.storeEndpoint(endpoint)
    this.tick()
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
        <div style={{
          width: 300,
          margin: 'auto'
        }}>
          <Dropdown
            options={[
              'services',
              'pods',
              'external',
            ]}
            value={this.state.endpointPortion}
            onChange={this.onDropdownChange}
          />
        </div>


        <h3 className="title">
          <PlayPause onChange={this.onPlayPauseChange} playing={this.state.playing}/>
          {this.state.totalCount} connections @{this.state.time.toLocaleTimeString('en-US')}
        </h3>
        <div className="content">
          <div className="graphs">
            <div className="chord-container">
              <Chord
                tooltipFormat={this.tooltipFormat}

                onClick={(e) => {alert(e)}}
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
            <List 
              title={
                <h3>
                  Inactive Endoints
                </h3>
              }
              empty={
                <div>
                  All Ednpoints Connected
                </div>
              }
              items={this.state.disconnectedLabels}/>
          </div>
        </div>
        <label>
          edit endpoint
          <input
            style={{
              border: 0,
              width: 400,
              margin: '0 5px',
            }}
           onChange={this.onInputChange} value={this.state.endpoint} type="text"/>
        </label>
      </div>
    );
  }
}

export default App;

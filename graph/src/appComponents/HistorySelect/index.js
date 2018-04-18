import React, { Component } from 'react';

class HistoryButton extends Component {
  onClick = () => {
    this.props.onClick(this.props.item)
  }
  render() {
    return (
      <button onClick={this.onClick}>
        {this.props.item.time.toLocaleTimeString('en-US')}
      </button>
    )
  }
}

export default class HistorySelect extends Component {
  onClickHistoryButton = (item) => {
    this.props.onChange(item);
  }
  render() {

    return (
      <div style={{
        display: 'flex',
        flex: 1,
        justifyContent: 'center',
      }}>
        {
          this.props.history.map((instance) =>{
            return (
              <HistoryButton onClick={this.onClickHistoryButton} key={instance.time} item={instance}/>
            )
          })
        }
      </div>
    );
  }
}
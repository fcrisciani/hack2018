import React, { Component } from 'react';

export default class List extends Component {
  render() {
    if (!this.props.items.length) {
      return (
        <div>
          All items connected
        </div>
      )
    }
    return (
      <ul>
        {
          this.props.items.map((item, i) => {
            return (
              <li key={i}>
              {item}
              </li>
            );
          })
        }
      </ul>
    );
  }
}
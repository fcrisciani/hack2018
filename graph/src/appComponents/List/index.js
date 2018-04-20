import React, { Component } from 'react';

export default class List extends Component {
  render() {
    if (!this.props.items.length) {
      return this.props.empty;
    }
    return (
      <div className="list-container">
      <h3 className="list-title">
        {this.props.title}
      </h3>
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
      </div>
    );
  }
}
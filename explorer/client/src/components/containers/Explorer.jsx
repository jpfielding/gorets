import React from 'react';
import MetadataService from 'services/MetadataService';

export default class Explorer extends React.Component {

  componentDidMount() {
    console.log('Mounting');
    MetadataService
      .get()
      .then(response => response.json())
      .then(json => {
        console.log(json);
      });
  }

  render() {
    return (
      <h1>Explorer</h1>
    );
  }
}
